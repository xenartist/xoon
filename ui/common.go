package ui

import (
	"github.com/rivo/tview"
)

var ModuleNames = []string{"Solana CLI", "XENBLOCKS"}

type ModuleUI struct {
	DashboardFlex *tview.Flex
	LogView       *tview.TextView
	ConfigFlex    *tview.Flex
}

func CreateModuleUI(name string, app *tview.Application) ModuleUI {
	logView := CreateLogView(name+" Logs", app)
	configFlex := CreateConfigFlex(name, app, logView)
	dashboardFlex := CreateDashboardFlex(name, app)
	return ModuleUI{
		DashboardFlex: dashboardFlex,
		LogView:       logView,
		ConfigFlex:    configFlex,
	}
}

func CreateDashboardFlex(title string, app *tview.Application) *tview.Flex {
	dashboardFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	// Add placeholder text for now
	placeholder := tview.NewTextView().SetText("Dashboard for " + title)
	dashboardFlex.AddItem(placeholder, 0, 1, false)

	dashboardFlex.SetBorder(true).SetTitle(title + " Dashboard")
	return dashboardFlex
}

func CreateLogView(title string, app *tview.Application) *tview.TextView {
	logView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	logView.SetBorder(true).SetTitle(title)
	return logView
}

func CreateConfigFlex(title string, app *tview.Application, logView *tview.TextView) *tview.Flex {

	switch title {
	case "XENBLOCKS":
		return CreateXenblocksConfigFlex(app, logView)
	case "Solana CLI":
		return CreateSolanaConfigFlex(app, logView)
	default:
		return createDefaultConfigFlex(title, app, logView)
	}
}

func createDefaultConfigFlex(title string, app *tview.Application, logView *tview.TextView) *tview.Flex {
	configFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	configFlex.SetBorder(true).SetTitle(title + " Config")
	return configFlex
}

func CreateSwitchViewFunc(rightFlex *tview.Flex, mainMenu *tview.List) func(*tview.Flex, *tview.Flex, *tview.TextView) {
	return func(dashboardFlex *tview.Flex, configFlex *tview.Flex, logView *tview.TextView) {
		rightFlex.Clear()
		rightFlex.
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(dashboardFlex, 0, 1, false).
				AddItem(configFlex, 0, 1, false),
				0, 1, false).
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
