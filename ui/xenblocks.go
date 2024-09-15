package ui

import (
	"xoon/utils"
	"xoon/xenblocks"

	"github.com/rivo/tview"
)

func CreateXenblocksUI(app *tview.Application) ModuleUI {
	var moduleUI = CreateModuleUI("XENBLOCKS", app)

	// Create form
	form := tview.NewForm().
		AddInputField("EIP-55 Address", "", 44, nil, nil).
		AddInputField("RPC Link", "http://xenblocks.io", 44, nil, nil).
		AddInputField("Dev Fee (0-1000)", "2", 4, nil, nil).
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
