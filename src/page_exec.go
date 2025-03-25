package main

import (
	"database/sql"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

type pageExecType struct {
	execArea *tview.TextArea
	outArea  *tview.TextArea
	*tview.Flex
}

var pageExec pageExecType

func (pageExec *pageExecType) build() {

	pageExec.execArea = tview.NewTextArea()

	pageExec.execArea.
		SetBorderPadding(2, 2, 2, 2).
		SetBackgroundColor(tcell.ColorBlue)

	pageExec.execArea.SetBorderColor(tcell.ColorBlue)

	pageExec.execArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Rune() == 's' && event.Modifiers() == tcell.ModAlt {
			saveExec()
			pageExec.execArea.SetText("", true)
			pageExec.outArea.SetText("", true)
			return nil
		}
		if event.Rune() == 'c' && event.Modifiers() == tcell.ModAlt {
			err := clipboard.WriteAll(pageExec.execArea.GetText())
			check(err)
		}
		if event.Rune() == 'v' && event.Modifiers() == tcell.ModAlt {
			clipBoardContent, err := clipboard.ReadAll()
			check(err)

			pageExec.execArea.SetText(pageExec.execArea.GetText()+clipBoardContent, true)
		}
		return event
	})

	pageExec.outArea = tview.NewTextArea()

	pageExec.outArea.
		SetBorderPadding(2, 2, 2, 2).
		SetBackgroundColor(tcell.ColorBlue)

	pageExec.outArea.SetBorderColor(tcell.ColorBlue)

	pageExec.outArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Rune() == 's' && event.Modifiers() == tcell.ModAlt {
			saveExec()
			pageExec.execArea.SetText("", true)
			pageExec.outArea.SetText("", true)
			return nil
		}
		if event.Rune() == 'c' && event.Modifiers() == tcell.ModAlt {
			err := clipboard.WriteAll(pageExec.outArea.GetText())
			check(err)
		}
		if event.Rune() == 'v' && event.Modifiers() == tcell.ModAlt {
			clipBoardContent, err := clipboard.ReadAll()
			check(err)

			pageExec.outArea.SetText(pageExec.outArea.GetText()+clipBoardContent, true)
		}
		return event
	})

	pageExec.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageExec.execArea, 0, 1, true).
		AddItem(pageExec.outArea, 0, 1, true)

	pageProTree.Pages.AddPage("exec", pageExec.Flex, true, false)
}

func saveExec() {

	query := "UPDATE obj" + "\n" +
		"SET exec = '" + pageExec.execArea.GetText() + "'\n" +
		", output = '" + pageExec.outArea.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))

	_, err := database.Exec(query)
	check(err)
}

func setExec() {
	query := `select exec
					, output
				from obj` +
		` where id = ` + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))

	obj, err := database.Query(query)
	check(err)

	obj.Next()
	var exec, output sql.NullString
	err = obj.Scan(&exec, &output)
	pageExec.execArea.SetText(exec.String, true)
	pageExec.outArea.SetText(output.String, true)
	obj.Close()
}

func (pageExec *pageExecType) show() {
	setExec()
	pageProTree.Pages.SwitchToPage("exec")
	app.SetFocus(pageExec.Flex)
}
