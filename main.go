package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create a text view for logs
	logView := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	logView.SetBorder(true).SetTitle("Logs")

	// Function to log messages
	logMessage := func(message string) {
		fmt.Fprintf(logView, "%s\n", message)
	}

	// Create the "Install Solana CLI" button
	installButton := tview.NewButton("Install Solana CLI")
	installButton.SetSelectedFunc(func() {

		var output []byte // Declare output here if not already declared
		var err error     // Declare err here if not already declared

		// Check if Solana is already installed
		cmd := exec.Command("solana", "--version")
		output, err = cmd.Output()
		if err == nil && strings.Contains(string(output), "solana-cli") {
			logMessage("Solana is already installed. No need to install again.")
			logMessage(fmt.Sprintf("Current version: %s", strings.TrimSpace(string(output))))
			return
		}

		logMessage("Starting Solana CLI installation...")

		// Download the installation script
		curlCmd := exec.Command("curl", "-sSfL", "https://release.solana.com/v1.18.22/install")
		output, err = curlCmd.Output()
		if err != nil {
			logMessage(fmt.Sprintf("Error downloading script: %v", err))
			return
		}

		// Save the script to a temporary file
		tmpfile, err := os.CreateTemp("", "solana-install-*.sh")
		if err != nil {
			logMessage(fmt.Sprintf("Error creating temp file: %v", err))
			return
		}
		defer tmpfile.Close()

		if _, err := tmpfile.Write(output); err != nil {
			logMessage(fmt.Sprintf("Error writing to temp file: %v", err))
			return
		}

		// Execute the installation script
		go func() {
			combinedCmd := exec.Command("sh", tmpfile.Name())

			// Create pipes for stdout and stderr
			stdout, err := combinedCmd.StdoutPipe()
			if err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(fmt.Sprintf("Error creating stdout pipe: %v", err))
				})
				return
			}
			stderr, err := combinedCmd.StderrPipe()
			if err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(fmt.Sprintf("Error creating stderr pipe: %v", err))
				})
				return
			}

			// Start the command
			if err := combinedCmd.Start(); err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(fmt.Sprintf("Error starting script: %v", err))
				})
				return
			}

			// Function to read from a pipe and update the UI
			readAndLog := func(pipe io.Reader) {
				scanner := bufio.NewScanner(pipe)
				for scanner.Scan() {
					line := scanner.Text()
					app.QueueUpdateDraw(func() {
						logMessage(line)
					})
				}
			}

			// Read from stdout and stderr concurrently
			go readAndLog(stdout)
			go readAndLog(stderr)

			// Wait for the command to finish
			if err := combinedCmd.Wait(); err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(fmt.Sprintf("Error executing script: %v", err))
				})
			}

			// Continue with the rest of the installation process
			err = addSolanaToPath(logMessage)
			if err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(fmt.Sprintf("Error adding Solana to PATH: %v", err))
				})
			} else {
				app.QueueUpdateDraw(func() {
					logMessage("Solana path added to user's profile.")
					logMessage("Please reload your shell configuration or restart your terminal for the changes to take effect.")
				})
			}

			app.QueueUpdateDraw(func() {
				logMessage("Solana CLI installation completed.")
			})
		}()

		logMessage("Solana CLI installation started. Please wait...")

	})

	// Create the "Quit" button
	quitButton := tview.NewButton("Quit")
	quitButton.SetSelectedFunc(func() {
		app.Stop()
	})

	// Create a flex layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(installButton, 1, 1, true).
		AddItem(logView, 0, 1, false).
		AddItem(quitButton, 1, 1, true)

		// Set a minimum size for the application
	// app.SetMinSize(60, 20)

	// Run the application
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func addSolanaToPath(logMessage func(string)) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logMessage(fmt.Sprintf("Error getting user home directory: %v", err))
		return err
	}

	profilePath := filepath.Join(homeDir, ".bashrc")
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		profilePath = filepath.Join(homeDir, ".profile")
	}

	logMessage(fmt.Sprintf("Using profile file: %s", profilePath))

	solanaPath := fmt.Sprintf("export PATH=\"%s/.local/share/solana/install/active_release/bin:$PATH\"", homeDir)

	content, err := os.ReadFile(profilePath)
	if err != nil {
		logMessage(fmt.Sprintf("Error reading profile file: %v", err))
		return err
	}

	if strings.Contains(string(content), solanaPath) {
		logMessage("Solana path already exists in the profile. No changes made.")
		return nil
	}

	f, err := os.OpenFile(profilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logMessage(fmt.Sprintf("Error opening profile file: %v", err))
		return err
	}
	defer f.Close()

	if _, err = f.WriteString("\n" + solanaPath + "\n"); err != nil {
		logMessage(fmt.Sprintf("Error writing to profile file: %v", err))
		return err
	}

	logMessage("Solana path successfully added to the profile.")
	return nil
}
