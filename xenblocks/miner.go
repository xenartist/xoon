package xenblocks

import (
	"bufio"
	"io"
	"os/exec"
	"runtime"
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
			scanner := bufio.NewScanner(pipe)
			scanner.Split(bufio.ScanLines)

			var lastLine string
			ticker := time.NewTicker(200 * time.Millisecond)
			defer ticker.Stop()

			go func() {
				for range ticker.C {
					if lastLine != "" {
						app.QueueUpdateDraw(func() {
							logMessage(logView, "Current status: "+lastLine)
						})
					}
				}
			}()

			for scanner.Scan() {
				line := scanner.Text()
				lastLine = line
				app.QueueUpdateDraw(func() {
					logMessage(logView, line)
				})
			}

			if err := scanner.Err(); err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(logView, "Error reading pipe: "+err.Error())
				})
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
