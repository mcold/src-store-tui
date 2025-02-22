package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

	pagePro.Pages.AddPage("proDesc", pageProDesc.Flex, true, false)
}
