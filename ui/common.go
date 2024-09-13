package ui

import (
	"github.com/rivo/tview"
)

type ModuleUI struct {
	LogView    *tview.TextView
	ConfigFlex *tview.Flex
}

func CreateModuleUI(name string, app *tview.Application, actions map[string]func()) ModuleUI {
	logView := CreateLogView(name+" Logs", app)
	configFlex := CreateConfigFlex(name, app, logView, actions)
	return ModuleUI{
		LogView:    logView,
		ConfigFlex: configFlex,
	}
}
