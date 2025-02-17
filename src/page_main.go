package main

import "github.com/rivo/tview"

type pageMainType struct {
	pages *tview.Pages
	*tview.Flex
}

var pageMain pageMainType

func (pageMain *pageMainType) build() {

	err := database.Connect()
	if err != nil {
		panic(err)
	}

	pageMain.pages = tview.NewPages()

	pagePro.build()

	pageMain.Flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pageMain.pages, 0, 1, true)

	app.SetFocus(pagePro.lPro)

	application.pages.AddPage("main", pageMain.Flex, true, true)
}
