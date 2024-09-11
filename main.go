package main

import (
	"xoon/solana"
	"xoon/ui"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	mainMenu := ui.CreateMainMenu()
	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	solanaLogView := ui.CreateLogView("Solana CLI Logs", app)
	solanaConfigFlex := ui.CreateConfigFlex("Solana CLI", app, solanaLogView, solana.InstallSolanaCLI)

	xenblockLogView := ui.CreateLogView("XENBLOCKS Logs", app)
	xenblockConfigFlex := ui.CreateConfigFlex("XENBLOCKS", app, xenblockLogView, nil)

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
