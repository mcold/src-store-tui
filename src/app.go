package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os"
)

type applicationType struct {
	pages *tview.Pages
}

var app *tview.Application

func (application *applicationType) init() {
	file, err := os.OpenFile("app.log", os.O_TRUNC|os.O_CREATE, 0666)
	check(err)
	log.SetOutput(file)

	app = tview.NewApplication()

	application.pages = tview.NewPages()
	pageMain.build()

	application.registerGlobalShortcuts()

	if err := app.SetRoot(application.pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}

func (application *applicationType) registerGlobalShortcuts() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			beforeSwitch()
			application.ConfirmQuit()
		case tcell.KeyF2:
			beforeSwitch()
			app.SetFocus(pagePro.lPro)
		case tcell.KeyF3:
			beforeSwitch()
			pageProTree.show()
		case tcell.KeyF4:
			beforeSwitch()
			pageSrc.show()
		case tcell.KeyF5:
			beforeSwitch()
			pageProDesc.show()
		case tcell.KeyF6:
			beforeSwitch()
			pageObjDesc.show()
		case tcell.KeyF10:
			beforeSwitch()
			pageExec.show()
		case tcell.KeyF11:
			pagePro.flexPro.RemoveItem(pagePro.flListPro)
		case tcell.KeyF12:
			if pagePro.flexPro.GetItemCount() < 2 {
				pagePro.flexPro.RemoveItem(pagePro.Pages)
				pagePro.flexPro.AddItem(pagePro.flListPro, 0, 1, true)
				pagePro.flexPro.AddItem(pagePro.Pages, 0, 4, true)
			}

		default:
			return event
		}
		return nil
	})
}

func (application *applicationType) ConfirmQuit() {
	pageConfirm.show("Are you sure you want to exit?", application.Quit)
}

func (application *applicationType) Quit() {
	if database.DB != nil {
		database.DB.Close()
	}
	app.Stop()
}

func beforeSwitch() {
	saveObjDesc()
}
