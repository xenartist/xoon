package ui

import (
	"xoon/utils"
	"xoon/xenblocks"

	"github.com/rivo/tview"
)

func CreateXenblocksUI(app *tview.Application) ModuleUI {
	var moduleUI ModuleUI
	actions := map[string]func(){
		"Install": func() { xenblocks.InstallXENBLOCKS(app, moduleUI.LogView, utils.LogMessage) },
		"Start Mining": func() {
			if !xenblocks.IsMining() {
				xenblocks.StartMining(app, moduleUI.LogView, utils.LogMessage)
				UpdateButtonLabel(moduleUI.ConfigFlex, "Start Mining", "Stop Mining")
			} else {
				xenblocks.StopMining(app, moduleUI.LogView, utils.LogMessage)
				UpdateButtonLabel(moduleUI.ConfigFlex, "Stop Mining", "Start Mining")
			}
		},
		// Add more actions as needed
	}
	moduleUI = CreateModuleUI("XENBLOCKS", app, actions)
	return moduleUI
}
