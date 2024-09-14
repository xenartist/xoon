package xenblocks

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"xoon/utils"

	"github.com/rivo/tview"
)

// InstallXENBLOCKS handles the installation of XENBLOCKS
func InstallXENBLOCKS(app *tview.Application, logView *tview.TextView, logMessage utils.LogMessageFunc) {
	// Log the start of the installation process
	if isXENBLOCKSInstalled(logView, logMessage) {
		return
	}

	logMessage(logView, "Starting XenblocksMiner installation...")

	downloadedPath, err := downloadXENBLOCKS(logView, logMessage)
	if err != nil {
		return
	}

	// Extract XenblocksMiner
	xenblocksMinerPath, err := extractXENBLOCKS(logView, logMessage, downloadedPath)
	if err != nil {
		return
	}

	logMessage(logView, fmt.Sprintf("XenblocksMiner installed successfully at: %s", xenblocksMinerPath))
}

func isXENBLOCKSInstalled(logView *tview.TextView, logMessage utils.LogMessageFunc) bool {
	cmd := exec.Command("./xenblocksMiner", "-h")
	output, err := cmd.Output()
	if err == nil && strings.Contains(string(output), "XenblocksMiner") {
		logMessage(logView, "XenblocksMiner is already installed. No need to install again.")
		return true
	}
	return false
}

func downloadXENBLOCKS(logView *tview.TextView, logMessage utils.LogMessageFunc) (string, error) {
	url := "https://github.com/woodysoil/XenblocksMiner/releases/download/v1.4.0/xenblocksMiner-1.4.0-Linux.tar.gz"
	fileName := "xenblocksMiner-1.4.0-Linux.tar.gz"

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error getting current directory: %v", err))
		return "", err
	}

	// Construct the full file path
	filePath := filepath.Join(cwd, fileName)

	// Prepare the curl command
	cmd := exec.Command("curl", "-L", "-o", filePath, url)

	// Capture the command output
	output, err := cmd.CombinedOutput()
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error downloading file: %v\nOutput: %s", err, string(output)))
		return "", err
	}

	logMessage(logView, fmt.Sprintf("File downloaded successfully to: %s", filePath))
	return filePath, nil
}

func extractXENBLOCKS(logView *tview.TextView, logMessage utils.LogMessageFunc, downloadedPath string) (string, error) {
	// Get the directory of the downloaded file
	dir := filepath.Dir(downloadedPath)

	// Extract the tar.gz file
	cmd := exec.Command("tar", "-zxvf", downloadedPath, "-C", dir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error extracting file: %v\nOutput: %s", err, string(output)))
		return "", err
	}
	logMessage(logView, "File extracted successfully")

	// Construct the path to the extracted executable
	executablePath := filepath.Join(dir, "xenblocksMiner")

	// Make the file executable
	cmd = exec.Command("chmod", "+x", executablePath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error making file executable: %v\nOutput: %s", err, string(output)))
		return "", err
	}
	logMessage(logView, "File permissions updated successfully")

	// Return the path to the executable
	return executablePath, nil
}
