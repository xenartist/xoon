package main

import (
	"xoon/ui"
	"xoon/utils"

	"github.com/rivo/tview"
)

func main() {
	utils.XoosInit()

	app := tview.NewApplication()

	mainMenu := ui.CreateMainMenu()
	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	solanaUI := ui.CreateSolanaUI(app)
	xenblocksUI := ui.CreateXenblocksUI(app)
	xolanaUI := ui.CreateXolanaUI(app)

	//click items in mainmenu to swith views
	switchView := ui.CreateSwitchViewFunc(rightFlex, mainMenu)

	modules := []ui.ModuleUI{
		{
			DashboardFlex: solanaUI.DashboardFlex,
			ConfigFlex:    solanaUI.ConfigFlex,
			LogView:       solanaUI.LogView,
		},
		{
			DashboardFlex: xolanaUI.DashboardFlex,
			ConfigFlex:    xolanaUI.ConfigFlex,
			LogView:       xolanaUI.LogView,
		},
		{
			DashboardFlex: xenblocksUI.DashboardFlex,
			ConfigFlex:    xenblocksUI.ConfigFlex,
			LogView:       xenblocksUI.LogView,
		},
	}

	ui.SetupMenuItemSelection(mainMenu, switchView, modules)

	mainFlex := tview.NewFlex().
		AddItem(mainMenu, 0, 1, true).
		AddItem(rightFlex, 0, 3, false)

	//input capture, eg. press 4 times q to quit
	ui.SetupInputCapture(app)

	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
