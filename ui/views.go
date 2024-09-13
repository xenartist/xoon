package ui

import (
	"strings"
	"sync"

	"github.com/rivo/tview"
)

type DynamicLogView struct {
	*tview.TextView
	lines    []string
	mutex    sync.Mutex
	app      *tview.Application
	maxLines int
}

func CreateDynamicLogView(title string, app *tview.Application, maxLines int) *DynamicLogView {
	logView := &DynamicLogView{
		TextView: tview.NewTextView().SetDynamicColors(true),
		lines:    make([]string, 0),
		app:      app,
		maxLines: maxLines,
	}
	logView.SetBorder(true).SetTitle(title)
	return logView
}

func (d *DynamicLogView) Write(p []byte) (n int, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	text := string(p)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if strings.Contains(line, "\r") {
			// This is likely a progress update
			parts := strings.Split(line, "\r")
			if len(d.lines) > 0 {
				d.lines[len(d.lines)-1] = parts[len(parts)-1]
			} else {
				d.lines = append(d.lines, parts[len(parts)-1])
			}
		} else {
			// Regular log line
			d.lines = append(d.lines, line)
		}
	}

	// Trim old lines if exceeding maxLines
	if len(d.lines) > d.maxLines {
		d.lines = d.lines[len(d.lines)-d.maxLines:]
	}

	d.app.QueueUpdateDraw(func() {
		d.Clear()
		for _, line := range d.lines {
			d.TextView.Write([]byte(line + "\n"))
		}
	})

	return len(p), nil
}

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
					//utils.LogMessage(logView, "Action '"+actionName+"' triggered")//temp disable
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
