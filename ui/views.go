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

func CreateConfigFlex(title string, app *tview.Application, logView *tview.TextView, actions map[string]func()) *tview.Flex {
	configFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	for actionName, actionFunc := range actions {
		button := tview.NewButton(actionName)
		button.SetSelectedFunc(func() {
			go func(action func()) {
				// Wrap the action with logging functionality
				action()
				app.QueueUpdateDraw(func() {
					utils.LogMessage(logView, "Action '"+actionName+"' triggered")
				})
			}(actionFunc)
		})
		configFlex.AddItem(button, 0, 1, false)

		// Add a spacer between buttons
		configFlex.AddItem(tview.NewBox(), 0, 1, false)
	}

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

func UpdateButtonLabel(flex *tview.Flex, buttonName string, newLabel string) {
	for i := 0; i < flex.GetItemCount(); i++ {
		item := flex.GetItem(i)
		if button, ok := item.(*tview.Button); ok {
			if button.GetLabel() == buttonName {
				button.SetLabel(newLabel)
				return
			}
		}
	}
}
