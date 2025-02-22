package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

	pageObjDesc.Flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pageObjDesc.descArea, 0, 1, true)

	pageObjDesc.Flex.SetTitle("F6/F5").
		SetTitleAlign(tview.AlignLeft)

	pageProTree.Pages.AddPage("proObjDesc", pageObjDesc.Flex, true, false)
}
