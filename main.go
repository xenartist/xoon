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

	//
	switchView := ui.CreateSwitchViewFunc(rightFlex, mainMenu)

	ui.SetupMenuItemSelection(mainMenu, switchView,
		solanaUI.ConfigFlex, xenblocksUI.ConfigFlex,
		solanaUI.LogView, xenblocksUI.LogView)

	mainFlex := tview.NewFlex().
		AddItem(mainMenu, 0, 1, true).
		AddItem(rightFlex, 0, 3, false)

	//Press q 4 times to quit this app
	ui.SetupInputCapture(app)
	// app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	if event.Rune() == 'q' {
	// 		now := time.Now()
	// 		if now.Sub(lastQuitTime) > time.Second {
	// 			quitCount = 1
	// 		} else {
	// 			quitCount++
	// 		}
	// 		lastQuitTime = now
	// 		if quitCount >= 4 {
	// 			xenblocks.KillMiningProcess() // Kill the mining process before exiting
	// 			app.Stop()
	// 			return nil
	// 		}
	// 	}
	// 	return event
	// })

	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
