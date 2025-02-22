package main

import (
	"database/sql"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
)

type pageSrcType struct {
	lSrc *tview.List
	*tview.Flex
}

var pageSrc pageSrcType

func (pageSrc *pageSrcType) build() {

	pageSrc.lSrc = tview.NewList()

	pageSrc.lSrc.SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pageSrc.lSrc.ShowSecondaryText(true).
		SetBorderPadding(1, 1, 1, 1)

	pageSrc.Flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pageSrc.lSrc, 0, 1, true)

	//pageSrc.Flex.SetBorder(true).SetBorderColor(tcell.ColorBlue)

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
				from srcc
			   where id_file = ` + strconv.Itoa(idFile) +
		` order by num asc`

	log.Println(query)

	lines, err := database.Query(query)
	check(err)

	for lines.Next() {
		var id sql.NullInt64
		var line, comment sql.NullString

		err := lines.Scan(&id, &line, &comment)
		check(err)

		pageSrc.lSrc.AddItem(line.String, comment.String, rune(0), func() {})

	}

	log.Println("-------------------------------")
}
