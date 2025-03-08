package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
)

type pageSrcDescType struct {
	descArea *tview.TextArea
	*tview.Flex
}

var pageSrcDesc pageSrcDescType

func (pageSrcDesc *pageSrcDescType) build() {

	pageSrcDesc.descArea = tview.NewTextArea()

	pageSrcDesc.descArea.SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 1, 1)

	pageSrcDesc.descArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Rune() == 's' && event.Modifiers() == tcell.ModAlt {
			saveSrcDesc()
			pageSrcDesc.descArea.SetText("", true)
			curNode := pageProTree.trPro.GetCurrentNode()
			curNode.ClearChildren()
			pageSrc.lSrc.Clear()
			setFileSrc(curNode.GetReference().(int))
			pageProTree.Pages.SwitchToPage("src")
			return nil
		}
		return event
	})

	pageSrcDesc.Flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pageSrcDesc.descArea, 0, 1, true)

	pageSrcDesc.Flex.SetTitle("F11").
		SetTitleAlign(tview.AlignLeft)

	pageProTree.Pages.AddPage("proSrcDesc", pageSrcDesc.Flex, true, false)
}

func saveSrcDesc() {
	log.Println("-------------------------------")
	log.Println("saveSrcDesc")
	log.Println("---------------------")

	query := "UPDATE src" + "\n" +
		"SET comment = '" + pageSrcDesc.descArea.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pageSrc.mPosId[pageSrc.lSrc.GetCurrentItem()])

	log.Println(query)

	_, err := database.Exec(query)
	check(err)

	log.Println("-------------------------------")

}

func (pageSrcDesc *pageSrcDescType) show() {
	pageProTree.Pages.SwitchToPage("proSrcDesc")
	app.SetFocus(pageSrcDesc.Flex)
}
