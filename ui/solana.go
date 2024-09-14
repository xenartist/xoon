package ui

import (
	"github.com/rivo/tview"
)

func CreateSolanaUI(app *tview.Application) ModuleUI {
	// var moduleUI ModuleUI
	// actions := map[string]func(){
	// 	"Install": func() { solana.InstallSolanaCLI(app, moduleUI.LogView, utils.LogMessage) },
	// 	"Airdrop": func() { solana.Airdrop(app, moduleUI.LogView, utils.LogMessage) },
	// 	// Add more actions as needed
	// }
	var moduleUI = CreateModuleUI("Solana CLI", app)
	return moduleUI
}

func CreateSolanaConfigFlex(app *tview.Application, logView *tview.TextView) *tview.Flex {
	configFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	// actions := map[string]func(){
	// 	"Check Balance":    checkSolanaBalance,
	// 	"Send Transaction": sendSolanaTransaction,
	// 	// Add other Solana-specific actions here
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

	configFlex.SetBorder(true).SetTitle("Solana CLI Config")
	return configFlex
}
