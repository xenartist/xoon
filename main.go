package main

import (
	"xoon/solana"
	"xoon/ui"
	"xoon/utils"
	"xoon/xenblocks"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	mainMenu := ui.CreateMainMenu()
	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	//Solana CLI
	solanaLogView := ui.CreateLogView("Solana CLI Logs", app)
	solanaActions := map[string]func(){
		"Install": func() { solana.InstallSolanaCLI(app, solanaLogView, utils.LogMessage) },
		"Airdrop": func() { solana.Airdrop(app, solanaLogView, utils.LogMessage) },
		// Add more actions as needed
	}
	solanaConfigFlex := ui.CreateConfigFlex("Solana CLI", app, solanaLogView, solanaActions)

	//XENBLOCKS
	xenblockLogView := ui.CreateLogView("XENBLOCKS Logs", app)
	xenblockActions := map[string]func(){
		"Install": func() { xenblocks.InstallXENBLOCKS(app, xenblockLogView, utils.LogMessage) },
		// Add more actions as needed
	}
	xenblockConfigFlex := ui.CreateConfigFlex("XENBLOCKS", app, xenblockLogView, xenblockActions)

	//
	switchView := ui.CreateSwitchViewFunc(rightFlex, mainMenu)

	ui.SetupMenuItemSelection(mainMenu, switchView, solanaConfigFlex, xenblockConfigFlex, solanaLogView, xenblockLogView)

	mainFlex := tview.NewFlex().
		AddItem(mainMenu, 0, 1, true).
		AddItem(rightFlex, 0, 3, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF10 && event.Modifiers() == tcell.ModCtrl {
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
