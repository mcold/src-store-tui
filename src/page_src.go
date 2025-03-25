package main

import (
	"database/sql"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	"strings"
)

type pageSrcType struct {
	lSrc     *tview.List
	descArea *tview.TextArea
	bDesc    bool
	bName    bool
	mPosId   map[int]int
	curPos   int
	*tview.Flex
}

var pageSrc pageSrcType

func (pageSrc *pageSrcType) build() {

	pageSrc.bDesc = false
	pageSrc.bName = false

	pageSrc.lSrc = tview.NewList()

	pageSrc.lSrc.SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pageSrc.lSrc.ShowSecondaryText(true).
		SetBorderPadding(1, 1, 1, 1)

	pageSrc.lSrc.SetSelectedBackgroundColor(tcell.ColorOrange)

	pageSrc.lSrc.SetTitle("F4").
		SetTitleAlign(tview.AlignLeft)

	pageSrc.descArea = tview.NewTextArea()
	pageSrc.descArea.SetBorderColor(tcell.ColorBlue)
	pageSrc.descArea.SetBorderPadding(1, 1, 1, 1)
	pageSrc.descArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("comment").
		SetTitleAlign(tview.AlignLeft)

	pageSrc.descArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			removeSrcDesc()
			return nil
		}
		return event
	})

	pageSrc.lSrc.SetSelectedFunc(func(i int, s string, s2 string, r rune) {

		pageSrc.descArea.SetText(s2, true)
		pageSrc.curPos = pageSrc.lSrc.GetCurrentItem()
	})

	pageSrc.mPosId = make(map[int]int)

	pageSrc.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageSrc.lSrc, 0, 10, true)

	pageSrc.Flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyDelete {
			delSrc()
			pageSrc.lSrc.Clear()
			setFileSrc(pageProTree.trPro.GetCurrentNode().GetReference().(int))
			if pageSrc.lSrc.GetItemCount() > pageSrc.curPos {
				pageSrc.lSrc.SetCurrentItem(pageSrc.curPos)
			}
		}
		if event.Key() == tcell.KeyCtrlQ {
			if pageSrc.bDesc {
				saveSrcDescCtrl()
			} else {
				if pageSrc.bName {
					saveSrcLineCtrl()
				}
				pageSrc.bDesc = true
				pageSrc.descArea.SetTitle("comment")

				_, desc := pageSrc.lSrc.GetItemText(pageSrc.curPos)
				pageSrc.descArea.SetText(strings.TrimSpace(desc)+" ", true)

				pageSrc.Flex.AddItem(pageSrc.descArea, 0, 1, false)
				app.SetFocus(pageSrc.descArea)
			}

		}
		if event.Key() == tcell.KeyCtrlW {
			if pageSrc.bName {
				saveSrcLineCtrl()
			} else {
				if pageSrc.bDesc {
					saveSrcDescCtrl()
				}
				pageSrc.bName = true
				pageSrc.descArea.SetTitle("name")
				setSrcLine()
				pageSrc.Flex.AddItem(pageSrc.descArea, 0, 1, true)
				app.SetFocus(pageSrc.descArea)
			}

		}

		return event
	})

	pageProTree.Pages.AddPage("src", pageSrc.Flex, true, true)
}

func setFileSrc(idFile int) {

	query := `select id
				   , line
				   , comment
				from src
			   where id_file = ` + strconv.Itoa(idFile) +
		` order by num asc`

	lines, err := database.Query(query)
	check(err)

	posNum := 0
	for lines.Next() {
		posNum++
		var id sql.NullInt64

		var line, comment sql.NullString
		err := lines.Scan(&id, &line, &comment)
		check(err)

		pageSrc.mPosId[posNum-1] = int(id.Int64)
		pageSrc.lSrc.AddItem(strings.ReplaceAll(line.String, "\t", "    "), comment.String, rune(0), func() {})

	}

	lines.Close()
}

func delSrc() {
	query := `DELETE FROM src
			  WHERE id = ` + strconv.Itoa(pageSrc.mPosId[pageSrc.curPos])

	_, err := database.Exec(query)
	check(err)
}

func (pageSrc *pageSrcType) show() {
	pageSrc.bDesc = false
	pageSrc.bName = false
	pageProTree.Pages.SwitchToPage("src")
	pageSrc.lSrc.Clear()
	setFileSrc(pageProTree.trPro.GetCurrentNode().GetReference().(int))
	app.SetFocus(pageSrc.Flex)
}

func saveSrc() {
	var query string
	if pageSrc.bDesc {
		query = "UPDATE src" + "\n" +
			"SET comment = '" + pageSrc.descArea.GetText() + "'\n" +
			"WHERE id = " + strconv.Itoa(pageSrc.mPosId[pageSrc.curPos])
	} else if pageSrc.bName {
		query = "UPDATE src" + "\n" +
			"SET line = '" + strings.TrimRight(pageSrc.descArea.GetText(), " ") + "'\n" +
			"WHERE id = " + strconv.Itoa(pageSrc.mPosId[pageSrc.curPos])
	}

	_, err := database.Exec(query)
	check(err)
}

func setSrcLine() {
	query := `select line
				from src` +
		` where id = ` + strconv.Itoa(pageSrc.mPosId[pageSrc.curPos])

	src, err := database.Query(query)
	check(err)

	src.Next()
	var line sql.NullString
	err = src.Scan(&line)
	// TODO: added mask - cause trim last token - I don't know why
	pageSrc.descArea.SetText(line.String+" <mask>\n", true)
	src.Close()
}

func saveSrcDescCtrl() {
	saveSrc()
	pageSrc.bDesc = false
	pageSrc.Flex.RemoveItem(pageSrc.descArea)
	pageSrc.show()
	pageSrc.lSrc.SetCurrentItem(pageSrc.curPos)
}

func saveSrcLineCtrl() {
	saveSrc()
	pageSrc.bName = false
	pageSrc.Flex.RemoveItem(pageSrc.descArea)
	pageSrc.show()
	pageSrc.lSrc.SetCurrentItem(pageSrc.curPos)
}

func removeSrcDesc() {
	if pageSrc.bName {
		saveSrcLineCtrl()
		pageSrc.bName = false
		return
	}
	if pageSrc.bDesc {
		saveSrcDescCtrl()
		pageSrc.bDesc = false
		return
	}
}
