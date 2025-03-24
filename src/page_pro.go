package main

import (
	"database/sql"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
	"strings"
)

type pageProType struct {
	lPro      *tview.List
	descArea  *tview.TextArea
	nameArea  *tview.TextArea
	mPosId    map[int]int
	flListPro *tview.Flex
	*tview.Pages
}

var pagePro pageProType

func (pagePro *pageProType) build() {

	pagePro.Pages = tview.NewPages()
	pageProTree.build()
	pageProDesc.build()

	pagePro.mPosId = make(map[int]int)

	pagePro.lPro = tview.NewList()

	pagePro.lPro.SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetBorderColor(tcell.ColorBlue).
		SetBorderPadding(1, 1, 1, 1)

	pagePro.lPro.SetTitle("F2").
		SetTitleAlign(tview.AlignLeft)

	setListPro()

	if pagePro.lPro.GetItemCount() > 0 {
		pagePro.lPro.SetCurrentItem(0)
		pagePro.Pages.SwitchToPage("proTree")
		pageProTree.rootPro.ClearChildren()

		setTreePro(0)
		setProComment()
	}

	pagePro.lPro.SetSelectedFunc(func(pos int, _ string, _ string, _ rune) {

		pagePro.Pages.SwitchToPage("proTree")
		pageProTree.rootPro.ClearChildren()

		setTreePro(pos)
		setProComment()

		//app.SetFocus(pageProTree.Flex)
	})

	pagePro.lPro.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDelete || event.Key() == tcell.KeyDEL {
			delPro(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])
			reloadProTree()
			return nil
		}
		return event
	})

	pagePro.flListPro = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pagePro.lPro, 0, 10, true)

	pagePro.descArea = tview.NewTextArea()
	pagePro.descArea.SetBorderColor(tcell.ColorBlue)
	pagePro.descArea.SetBorderPadding(1, 1, 1, 1)
	pagePro.descArea.SetDisabled(true)

	pagePro.descArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("comment").
		SetTitleAlign(tview.AlignLeft)

	pagePro.descArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			hideProDesc()
			return nil
		}
		return event
	})

	pagePro.nameArea = tview.NewTextArea()
	pagePro.nameArea.SetBorderColor(tcell.ColorBlue)
	pagePro.nameArea.SetBorderPadding(1, 1, 1, 1)
	pagePro.nameArea.SetDisabled(true)

	pagePro.nameArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("name").
		SetTitleAlign(tview.AlignLeft)

	pagePro.nameArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			hideProName()
			return nil
		}
		return event
	})

	flexPro := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pagePro.flListPro, 0, 1, true).
		AddItem(pagePro.Pages, 0, 4, true)

	pagePro.flListPro.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyCtrlQ {
			if pagePro.descArea.GetDisabled() == true {
				pagePro.flListPro.AddItem(pagePro.descArea, 0, 1, false)
				pagePro.descArea.SetText(getProDesc(), true)
				app.SetFocus(pagePro.descArea)
				pagePro.descArea.SetDisabled(false)
			} else {
				hideProDesc()
			}

		}
		if event.Key() == tcell.KeyCtrlW {
			if pagePro.nameArea.GetDisabled() == true {
				pagePro.nameArea.SetTitle("name")
				pagePro.flListPro.AddItem(pagePro.nameArea, 0, 1, false)
				itemName, _ := pagePro.lPro.GetItemText(pagePro.lPro.GetCurrentItem())
				pagePro.nameArea.SetText(itemName+" <mask>", true)
				app.SetFocus(pagePro.nameArea)
				pagePro.nameArea.SetDisabled(false)
			} else {
				hideProName()
			}

		}

		return event
	})

	pageMain.pages.AddPage("pro", flexPro, true, true)

}

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

func saveProName() {
	log.Println("-------------------------------")
	log.Println("saveProName")
	log.Println("---------------------")

	query := "UPDATE prj" + "\n" +
		"SET name = '" + strings.TrimSpace(pagePro.nameArea.GetText()) + "'\n" +
		"WHERE id = " + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])

	log.Println(query)

	_, err := database.Exec(query)
	check(err)

	database.Close()
	database.Connect()

	log.Println("-------------------------------")
}

func saveProProDesc() {
	log.Println("-------------------------------")
	log.Println("saveProProDesc")
	log.Println("---------------------")

	query := "UPDATE prj" + "\n" +
		"SET comment = '" + strings.TrimSpace(pagePro.descArea.GetText()) + "'\n" +
		"WHERE id = " + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])

	log.Println(query)

	_, err := database.Exec(query)
	check(err)

	log.Println("-------------------------------")
}

func getProDesc() string {
	query := `select comment
				from prj` +
		` where id = ` + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])

	log.Println(query)

	pro, err := database.Query(query)
	check(err)

	pro.Next()
	var comment sql.NullString
	err = pro.Scan(&comment)

	pro.Close()

	return comment.String
}

func saveProDescCtrl() {
	log.Println("-------------------------------")
	log.Println("saveProDescCtrl")
	log.Println("---------------------")

	var query string
	query = "UPDATE prj" + "\n" +
		"SET comment = '" + pagePro.descArea.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])
	log.Println(query)

	_, err := database.Exec(query)
	check(err)

	database.Close()
	database.Connect()

	log.Println("-------------------------------")
}

func hideProDesc() {
	pagePro.descArea.SetDisabled(true)
	curPos := pagePro.lPro.GetCurrentItem()
	saveProProDesc()
	pagePro.flListPro.RemoveItem(pagePro.descArea)
	pagePro.lPro.Clear()
	setListPro()
	pagePro.lPro.SetCurrentItem(curPos)
	app.SetFocus(pagePro.lPro)
}

func hideProName() {
	pagePro.nameArea.SetDisabled(true)
	curPos := pagePro.lPro.GetCurrentItem()
	saveProName()
	pagePro.flListPro.RemoveItem(pagePro.nameArea)
	pagePro.lPro.Clear()
	setListPro()
	pagePro.lPro.SetCurrentItem(curPos)
	app.SetFocus(pagePro.lPro)
}
