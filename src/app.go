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
			application.ConfirmQuit()
		case tcell.KeyF2:
			app.SetFocus(pagePro.lPro)
		case tcell.KeyF3:
			pagePro.Pages.SwitchToPage("proTree")
			app.SetFocus(pageProTree.trPro)
		case tcell.KeyF4:
			pagePro.Pages.SwitchToPage("proDesc")
			setProComment()
			app.SetFocus(pageProDesc.Flex)
		case tcell.KeyF5:
			pageProTree.Pages.SwitchToPage("src")
			app.SetFocus(pageSrc.Flex)
		case tcell.KeyF6:
			pageProTree.Pages.SwitchToPage("proObjDesc")
			app.SetFocus(pageObjDesc.Flex)
		case tcell.KeyF11:
			pageProTree.Pages.SwitchToPage("proSrcDesc")
			app.SetFocus(pageSrcDesc.Flex)
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
