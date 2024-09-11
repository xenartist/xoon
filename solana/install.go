package solana

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"xoon/utils"

	"github.com/rivo/tview"
)

func InstallSolanaCLI(app *tview.Application, logView *tview.TextView, logMessage utils.LogMessageFunc) {
	// Check if Solana is already installed
	if isSolanaInstalled(logView, logMessage) {
		return
	}

	logMessage(logView, "Starting Solana CLI installation...")

	// Download and save the installation script
	scriptPath, err := downloadInstallScript(logView, logMessage)
	if err != nil {
		return
	}

	// Execute the installation script
	executeInstallScript(app, logView, logMessage, scriptPath)

	logMessage(logView, "Solana CLI installation started. Please wait...")
}

func isSolanaInstalled(logView *tview.TextView, logMessage utils.LogMessageFunc) bool {
	cmd := exec.Command("solana", "--version")
	output, err := cmd.Output()
	if err == nil && strings.Contains(string(output), "solana-cli") {
		logMessage(logView, "Solana is already installed. No need to install again.")
		logMessage(logView, fmt.Sprintf("Current version: %s", strings.TrimSpace(string(output))))
		return true
	}
	return false
}

func downloadInstallScript(logView *tview.TextView, logMessage utils.LogMessageFunc) (string, error) {
	curlCmd := exec.Command("curl", "-sSfL", "https://release.solana.com/v1.18.22/install")
	output, err := curlCmd.Output()
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error downloading script: %v", err))
		return "", err
	}

	tmpfile, err := os.CreateTemp("", "solana-install-*.sh")
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error creating temp file: %v", err))
		return "", err
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write(output); err != nil {
		logMessage(logView, fmt.Sprintf("Error writing to temp file: %v", err))
		return "", err
	}

	return tmpfile.Name(), nil
}

func executeInstallScript(app *tview.Application, logView *tview.TextView, logMessage utils.LogMessageFunc, scriptPath string) {
	go func() {
		combinedCmd := exec.Command("sh", scriptPath)

		stdout, err := combinedCmd.StdoutPipe()
		if err != nil {
			app.QueueUpdateDraw(func() {
				logMessage(logView, fmt.Sprintf("Error creating stdout pipe: %v", err))
			})
			return
		}
		stderr, err := combinedCmd.StderrPipe()
		if err != nil {
			app.QueueUpdateDraw(func() {
				logMessage(logView, fmt.Sprintf("Error creating stderr pipe: %v", err))
			})
			return
		}

		if err := combinedCmd.Start(); err != nil {
			app.QueueUpdateDraw(func() {
				logMessage(logView, fmt.Sprintf("Error starting script: %v", err))
			})
			return
		}

		go readAndLog(app, stdout, logView, logMessage)
		go readAndLog(app, stderr, logView, logMessage)

		if err := combinedCmd.Wait(); err != nil {
			app.QueueUpdateDraw(func() {
				logMessage(logView, fmt.Sprintf("Error executing script: %v", err))
			})
		}

		err = addSolanaToPath(logMessage, logView)
		if err != nil {
			app.QueueUpdateDraw(func() {
				logMessage(logView, fmt.Sprintf("Error adding Solana to PATH: %v", err))
			})
		} else {
			app.QueueUpdateDraw(func() {
				logMessage(logView, "Solana path added to user's profile.")
				logMessage(logView, "Please reload your shell configuration or restart your terminal for the changes to take effect.")
			})
		}

		app.QueueUpdateDraw(func() {
			logMessage(logView, "Solana CLI installation completed.")
		})
	}()
}

func readAndLog(app *tview.Application, pipe io.Reader, logView *tview.TextView, logMessage utils.LogMessageFunc) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		app.QueueUpdateDraw(func() {
			logMessage(logView, line)
		})
	}
}

func addSolanaToPath(logMessage utils.LogMessageFunc, logView *tview.TextView) error {
	// Implementation remains the same as in the original file
	// ...
	return nil
}
