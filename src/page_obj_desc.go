package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
)

type pageObjDescType struct {
	descArea *tview.TextArea
	*tview.Flex
}

var pageObjDesc pageObjDescType

func (pageObjDesc *pageObjDescType) build() {

	pageObjDesc.descArea = tview.NewTextArea()

	pageObjDesc.descArea.SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 1, 1)

	pageObjDesc.descArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Rune() == 's' && event.Modifiers() == tcell.ModAlt {
			saveProDesc()
			pageObjDesc.descArea.SetText("", true)
			pageProTree.Pages.SwitchToPage("src")
			return nil
		}
		return event
	})

	pageObjDesc.Flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pageObjDesc.descArea, 0, 1, true)

	pageObjDesc.Flex.SetTitle("F6/F5").
		SetTitleAlign(tview.AlignLeft)

	pageProTree.Pages.AddPage("proObjDesc", pageObjDesc.Flex, true, false)
}

func saveObjDesc() {

	query := "UPDATE prj" + "\n" +
		"SET comment = '" + pageProDesc.descArea.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])

	log.Println(query)

	_, err := database.Exec(query, "PRAGMA busy_timeout=30000;")
	check(err)

}

func (pageObjDesc *pageObjDescType) show() {
	pageProTree.Pages.SwitchToPage("proObjDesc")
	app.SetFocus(pageObjDesc.Flex)
}
