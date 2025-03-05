package main

import (
	"database/sql"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
)

type pageSrcType struct {
	lSrc   *tview.List
	mPosId map[int]int
	*tview.Flex
}

var pageSrc pageSrcType

func (pageSrc *pageSrcType) build() {

	pageSrc.lSrc = tview.NewList()

	pageSrc.lSrc.SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pageSrc.lSrc.ShowSecondaryText(true).
		SetBorderPadding(1, 1, 1, 1)

	pageSrc.lSrc.SetSelectedBackgroundColor(tcell.ColorOrange)

	pageSrc.lSrc.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		pageSrcDesc.descArea.SetText(s2, true)
	})

	pageSrc.lSrc.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyDelete {
			curPos := pageSrc.lSrc.GetCurrentItem()
			delSrc()
			pageSrc.lSrc.Clear()
			setFileSrc(pageProTree.trPro.GetCurrentNode().GetReference().(int))
			if pageSrc.lSrc.GetItemCount() > curPos {
				pageSrc.lSrc.SetCurrentItem(curPos)
			}
		}

		return event
	})

	pageSrc.mPosId = make(map[int]int)

	pageSrc.Flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pageSrc.lSrc, 0, 1, true)

	pageSrc.Flex.SetTitle("F5/F6").
		SetTitleAlign(tview.AlignLeft)

	pageProTree.Pages.AddPage("src", pageSrc.lSrc, true, true)
}

func setFileSrc(idFile int) {
	log.Println("-------------------------------")
	log.Println("setFileSrc")
	log.Println("--------------------")

	query := `select id
				   , line
				   , comment
				from src
			   where id_file = ` + strconv.Itoa(idFile) +
		` order by num asc`

	log.Println(query)

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
		pageSrc.lSrc.AddItem(line.String, comment.String, rune(0), func() {})

	}

	lines.Close()

	log.Println("-------------------------------")
}

func delSrc() {
	log.Println("delSrc")
	query := `DELETE FROM src
			  WHERE id = ` + strconv.Itoa(pageSrc.mPosId[pageSrc.lSrc.GetCurrentItem()])

	log.Println(query)
	_, err := database.Exec(query)
	check(err)
}
