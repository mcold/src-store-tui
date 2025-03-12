package main

import (
	"database/sql"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
)

type pageSrcType struct {
	lSrc     *tview.List
	descArea *tview.TextArea
	bDesc    bool
	mPosId   map[int]int
	*tview.Flex
}

var pageSrc pageSrcType

func (pageSrc *pageSrcType) build() {

	pageSrc.bDesc = false

	pageSrc.lSrc = tview.NewList()

	pageSrc.lSrc.SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pageSrc.lSrc.ShowSecondaryText(true).
		SetBorderPadding(1, 1, 1, 1)

	pageSrc.lSrc.SetSelectedBackgroundColor(tcell.ColorOrange)

	pageSrc.lSrc.SetTitle("F5/F6").
		SetTitleAlign(tview.AlignLeft)

	pageSrc.descArea = tview.NewTextArea()
	pageSrc.descArea.SetBorderColor(tcell.ColorBlue)
	pageSrc.descArea.SetBorderPadding(1, 1, 1, 1)
	pageSrc.descArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("comment").
		SetTitleAlign(tview.AlignLeft)

	pageSrc.lSrc.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		pageSrc.descArea.SetText(s2, true)
	})

	pageSrc.mPosId = make(map[int]int)

	pageSrc.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageSrc.lSrc, 0, 10, true)

	pageSrc.Flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyDelete {
			curPos := pageSrc.lSrc.GetCurrentItem()
			delSrc()
			pageSrc.lSrc.Clear()
			setFileSrc(pageProTree.trPro.GetCurrentNode().GetReference().(int))
			if pageSrc.lSrc.GetItemCount() > curPos {
				pageSrc.lSrc.SetCurrentItem(curPos)
			}
		}
		if event.Key() == tcell.KeyCtrlQ {
			if pageSrc.bDesc {
				pageSrc.bDesc = false
				curPos := pageSrc.lSrc.GetCurrentItem()
				saveSrc()
				pageSrc.Flex.RemoveItem(pageSrc.descArea)
				pageSrc.show()
				pageSrc.lSrc.SetCurrentItem(curPos)
			} else {
				pageSrc.bDesc = true
				setSrc()
				pageSrc.Flex.AddItem(pageSrc.descArea, 0, 1, false)
				app.SetFocus(pageSrc.descArea)
			}

		}

		//if event.Key() == tcell.KeyCtrlW {
		//	if pageSrc.bName {
		//		if pageSrc.bDesc {
		//			pageSrc.bDesc = false
		//			curPos := pageSrc.lSrc.GetCurrentItem()
		//			saveSrc()
		//			pageSrc.Flex.RemoveItem(pageSrc.descArea)
		//			pageSrc.show()
		//			pageSrc.lSrc.SetCurrentItem(curPos)
		//		}
		//		pageSrc.bName = false
		//		curPos := pageSrc.lSrc.GetCurrentItem()
		//		saveSrc()
		//		pageSrc.Flex.RemoveItem(pageSrc.nameArea)
		//		pageSrc.show()
		//		pageSrc.lSrc.SetCurrentItem(curPos)
		//	} else {
		//		pageSrc.bName = true
		//		//sName := pageSrc.nameArea.GetText()
		//		pageSrc.Flex.AddItem(pageSrc.nameArea, 0, 1, false)
		//		//pageSrc.nameArea.SetText(sName, true)
		//		setSrc()
		//		app.SetFocus(pageSrc.nameArea)
		//	}
		//
		//}

		return event
	})

	pageProTree.Pages.AddPage("src", pageSrc.Flex, true, true)
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

func (pageSrc *pageSrcType) show() {
	pageSrc.bDesc = false
	pageProTree.Pages.SwitchToPage("src")
	pageSrc.lSrc.Clear()
	setFileSrc(pageProTree.trPro.GetCurrentNode().GetReference().(int))
	app.SetFocus(pageSrc.Flex)
}

func saveSrc() {
	log.Println("-------------------------------")
	log.Println("saveSrc")
	log.Println("---------------------")

	query := "UPDATE src" + "\n" +
		"SET comment = '" + pageSrc.descArea.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pageSrc.mPosId[pageSrc.lSrc.GetCurrentItem()])

	log.Println(query)

	_, err := database.Exec(query)
	check(err)

	log.Println("-------------------------------")

}

func setSrc() {
	log.Println("-------------------------------")
	log.Println("setObjDesc")
	query := `select line
       				 , comment
				from src` +
		` where id = ` + strconv.Itoa(pageSrc.mPosId[pageSrc.lSrc.GetCurrentItem()])

	src, err := database.Query(query)
	check(err)

	src.Next()
	var line, comment sql.NullString
	err = src.Scan(&line, &comment)
	pageSrc.descArea.SetText(comment.String, true)
	src.Close()

	log.Println("-------------------------------")
}
