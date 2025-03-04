package main

import (
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
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
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(2, 1, 2, 1).
		SetBackgroundColor(tcell.ColorBlue)

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

			pageExec.execArea.SetText(clipBoardContent, true)
		}
		return event
	})

	pageExec.outArea = tview.NewTextArea()

	pageExec.outArea.
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(2, 2, 1, 1).
		SetBackgroundColor(tcell.ColorBlue)

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

			pageExec.outArea.SetText(clipBoardContent, true)
		}
		return event
	})

	pageExec.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageExec.execArea, 0, 1, true).
		AddItem(pageExec.outArea, 0, 1, true)

	pageExec.Flex.SetTitle("F6").
		SetTitleAlign(tview.AlignLeft)

	pageProTree.Pages.AddPage("exec", pageExec.Flex, true, false)
}

func saveExec() {
	log.Println("-------------------------------")
	log.Println("saveExec")
	log.Println("---------------------")

	query := "UPDATE obj" + "\n" +
		"SET exec = '" + pageExec.execArea.GetText() + "'\n" +
		", output = '" + pageExec.outArea.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))

	log.Println(query)

	_, err := database.Exec(query)
	check(err)

	log.Println("-------------------------------")

}
