package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

type pageProDescType struct {
	descArea *tview.TextArea
	*tview.Flex
}

var pageProDesc pageProDescType

func (pageProDesc *pageProDescType) build() {

	pageProDesc.descArea = tview.NewTextArea()

	pageProDesc.descArea.SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 1, 1)

	pageProDesc.Flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pageProDesc.descArea, 0, 1, true)

	pageProDesc.Flex.SetTitle("F4/F3").
		SetTitleAlign(tview.AlignLeft)

	pageProDesc.descArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 's' && event.Modifiers() == tcell.ModAlt {
			saveProDesc()
			pageProDesc.descArea.SetText("", true)
			pagePro.Pages.SwitchToPage("proTree")
			return nil
		}
		if event.Key() == tcell.KeyEsc {
			pageProTree.Pages.SwitchToPage("src")
			app.SetFocus(pageSrc.Flex)
		}
		return event
	})

	pagePro.Pages.AddPage("proDesc", pageProDesc.Flex, true, false)
}

func saveProDesc() {

	query := "UPDATE prj" + "\n" +
		"SET comment = '" + pageProDesc.descArea.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])

	_, err := database.Exec(query, "PRAGMA busy_timeout=30000;")

	check(err)

}

func (pageProDesc *pageProDescType) show() {
	pagePro.Pages.SwitchToPage("proDesc")
	setProComment()
	app.SetFocus(pageProDesc.Flex)

}
