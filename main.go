package main

import (
	"time"
	"xoon/solana"
	"xoon/ui"
	"xoon/utils"
	"xoon/xenblocks"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	var quitCount int
	var lastQuitTime time.Time

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
	var xenblockConfigFlex *tview.Flex // Declare xenblockConfigFlex here
	xenblockActions := map[string]func(){
		"Install": func() { xenblocks.InstallXENBLOCKS(app, xenblockLogView, utils.LogMessage) },
		"Start Mining": func() {
			if !xenblocks.IsMining() {
				xenblocks.StartMining(app, xenblockLogView, utils.LogMessage)
				ui.UpdateButtonLabel(xenblockConfigFlex, "Start Mining", "Stop Mining")
			} else {
				xenblocks.StopMining(app, xenblockLogView, utils.LogMessage)
				ui.UpdateButtonLabel(xenblockConfigFlex, "Stop Mining", "Start Mining")
			}
		},
		// Add more actions as needed
	}
	xenblockConfigFlex = ui.CreateConfigFlex("XENBLOCKS", app, xenblockLogView, xenblockActions)

	//
	switchView := ui.CreateSwitchViewFunc(rightFlex, mainMenu)

	ui.SetupMenuItemSelection(mainMenu, switchView, solanaConfigFlex, xenblockConfigFlex, solanaLogView, xenblockLogView)

	mainFlex := tview.NewFlex().
		AddItem(mainMenu, 0, 1, true).
		AddItem(rightFlex, 0, 3, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			now := time.Now()
			if now.Sub(lastQuitTime) > time.Second {
				quitCount = 1
			} else {
				quitCount++
			}
			lastQuitTime = now
			if quitCount >= 4 {
				xenblocks.KillMiningProcess() // Kill the mining process before exiting
				app.Stop()
				return nil
			}
		}
		return event
	})

	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
