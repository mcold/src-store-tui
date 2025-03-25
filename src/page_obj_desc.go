package main

import (
	"database/sql"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
			saveObjDesc()
			pageObjDesc.descArea.SetText("", true)
			pageProTree.Pages.SwitchToPage("src")
			if pageSrc.lSrc.GetItemCount() > 0 {
				pageSrc.lSrc.SetCurrentItem(1)
			}
			return nil
		}
		return event
	})

	pageObjDesc.Flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pageObjDesc.descArea, 0, 1, true)

	pageObjDesc.Flex.SetTitle("F4").
		SetTitleAlign(tview.AlignLeft)

	pageProTree.Pages.AddPage("proObjDesc", pageObjDesc.Flex, true, false)
}

func saveObjDesc() {
	query := "UPDATE obj" + "\n" +
		"SET comment = '" + pageObjDesc.descArea.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))

	_, err := database.Exec(query, "PRAGMA busy_timeout=30000;")
	check(err)

}

func (pageObjDesc *pageObjDescType) show() {
	setObjDesc()
	pageProTree.Pages.SwitchToPage("proObjDesc")
	app.SetFocus(pageObjDesc.Flex)
}

func setObjDesc() {
	query := `select comment
				from obj` +
		` where id = ` + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))

	obj, err := database.Query(query)
	check(err)

	obj.Next()
	var comment sql.NullString
	err = obj.Scan(&comment)

	pageObjDesc.descArea.SetText(comment.String, true)
	pageProTree.descArea.SetText(comment.String, true)
	obj.Close()

	if len(comment.String) == 0 {
		pageProTree.flTree.RemoveItem(pageProTree.descArea)
	} else {
		if pageProTree.descArea.GetDisabled() == true {
			pageProTree.descArea.SetDisabled(false)
			pageProTree.flTree.AddItem(pageProTree.descArea, 0, 3, false)
		}
	}
}
