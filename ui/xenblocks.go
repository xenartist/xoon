package ui

import (
	"strings"
	"xoon/utils"
	"xoon/xenblocks"

	"github.com/rivo/tview"
)

func CreateXenblocksUI(app *tview.Application) ModuleUI {
	var moduleUI = CreateModuleUI("XENBLOCKS", app)

	// Ensure xenblocksMiner directory and config.txt exist
	if err := xenblocks.CreateXenblocksMinerDir(moduleUI.LogView, utils.LogMessage); err != nil {
		utils.LogMessage(moduleUI.LogView, "Error creating xenblocksMiner directory: "+err.Error())
	}

	// Read config file
	config, err := xenblocks.ReadConfigFile(moduleUI.LogView, utils.LogMessage)

	// Default values
	accountAddress := ""
	devFee := "2"

	if err == nil {
		// Parse config
		for _, line := range strings.Split(config, "\n") {
			if strings.HasPrefix(line, "account_address=") {
				accountAddress = strings.TrimPrefix(line, "account_address=")
			} else if strings.HasPrefix(line, "devfee_permillage=") {
				devFee = strings.TrimPrefix(line, "devfee_permillage=")
			}
		}
	} else {
		utils.LogMessage(moduleUI.LogView, "Error reading config file: "+err.Error())
	}

	// Create form
	form := tview.NewForm().
		AddInputField("EIP-55 Address", accountAddress, 44, nil, nil).
		AddInputField("RPC Link", "http://xenblocks.io", 44, nil, nil).
		AddInputField("Dev Fee (0-1000)", devFee, 4, nil, nil).
		AddButton("Install Miner", func() { xenblocks.InstallXENBLOCKS(app, moduleUI.LogView, utils.LogMessage) }).
		AddButton("Save Config", nil).
		AddButton("Start Mining", func() {
			if !xenblocks.IsMining() {
				xenblocks.StartMining(app, moduleUI.LogView, utils.LogMessage)
			}
		}).
		AddButton("Stop Mining", func() {
			if xenblocks.IsMining() {
				xenblocks.StopMining(app, moduleUI.LogView, utils.LogMessage)
			}
		})

	contentFlex := tview.NewFlex().AddItem(form, 0, 1, true)

	moduleUI.ConfigFlex.AddItem(contentFlex, 0, 1, true)

	return moduleUI
}

func CreateXenblocksConfigFlex(app *tview.Application, logView *tview.TextView) *tview.Flex {
	configFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	configFlex.SetBorder(true).SetTitle("Xenblocks Config")
	return configFlex
}
