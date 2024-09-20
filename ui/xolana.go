package ui

import (
	"xoon/utils"
	"xoon/xolana"

	"github.com/rivo/tview"
)

func CreateXolanaUI(app *tview.Application) ModuleUI {
	var moduleUI = CreateModuleUI("Xolana", app)

	form := tview.NewForm()

	var publicKey string
	form.AddInputField("SVM Public Address", "NOTE: 5 SOL per hour is allowed", 44, nil, func(text string) {
		publicKey = text
	})

	form.AddButton("Get Faucet", func() {
		xolana.GetFaucet(app, moduleUI.LogView, utils.LogMessage, publicKey)
	})

	contentFlex := tview.NewFlex().AddItem(form, 0, 1, true)

	moduleUI.ConfigFlex.AddItem(contentFlex, 0, 1, true)

	return moduleUI
}

func CreateXolanaConfigFlex(app *tview.Application, logView *tview.TextView) *tview.Flex {
	configFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	configFlex.SetBorder(true).SetTitle("Xolana Config")
	return configFlex
}
