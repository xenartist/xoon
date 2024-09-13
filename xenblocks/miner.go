package xenblocks

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
	"xoon/utils"

	"github.com/rivo/tview"
)

func StartMining(app *tview.Application, logView *tview.TextView, logMessage utils.LogMessageFunc) {
	go func() {
		// Create the command
		cmd := exec.Command("./xenblocksMiner", "--minerAddr", "0x970Ce544847B0E314eA357e609A0C0cA4D9fD823", "--totalDevFee", "1")

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
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						app.QueueUpdateDraw(func() {
							logMessage(logView, "Error reading pipe: "+err.Error())
						})
					}
					break
				}
				line = strings.TrimSpace(line)
				if line != "" {
					app.QueueUpdateDraw(func() {
						logMessage(logView, line)
					})
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
