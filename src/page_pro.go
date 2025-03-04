package main

import (
	"database/sql"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
)

type pageProType struct {
	lPro   *tview.List
	mPosId map[int]int
	*tview.Pages
}

var pagePro pageProType

func (pagePro *pageProType) build() {

	pagePro.Pages = tview.NewPages()
	pageProTree.build()
	pageProDesc.build()

	pagePro.mPosId = make(map[int]int)

	// list
	pagePro.lPro = tview.NewList()

	pagePro.lPro.SetBorder(true)
	pagePro.lPro.SetBorderColor(tcell.ColorBlue)
	pagePro.lPro.SetBorderColor(tcell.ColorBlue)

	setListPro()

	pagePro.lPro.SetSelectedFunc(func(pos int, _ string, _ string, _ rune) {

		pagePro.Pages.SwitchToPage("proTree")
		pageProTree.rootPro.ClearChildren()

		setTreePro(pos)
		setProComment()

		app.SetFocus(pageProTree.Flex)
	})

	pagePro.lPro.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDelete || event.Key() == tcell.KeyDEL {
			delPro(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])
			reloadProTree()
			return nil
		}
		return event
	})

	flexPro := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pagePro.lPro, 0, 1, true).
		AddItem(pagePro.Pages, 0, 4, true)

	flexPro.SetBorder(true).SetBorderColor(tcell.ColorBlue).
		SetTitle("F2").
		SetTitleAlign(tview.AlignLeft)

	pageMain.pages.AddPage("pro", flexPro, true, true)

}

// fulfill project list
func setListPro() {
	log.Println("-------------------------------")
	log.Println("setListPro")
	query := `select id
				   , name
				   , comment
				from prj
			   order by name`

	pros, err := database.Query(query)
	check(err)

	posNum := 0
	for pros.Next() {
		posNum++
		var id sql.NullInt64
		var name, comment sql.NullString

		err := pros.Scan(&id, &name, &comment)
		check(err)

		pagePro.mPosId[posNum-1] = int(id.Int64)
		pagePro.lPro.AddItem(name.String, comment.String, rune(0), func() {})
	}

	pros.Close()

	log.Println("-------------------------------")
}

func setProComment() {
	log.Println("-------------------------------")
	log.Println("setProComment")
	query := `select comment
				from prj` +
		` where id = ` + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])

	pros, err := database.Query(query)
	check(err)

	pros.Next()
	var comment sql.NullString
	err = pros.Scan(&comment)

	pageProDesc.descArea.SetText(comment.String, true)
	pros.Close()

	log.Println("-------------------------------")
}

func delPro(idPro int) {

	querySrc := `DELETE FROM src
			    WHERE id_prj = ` + strconv.Itoa(idPro)

	log.Println("delPro", querySrc)
	_, err := database.Exec(querySrc)
	check(err)

	queryObj := `DELETE FROM obj
			    WHERE id_prj = ` + strconv.Itoa(idPro)

	log.Println("delPro", queryObj)
	_, err = database.Exec(queryObj)

	queryPro := `DELETE FROM prj
			    WHERE id = ` + strconv.Itoa(idPro)

	log.Println("delPro", queryPro)
	_, err = database.Exec(queryPro)

	check(err)
}

func reloadProTree() {
	pagePro.lPro.Clear()
	pageProTree.rootPro.ClearChildren()
	setListPro()

	if pagePro.lPro.GetItemCount() > 0 {
		pagePro.lPro.SetCurrentItem(0)
	}
}
