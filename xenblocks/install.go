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

const (
	XENBLOCKS_MINER_DIR = "xenblocksMiner"
	CONFIG_FILE_NAME    = "config.txt"
)

// ReadConfigFile reads the content of config.txt
func ReadConfigFile(logView *tview.TextView, logMessage utils.LogMessageFunc) (string, error) {
	cwd, err := os.Getwd()
	logMessage(logView, "Debug: ReadConfigFile cwd is "+cwd)
	if err != nil {
		logMessage(logView, "Error getting current working directory: "+err.Error())
		return "", err
	}

	configPath := filepath.Join(cwd, XENBLOCKS_MINER_DIR, CONFIG_FILE_NAME)
	content, err := os.ReadFile(configPath)
	if err != nil {
		logMessage(logView, "Error reading config file: "+err.Error())
		return "", err
	}
	logMessage(logView, "Config file read successfully")
	return string(content), nil
}

// WriteConfigFile writes or updates the content of config.txt
func WriteConfigFile(content string, logView *tview.TextView, logMessage utils.LogMessageFunc) error {
	cwd, err := os.Getwd()
	if err != nil {
		logMessage(logView, "Error getting current working directory: "+err.Error())
		return err
	}

	configPath := filepath.Join(cwd, XENBLOCKS_MINER_DIR, CONFIG_FILE_NAME)
	err = os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		logMessage(logView, "Error writing config file: "+err.Error())
		return err
	}
	logMessage(logView, "Config file written successfully")
	return nil
}

// CreateXenblocksMinerDir creates the xenblocksMiner directory and config file if they don't exist
func CreateXenblocksMinerDir(logView *tview.TextView, logMessage utils.LogMessageFunc) error {
	cwd, err := os.Getwd()
	logMessage(logView, "debug:cwd is "+cwd)
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error getting current directory: %v", err))
		return err
	}

	xenblocksMinerPath := filepath.Join(cwd, XENBLOCKS_MINER_DIR)
	if _, err := os.Stat(xenblocksMinerPath); os.IsNotExist(err) {
		err = os.Mkdir(xenblocksMinerPath, 0755)
		if err != nil {
			logMessage(logView, fmt.Sprintf("Error creating xenblocksMiner directory: %v", err))
			return err
		}
		logMessage(logView, "xenblocksMiner directory created successfully")
	}

	configPath := filepath.Join(xenblocksMinerPath, CONFIG_FILE_NAME)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		content := "account_address=ETH address with uppercase and lowercase\ndevfee_permillage=2"
		err = os.WriteFile(configPath, []byte(content), 0644)
		if err != nil {
			logMessage(logView, fmt.Sprintf("Error creating config file: %v", err))
			return err
		}
		logMessage(logView, "Config file created successfully")
	}

	return nil
}

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
	cwd, err := os.Getwd()
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error getting current directory: %v", err))
		return false
	}

	cmd := exec.Command(filepath.Join(cwd, XENBLOCKS_MINER_DIR, "xenblocksMiner"), "-h")
	output, err := cmd.Output()
	if err == nil && strings.Contains(string(output), "XenblocksMiner") {
		logMessage(logView, "XenblocksMiner is already installed. No need to install again.")
		return true
	}
	return false
}

func downloadXENBLOCKS(logView *tview.TextView, logMessage utils.LogMessageFunc) (string, error) {
	url := "https://github.com/woodysoil/XenblocksMiner/releases/download/v1.3.1/xenblocksMiner-1.3.1-Linux.tar.gz"
	fileName := "xenblocksMiner-1.3.1-Linux.tar.gz"

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error getting current directory: %v", err))
		return "", err
	}

	// Construct the full file path
	filePath := filepath.Join(cwd, XENBLOCKS_MINER_DIR, fileName)

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
