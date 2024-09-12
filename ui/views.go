package ui

import (
	"xoon/utils"

	"github.com/rivo/tview"
)

func CreateLogView(title string, app *tview.Application) *tview.TextView {
	logView := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	logView.SetBorder(true).SetTitle(title)
	return logView
}

func CreateConfigFlex(title string, app *tview.Application, logView *tview.TextView, installFunc func(*tview.Application, *tview.TextView, utils.LogMessageFunc)) *tview.Flex {
	configFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	if title == "Solana CLI" {
		installButton := tview.NewButton("Install CLI")
		installButton.SetSelectedFunc(func() {
			installFunc(app, logView, utils.LogMessage)
		})
		configFlex.AddItem(installButton, 0, 1, false)
	} else {
		installButton := tview.NewButton("Install Miner")
		installButton.SetSelectedFunc(func() {
			installFunc(app, logView, utils.LogMessage)
		})
		configFlex.AddItem(installButton, 0, 1, false)

		// Add a spacer between buttons
		configFlex.AddItem(tview.NewBox(), 0, 1, false)

		// Add new Start/Stop Mining button
		miningButton := tview.NewButton("Start Mining")
		miningButton.SetSelectedFunc(func() {
			//TODO
		})
		configFlex.AddItem(miningButton, 0, 1, false)
	}

	configFlex.AddItem(tview.NewBox(), 0, 1, false)
	configFlex.SetBorder(true).SetTitle(title + " Config")
	return configFlex
}

func CreateSwitchViewFunc(rightFlex *tview.Flex, mainMenu *tview.List) func(*tview.Flex, *tview.TextView) {
	return func(configFlex *tview.Flex, logView *tview.TextView) {
		rightFlex.Clear()
		rightFlex.
			AddItem(configFlex, 3, 0, false).
			AddItem(logView, 0, 1, false)
		mainMenu.SetCurrentItem(0)
	}
}
