package xenblocks

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"xoon/utils"

	"github.com/rivo/tview"
)

var isMining bool = false

func StartMining(app *tview.Application, logView *tview.TextView, logMessage utils.LogMessageFunc) {

	isMining = true

	go func() {
		// Create the command
		cmd := exec.Command("./xenblocksMiner", "--minerAddr", "0x970Ce544847B0E314eA357e609A0C0cA4D9fD823", "--totalDevFee", "1", "--execute")

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
					logMessage(logView, fmt.Sprintf("Debug: Read %d bytes", n))

					lines := strings.Split(output, "\n")
					for _, line := range lines {
						line = strings.TrimSpace(line)
						if line != "" {
							logMessage(logView, line)
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
