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
	nameArea *tview.TextArea
	mPosId   map[int]int
	curPos   int
	*tview.Flex
}

var pageSrc pageSrcType

func (pageSrc *pageSrcType) build() {
	pageSrc.mPosId = make(map[int]int)

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
	pageSrc.descArea.SetBackgroundColor(tcell.ColorDarkSlateGrey)

	pageSrc.descArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			hideSrcDesc()
			return nil
		}
		return event
	})

	pageSrc.nameArea = tview.NewTextArea()
	pageSrc.nameArea.SetBorderColor(tcell.ColorBlue)
	pageSrc.nameArea.SetBorderPadding(1, 1, 1, 1)
	pageSrc.nameArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("name").
		SetTitleAlign(tview.AlignLeft)

	pageSrc.nameArea.SetBackgroundColor(tcell.ColorDarkSlateGrey)

	pageSrc.nameArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			hideSrcName()
			return nil
		}
		return event
	})

	pageSrc.lSrc.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		pageSrc.nameArea.SetText(s, true)
		pageSrc.descArea.SetText(s2, true)
		pageSrc.curPos = pageSrc.lSrc.GetCurrentItem()
	})

	pageSrc.lSrc.SetBackgroundColor(tcell.ColorDarkSlateGrey)

	pageSrc.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageSrc.lSrc, 0, 10, true)

	pageSrc.Flex.SetBackgroundColor(tcell.ColorDarkSlateGrey)
	pageSrc.SetBackgroundColor(tcell.ColorDarkSlateGrey)

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
			if pageSrc.descArea.GetDisabled() == true {
				if pageSrc.nameArea.GetDisabled() == false {
					hideSrcName()
				}
				pageSrc.Flex.AddItem(pageSrc.descArea, 0, 1, false)
				pageSrc.descArea.SetText(getSrcDesc(), true)
				app.SetFocus(pageSrc.descArea)
				pageSrc.descArea.SetDisabled(false)
				saveSrcDesc()
			} else {
				hideSrcDesc()
			}
		}
		if event.Key() == tcell.KeyCtrlW {
			if pageSrc.nameArea.GetDisabled() == true {
				if pageSrc.descArea.GetDisabled() == false {
					hideSrcDesc()
				}
				pageSrc.Flex.AddItem(pageSrc.nameArea, 0, 1, false)
				srcName, _ := pageSrc.lSrc.GetItemText(pageSrc.lSrc.GetCurrentItem())
				//log.Println(srcName)
				srcName = srcName + "<mask>" + "\n"
				pageSrc.nameArea.SetText(srcName+" ", true)
				app.SetFocus(pageSrc.nameArea)
				pageSrc.nameArea.SetDisabled(false)
			} else {
				hideSrcName()
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
	pageProTree.Pages.SwitchToPage("src")
	pageSrc.lSrc.Clear()
	setFileSrc(pageProTree.trPro.GetCurrentNode().GetReference().(int))
	app.SetFocus(pageSrc.Flex)
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
	pageSrc.descArea.SetText(line.String, true)
	src.Close()
}

func saveSrcDesc() {
	var query string
	query = "UPDATE src" + "\n" +
		"SET comment = '" + pageSrc.descArea.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pageSrc.mPosId[pageSrc.curPos])

	_, err := database.Exec(query)
	check(err)
}

func saveSrcName() {
	var query, val string
	val = strings.TrimRight(pageSrc.nameArea.GetText(), " ")
	if len(val) > 0 {
		val = val + " "
	}
	query = "UPDATE src" + "\n" +
		"SET line = '" + val + "'\n" +
		"WHERE id = " + strconv.Itoa(pageSrc.mPosId[pageSrc.curPos])

	_, err := database.Exec(query)
	check(err)
}

func getSrcDesc() string {
	query := `select comment
				from src` +
		` where id = ` + strconv.Itoa(pageSrc.mPosId[pageSrc.lSrc.GetCurrentItem()])

	pro, err := database.Query(query)
	check(err)

	pro.Next()
	var comment sql.NullString
	err = pro.Scan(&comment)

	pro.Close()

	return comment.String
}

func hideSrcDesc() {
	pageSrc.descArea.SetDisabled(true)
	curPos := pageSrc.lSrc.GetCurrentItem()
	saveSrcDesc()
	pageSrc.Flex.RemoveItem(pageSrc.descArea)
	pageSrc.show()
	pageSrc.lSrc.SetCurrentItem(curPos)
	app.SetFocus(pageSrc.lSrc)
}

func hideSrcName() {
	pageSrc.nameArea.SetDisabled(true)
	curPos := pageSrc.lSrc.GetCurrentItem()
	saveSrcName()
	pageSrc.Flex.RemoveItem(pageSrc.nameArea)
	pageSrc.show()
	pageSrc.lSrc.SetCurrentItem(curPos)
	app.SetFocus(pageSrc.lSrc)
}
