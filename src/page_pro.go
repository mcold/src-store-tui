package main

import (
	"database/sql"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
)

type pageProType struct {
	lPro *tview.List
	//trPro       *tview.TreeView
	//rootPro     *tview.TreeNode
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
		//log.Println(pos)
		//log.Println(pagePro.mPosId[pos])
		//log.Println("-------")
		pageProTree.rootPro.ClearChildren()
		log.Println("---- mPos start ----")
		for i, k := range pagePro.mPosId {
			log.Println(strconv.Itoa(i) + ": " + strconv.Itoa(k))
		}
		log.Println("---- mPos end ----")

		setTreePro(pos)

		app.SetFocus(pageProTree.Flex)
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
	log.Println("--------- files -----------")
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
		id := -1
		var name, comment sql.NullString

		err := pros.Scan(&id, &name, &comment)
		sName, sComment := name.String, comment.String
		check(err)

		pagePro.mPosId[posNum-1] = id
		pagePro.lPro.AddItem(sName, sComment, rune(0), func() {})
	}

	log.Println("-------------------------------")
}
