package xenblocks

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	"xoon/utils"

	"github.com/rivo/tview"
)

var isMining bool = false

func StartMining(app *tview.Application, logView *tview.TextView, logMessage utils.LogMessageFunc) {

	isMining = true

	go func() {
		// Read config file
		config, err := ReadConfigFile(logView, logMessage)
		if err != nil {
			logMessage(logView, "Error reading config file: "+err.Error())
			return
		}

		// Parse config
		devFee := "1"
		minerAddr := "0x970Ce544847B0E314eA357e609A0C0cA4D9fD823"
		for _, line := range strings.Split(config, "\n") {
			if strings.HasPrefix(line, "devfee_permillage=") {
				devFee = strings.TrimPrefix(line, "devfee_permillage=")
			} else if strings.HasPrefix(line, "account_address=") {
				minerAddr = strings.TrimPrefix(line, "account_address=")
			}
		}

		//Change working directory to xenblocksMiner
		cwd, err := os.Getwd()
		if err != nil {
			logMessage(logView, "Error getting current directory: "+err.Error())
			return
		}

		var xenblocksMinerPath = cwd
		if !strings.HasSuffix(cwd, XENBLOCKS_MINER_DIR) {
			xenblocksMinerPath = filepath.Join(cwd, XENBLOCKS_MINER_DIR)
		}

		err = os.Chdir(xenblocksMinerPath)
		if err != nil {
			logMessage(logView, "Error changing to xenblocksMiner directory: "+err.Error())
			return
		}

		// Create the command
		args := []string{
			"--minerAddr", minerAddr,
			"--totalDevFee", devFee,
		}

		if devFee != "0" {
			args = append(args, "--ecoDevAddr", "0x970Ce544847B0E314eA357e609A0C0cA4D9fD823")
		}

		args = append(args, "--execute")

		var executableName string
		if runtime.GOOS == "windows" {
			executableName = ".\\xenblocksMiner.exe"
		} else {
			executableName = "./xenblocksMiner"
		}

		cmd := exec.Command(executableName, args...)

		// Create pipes for both stdout and stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			logMessage(logView, "Error creating StdoutPipe: "+err.Error())
			return
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			logMessage(logView, "Error creating StderrPipe: "+err.Error())
			return
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			logMessage(logView, "Error starting miner: "+err.Error())
			return
		}

		var (
			lastUpdateTime time.Time
			mutex          sync.Mutex
		)

		logMessage(logView, "Debug: StartMining function called")

		// Function to read from a pipe and send to UI
		readPipe := func(pipe io.Reader) {
			reader := bufio.NewReader(pipe)
			buffer := make([]byte, 1024)

			logMessage(logView, "Debug: Starting to read pipe")

			for {
				n, err := reader.Read(buffer)
				if err != nil {
					if err == io.EOF {
						logMessage(logView, "Debug: EOF reached, waiting...")
						time.Sleep(100 * time.Millisecond)
						continue
					}
					logMessage(logView, fmt.Sprintf("Error reading pipe: %v", err))
					break
				}

				if n > 0 {
					output := string(buffer[:n])
					//logMessage(logView, fmt.Sprintf("Debug: Read %d bytes", n))

					lines := strings.Split(output, "\n")
					for _, line := range lines {
						line = strings.TrimSpace(line)
						if line != "" {
							if strings.Contains(line, "Mining:") {
								mutex.Lock()
								now := time.Now()
								if now.Sub(lastUpdateTime) >= time.Minute {
									lastUpdateTime = now
									logMessage(logView, line)
								}
								mutex.Unlock()
							} else if strings.Contains(line, "Ecosystem") {
								//skip
							} else {
								logMessage(logView, line)
							}
						}
					}
				} else {
					logMessage(logView, "Debug: No data read")
				}
			}
		}

		// Start goroutines to read from stdout and stderr
		go readPipe(stdout)
		go readPipe(stderr)

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			app.QueueUpdateDraw(func() {
				logMessage(logView, "Miner exited with error: "+err.Error())
			})
		} else {
			app.QueueUpdateDraw(func() {
				logMessage(logView, "Mining completed successfully")
			})
		}
	}()
}

func StopMining(app *tview.Application, logView *tview.TextView, logMessage utils.LogMessageFunc) {
	KillMiningProcess()
	logMessage(logView, "Mining stopped")
	isMining = false

	// Change directory to parent
	if err := os.Chdir(".."); err != nil {
		logMessage(logView, "Error changing to parent directory: "+err.Error())
	}
}

// KillMiningProcess stops all running xenblocksMiner processes
func KillMiningProcess() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("taskkill", "/F", "/IM", "xenblocksMiner*")
	} else {
		cmd = exec.Command("pkill", "-f", "xenblocksMiner")
	}
	_ = cmd.Run()
}

// IsMining returns the current mining status
func IsMining() bool {
	return isMining
}
