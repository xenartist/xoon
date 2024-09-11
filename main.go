package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create the main menu list
	mainMenu := tview.NewList().
		AddItem("Solana CLI", "", 'a', nil).
		AddItem("XENBLOCKS", "", 'e', nil).
		AddItem("Quit", "Press Ctrl+F10 to quit", 0, nil)

	mainMenu.SetBorder(true).SetTitle("xoon")

	// Function to log messages
	logMessage := func(logView *tview.TextView, message string) {
		// logView.Clear()
		fmt.Fprintf(logView, "%s\n", message)
	}

	// Create separate Config and Logs views for Solana CLI and XENBLOCKS
	solanaLogView := createLogView("Solana CLI Logs", app)
	solanaConfigFlex := createConfigFlex("Solana CLI", app, solanaLogView, logMessage)
	xenblockLogView := createLogView("XENBLOCKS Logs", app)
	xenblockConfigFlex := createConfigFlex("XENBLOCKS", app, xenblockLogView, logMessage)

	// Create a flex for the right side (config and log view)
	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)

		// Function to switch views
	switchView := func(configFlex *tview.Flex, logView *tview.TextView) {
		rightFlex.Clear()
		rightFlex.
			AddItem(configFlex, 3, 0, false).
			AddItem(logView, 0, 1, false)
		app.SetFocus(mainMenu)
	}

	// Create a text view for logs
	// logView := tview.NewTextView().
	// 	SetDynamicColors(true).
	// 	SetChangedFunc(func() {
	// 		app.Draw()
	// 	})
	// logView.SetBorder(true).SetTitle("Logs")

	// Set up menu item selection
	mainMenu.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		switch mainText {
		case "Solana CLI":
			switchView(solanaConfigFlex, solanaLogView)
		case "XENBLOCKS":
			switchView(xenblockConfigFlex, xenblockLogView)
		}
	})

	// Create the "Install Solana CLI" button
	installButton := tview.NewButton("Install Solana CLI")
	installButton.SetSelectedFunc(func() {

		var output []byte // Declare output here if not already declared
		var err error     // Declare err here if not already declared

		// Check if Solana is already installed
		cmd := exec.Command("solana", "--version")
		output, err = cmd.Output()
		if err == nil && strings.Contains(string(output), "solana-cli") {
			logMessage(solanaLogView, "Solana is already installed. No need to install again.")
			logMessage(solanaLogView, fmt.Sprintf("Current version: %s", strings.TrimSpace(string(output))))
			return
		}

		logMessage(solanaLogView, "Starting Solana CLI installation...")

		// Download the installation script
		curlCmd := exec.Command("curl", "-sSfL", "https://release.solana.com/v1.18.22/install")
		output, err = curlCmd.Output()
		if err != nil {
			logMessage(solanaLogView, fmt.Sprintf("Error downloading script: %v", err))
			return
		}

		// Save the script to a temporary file
		tmpfile, err := os.CreateTemp("", "solana-install-*.sh")
		if err != nil {
			logMessage(solanaLogView, fmt.Sprintf("Error creating temp file: %v", err))
			return
		}
		defer tmpfile.Close()

		if _, err := tmpfile.Write(output); err != nil {
			logMessage(solanaLogView, fmt.Sprintf("Error writing to temp file: %v", err))
			return
		}

		// Execute the installation script
		go func() {
			combinedCmd := exec.Command("sh", tmpfile.Name())

			// Create pipes for stdout and stderr
			stdout, err := combinedCmd.StdoutPipe()
			if err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(solanaLogView, fmt.Sprintf("Error creating stdout pipe: %v", err))
				})
				return
			}
			stderr, err := combinedCmd.StderrPipe()
			if err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(solanaLogView, fmt.Sprintf("Error creating stderr pipe: %v", err))
				})
				return
			}

			// Start the command
			if err := combinedCmd.Start(); err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(solanaLogView, fmt.Sprintf("Error starting script: %v", err))
				})
				return
			}

			// Function to read from a pipe and update the UI
			readAndLog := func(pipe io.Reader) {
				scanner := bufio.NewScanner(pipe)
				for scanner.Scan() {
					line := scanner.Text()
					app.QueueUpdateDraw(func() {
						logMessage(solanaLogView, line)
					})
				}
			}

			// Read from stdout and stderr concurrently
			go readAndLog(stdout)
			go readAndLog(stderr)

			// Wait for the command to finish
			if err := combinedCmd.Wait(); err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(solanaLogView, fmt.Sprintf("Error executing script: %v", err))
				})
			}

			// Continue with the rest of the installation process
			err = addSolanaToPath(func(message string) {
				logMessage(solanaLogView, message)
			})
			if err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(solanaLogView, fmt.Sprintf("Error adding Solana to PATH: %v", err))
				})
			} else {
				app.QueueUpdateDraw(func() {
					logMessage(solanaLogView, "Solana path added to user's profile.")
					logMessage(solanaLogView, "Please reload your shell configuration or restart your terminal for the changes to take effect.")
				})
			}

			app.QueueUpdateDraw(func() {
				logMessage(solanaLogView, "Solana CLI installation completed.")
			})
		}()

		logMessage(solanaLogView, "Solana CLI installation started. Please wait...")

	})

	// Initially show Solana CLI view
	switchView(solanaConfigFlex, solanaLogView)

	// Create the main flex layout
	mainFlex := tview.NewFlex().
		AddItem(mainMenu, 0, 1, true).
		AddItem(rightFlex, 0, 3, false)

	// Set up custom input capture for Ctrl+Shift+F10
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF10 && event.Modifiers() == tcell.ModCtrl {
			app.Stop()
			return nil
		}
		return event
	})

	// Run the application
	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

// Helper function to create a config flex
func createConfigFlex(title string, app *tview.Application, logView *tview.TextView, logMessage func(*tview.TextView, string)) *tview.Flex {
	configFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	if title == "Solana CLI" {
		installButton := createInstallSolanaButton(app, logView, logMessage)
		configFlex.AddItem(installButton, 0, 1, false)
	} else {
		// Add XENBLOCKS specific config items here
		configFlex.AddItem(tview.NewButton("XENBLOCKS Config"), 0, 1, false)
	}

	configFlex.AddItem(tview.NewBox(), 0, 1, false) // Placeholder for future controls
	configFlex.SetBorder(true).SetTitle(title + " Config")
	return configFlex
}

// Helper function to create a log view
func createLogView(title string, app *tview.Application) *tview.TextView {
	logView := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	logView.SetBorder(true).SetTitle(title)
	return logView
}

// Helper function to create the Install Solana CLI button
func createInstallSolanaButton(app *tview.Application, solanaLogView *tview.TextView, logMessage func(*tview.TextView, string)) *tview.Button {
	installButton := tview.NewButton("Install Solana CLI")
	installButton.SetSelectedFunc(func() {

		var output []byte // Declare output here if not already declared
		var err error     // Declare err here if not already declared

		// Check if Solana is already installed
		cmd := exec.Command("solana", "--version")
		output, err = cmd.Output()
		if err == nil && strings.Contains(string(output), "solana-cli") {
			logMessage(solanaLogView, "Solana is already installed. No need to install again.")
			logMessage(solanaLogView, fmt.Sprintf("Current version: %s", strings.TrimSpace(string(output))))
			return
		}

		logMessage(solanaLogView, "Starting Solana CLI installation...")

		// Download the installation script
		curlCmd := exec.Command("curl", "-sSfL", "https://release.solana.com/v1.18.22/install")
		output, err = curlCmd.Output()
		if err != nil {
			logMessage(solanaLogView, fmt.Sprintf("Error downloading script: %v", err))
			return
		}

		// Save the script to a temporary file
		tmpfile, err := os.CreateTemp("", "solana-install-*.sh")
		if err != nil {
			logMessage(solanaLogView, fmt.Sprintf("Error creating temp file: %v", err))
			return
		}
		defer tmpfile.Close()

		if _, err := tmpfile.Write(output); err != nil {
			logMessage(solanaLogView, fmt.Sprintf("Error writing to temp file: %v", err))
			return
		}

		// Execute the installation script
		go func() {
			combinedCmd := exec.Command("sh", tmpfile.Name())

			// Create pipes for stdout and stderr
			stdout, err := combinedCmd.StdoutPipe()
			if err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(solanaLogView, fmt.Sprintf("Error creating stdout pipe: %v", err))
				})
				return
			}
			stderr, err := combinedCmd.StderrPipe()
			if err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(solanaLogView, fmt.Sprintf("Error creating stderr pipe: %v", err))
				})
				return
			}

			// Start the command
			if err := combinedCmd.Start(); err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(solanaLogView, fmt.Sprintf("Error starting script: %v", err))
				})
				return
			}

			// Function to read from a pipe and update the UI
			readAndLog := func(pipe io.Reader) {
				scanner := bufio.NewScanner(pipe)
				for scanner.Scan() {
					line := scanner.Text()
					app.QueueUpdateDraw(func() {
						logMessage(solanaLogView, line)
					})
				}
			}

			// Read from stdout and stderr concurrently
			go readAndLog(stdout)
			go readAndLog(stderr)

			// Wait for the command to finish
			if err := combinedCmd.Wait(); err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(solanaLogView, fmt.Sprintf("Error executing script: %v", err))
				})
			}

			// Continue with the rest of the installation process
			err = addSolanaToPath(func(message string) {
				logMessage(solanaLogView, message)
			})
			if err != nil {
				app.QueueUpdateDraw(func() {
					logMessage(solanaLogView, fmt.Sprintf("Error adding Solana to PATH: %v", err))
				})
			} else {
				app.QueueUpdateDraw(func() {
					logMessage(solanaLogView, "Solana path added to user's profile.")
					logMessage(solanaLogView, "Please reload your shell configuration or restart your terminal for the changes to take effect.")
				})
			}

			app.QueueUpdateDraw(func() {
				logMessage(solanaLogView, "Solana CLI installation completed.")
			})
		}()

		logMessage(solanaLogView, "Solana CLI installation started. Please wait...")

	})

	return installButton
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
