package xolana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"xoon/utils"

	"github.com/rivo/tview"
)

func GetFaucet(app *tview.Application, logView *tview.TextView, logMessage utils.LogMessageFunc, publicKey string) {
	// Prepare the request payload
	payload := map[string]string{"pubkey": publicKey}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error creating JSON payload: %v", err))
		return
	}

	// Send POST request to the faucet
	resp, err := http.Post("https://xolana.xen.network/faucet", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error sending request to faucet: %v", err))
		return
	}
	defer resp.Body.Close()

	// Log the response status
	logMessage(logView, fmt.Sprintf("Faucet request status: %s", resp.Status))

	// TODO: Handle the response body if needed

}
