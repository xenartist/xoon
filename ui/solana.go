package ui

import (
	"github.com/rivo/tview"
)

func CreateSolanaUI(app *tview.Application) ModuleUI {
	var moduleUI = CreateModuleUI("Solana CLI", app)

	form := tview.NewForm()

	contentFlex := tview.NewFlex().AddItem(form, 0, 1, true)

	moduleUI.ConfigFlex.AddItem(contentFlex, 0, 1, true)

	return moduleUI
}

func CreateSolanaConfigFlex(app *tview.Application, logView *tview.TextView) *tview.Flex {
	configFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	configFlex.SetBorder(true).SetTitle("Solana CLI Config")
	return configFlex
}
