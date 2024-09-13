package ui

import (
	"time"
	"xoon/xenblocks"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateMainMenu() *tview.List {
	mainMenu := tview.NewList().
		AddItem("XOON!(Two Weeks)", "", 0, nil).
		AddItem("Solana CLI", "", 'a', nil).
		AddItem("XENBLOCKS", "", 'e', nil).
		AddItem("QUIT", "Press 'q' 4 times", 'q', nil)

	mainMenu.SetBorder(true).SetTitle("xoon")
	return mainMenu
}

func SetupMenuItemSelection(mainMenu *tview.List, switchView func(*tview.Flex, *tview.TextView), solanaConfigFlex, xenblockConfigFlex *tview.Flex, solanaLogView, xenblockLogView *tview.TextView) {
	mainMenu.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		switch mainText {
		case "Solana CLI":
			switchView(solanaConfigFlex, solanaLogView)
		case "XENBLOCKS":
			switchView(xenblockConfigFlex, xenblockLogView)
		}
	})
}

func SetupInputCapture(app *tview.Application) {
	var quitCount int
	var lastQuitTime time.Time

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
}
