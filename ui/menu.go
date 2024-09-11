package ui

import (
	"github.com/rivo/tview"
)

func CreateMainMenu() *tview.List {
	mainMenu := tview.NewList().
		AddItem("XOON!(Two Weeks)", "", 0, nil).
		AddItem("Solana CLI", "", 'a', nil).
		AddItem("XENBLOCKS", "", 'e', nil).
		AddItem("QUIT", "Ctrl+F10 to Quit", 0, nil)

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
