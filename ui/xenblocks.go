package ui

import (
	"xoon/utils"
	"xoon/xenblocks"

	"github.com/rivo/tview"
)

func CreateXenblocksUI(app *tview.Application) ModuleUI {
	// var moduleUI ModuleUI
	// actions := map[string]func(){
	// 	"Install": func() { xenblocks.InstallXENBLOCKS(app, moduleUI.LogView, utils.LogMessage) },
	// 	"Start Mining": func() {
	// 		if !xenblocks.IsMining() {
	// 			xenblocks.StartMining(app, moduleUI.LogView, utils.LogMessage)
	// 			UpdateButtonLabel(moduleUI.ConfigFlex, "Start Mining", "Stop Mining")
	// 		} else {
	// 			xenblocks.StopMining(app, moduleUI.LogView, utils.LogMessage)
	// 			UpdateButtonLabel(moduleUI.ConfigFlex, "Stop Mining", "Start Mining")
	// 		}
	// 	},
	// 	// Add more actions as needed
	// }
	var moduleUI = CreateModuleUI("XENBLOCKS", app)

	// Create form
	form := tview.NewForm().
		AddInputField("EIP-55 Address", "", 44, nil, nil).
		AddInputField("RPC Link", "http://xenblocks.io", 44, nil, nil).
		AddInputField("Dev Fee (0-1000)", "2", 4, nil, nil).
		AddButton("Save Config", nil).
		AddButton("Install CLI", func() { xenblocks.InstallXENBLOCKS(app, moduleUI.LogView, utils.LogMessage) }).
		AddButton("Start Mining", func() { xenblocks.StartMining(app, moduleUI.LogView, utils.LogMessage) })

	contentFlex := tview.NewFlex().AddItem(form, 0, 1, true)

	moduleUI.ConfigFlex.AddItem(contentFlex, 0, 1, true)

	return moduleUI
}

func CreateXenblocksConfigFlex(app *tview.Application, logView *tview.TextView) *tview.Flex {
	configFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	// actions := map[string]func(){
	// 	"Start Mining": startXenblocksMining,
	// 	"Stop Mining":  stopXenblocksMining,
	// 	// Add other Xenblocks-specific actions here
	// }

	// for actionName, actionFunc := range actions {
	// 	button := tview.NewButton(actionName)
	// 	button.SetSelectedFunc(func() {
	// 		go func(action func()) {
	// 			action()
	// 			app.QueueUpdateDraw(func() {
	// 				utils.LogMessage(logView, "Action '"+actionName+"' triggered")
	// 			})
	// 		}(actionFunc)
	// 	})
	// 	configFlex.AddItem(button, 0, 1, false)
	// 	configFlex.AddItem(tview.NewBox(), 0, 1, false) // Spacer
	// }

	configFlex.SetBorder(true).SetTitle("Xenblocks Config")
	return configFlex
}
