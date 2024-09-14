package main

import (
	"xoon/ui"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	mainMenu := ui.CreateMainMenu()
	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	solanaUI := ui.CreateSolanaUI(app)
	xenblocksUI := ui.CreateXenblocksUI(app)

	//click items in mainmenu to swith views
	switchView := ui.CreateSwitchViewFunc(rightFlex, mainMenu)

	ui.SetupMenuItemSelection(mainMenu, switchView,
		ui.ModuleUI{ConfigFlex: solanaUI.ConfigFlex, LogView: solanaUI.LogView},
		ui.ModuleUI{ConfigFlex: xenblocksUI.ConfigFlex, LogView: xenblocksUI.LogView})

	mainFlex := tview.NewFlex().
		AddItem(mainMenu, 0, 1, true).
		AddItem(rightFlex, 0, 3, false)

	//input capture, eg. press 4 times q to quit
	ui.SetupInputCapture(app)

	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
