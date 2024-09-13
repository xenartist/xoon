package ui

import (
	"xoon/solana"
	"xoon/utils"

	"github.com/rivo/tview"
)

func CreateSolanaUI(app *tview.Application) ModuleUI {
	var moduleUI ModuleUI
	actions := map[string]func(){
		"Install": func() { solana.InstallSolanaCLI(app, moduleUI.LogView, utils.LogMessage) },
		"Airdrop": func() { solana.Airdrop(app, moduleUI.LogView, utils.LogMessage) },
		// Add more actions as needed
	}
	moduleUI = CreateModuleUI("Solana CLI", app, actions)
	return moduleUI
}
