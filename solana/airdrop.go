package solana

import (
	"fmt"
	"os/exec"
	"xoon/utils"

	"github.com/rivo/tview"
)

func Airdrop(app *tview.Application, logView *tview.TextView, logMessage utils.LogMessageFunc) {
	logMessage(logView, "Starting Solana airdrop...")

	// Check if Solana CLI is installed
	if !isSolanaInstalled(logView, logMessage) {
		logMessage(logView, "Solana CLI is not installed. Please install it first.")
		return
	}

	// Get the current wallet address
	address, err := getCurrentWalletAddress(logView, logMessage)
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error getting wallet address: %v", err))
		return
	}

	// Request airdrop
	amount := "1" // Amount in SOL
	cmd := exec.Command("solana", "airdrop", amount, address)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logMessage(logView, fmt.Sprintf("Airdrop failed: %v", err))
		logMessage(logView, string(output))
		return
	}

	logMessage(logView, fmt.Sprintf("Airdrop successful: %s SOL to %s", amount, address))
	logMessage(logView, string(output))

	// Check balance after airdrop
	checkBalance(address, logView, logMessage)

	logMessage(logView, "Airdrop process completed.")
}

func getCurrentWalletAddress(logView *tview.TextView, logMessage utils.LogMessageFunc) (string, error) {
	cmd := exec.Command("solana", "address")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	address := string(output)
	logMessage(logView, fmt.Sprintf("Current wallet address: %s", address))
	return address, nil
}

func checkBalance(address string, logView *tview.TextView, logMessage utils.LogMessageFunc) {
	cmd := exec.Command("solana", "balance", address)
	output, err := cmd.Output()
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error checking balance: %v", err))
		return
	}
	logMessage(logView, fmt.Sprintf("Current balance: %s", string(output)))
}
