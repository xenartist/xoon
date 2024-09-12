package xenblocks

import (
	"bufio"
	"os/exec"
	"xoon/utils"

	"github.com/rivo/tview"
)

func StartMining(app *tview.Application, logView *tview.TextView, logMessage utils.LogMessageFunc) {
	go func() {
		// Create the command
		cmd := exec.Command("./xenblocksMiner", "--minerAddr", "0x970Ce544847B0E314eA357e609A0C0cA4D9fD823", "--totalDevFee", "1")

		// Create a pipe for the output
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			logMessage(logView, "Error creating StdoutPipe: "+err.Error())
			return
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			logMessage(logView, "Error starting miner: "+err.Error())
			return
		}

		// Create a scanner to read the output line by line
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			// Use app.QueueUpdateDraw to update UI from goroutine
			app.QueueUpdateDraw(func() {
				logMessage(logView, line)
			})
		}

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			logMessage(logView, "Miner exited with error: "+err.Error())
		} else {
			logMessage(logView, "Mining completed successfully")
		}
	}()
}
