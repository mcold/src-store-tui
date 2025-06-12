package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
)

type applicationType struct {
	pages *tview.Pages
}

var app *tview.Application

func (application *applicationType) init() {
	_, err := os.OpenFile("app.log", os.O_TRUNC|os.O_CREATE, 0666)
	check(err)
	
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
			application.ConfirmQuit()
		case tcell.KeyF2:
			app.SetFocus(pagePro.lPro)
		case tcell.KeyF3:
			pageProTree.show()
		case tcell.KeyF4:
			pageSrc.show()
		case tcell.KeyF5:
			pageProDesc.show()
		case tcell.KeyF6:
			pageObjDesc.show()
		case tcell.KeyF10:
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
