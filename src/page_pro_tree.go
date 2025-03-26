package main

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageProTreeType struct {
	trPro       *tview.TreeView
	rootPro     *tview.TreeNode
	descArea    *tview.TextArea
	nameArea    *tview.TextArea
	curFolderID int
	flTree      *tview.Flex
	*tview.Flex
	*tview.Pages
}

var pageProTree pageProTreeType

func (pageProTree *pageProTreeType) build() {

	pageProTree.Pages = tview.NewPages()
	pageSrc.build()
	pageObjDesc.build()
	pageExec.build()

	pageProTree.curFolderID = 0

	pageProTree.rootPro = tview.NewTreeNode(".").
		SetColor(tcell.ColorBlue)

	pageProTree.trPro = tview.NewTreeView().
		SetCurrentNode(pageProTree.rootPro)

	pageProTree.trPro.SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pageProTree.trPro.SetRoot(pageProTree.rootPro)

	pageProTree.trPro.SetTopLevel(1)

	pageProTree.trPro.SetTitle("F3").
		SetTitleAlign(tview.AlignLeft)

	pageProTree.trPro.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDelete {
			delObj(pageProTree.trPro.GetCurrentNode().GetReference().(int))
			pageProTree.rootPro.ClearChildren()
			setTreePro(pagePro.lPro.GetCurrentItem())
			setProComment()
			pageObjDesc.descArea.SetText("", true)
		}

		return event
	})

	pageProTree.rootPro.SetExpanded(true)

	pageObjDesc.Flex.SetFocusFunc(func() {
		app.SetFocus(pageObjDesc.descArea)
	})
	pageProTree.flTree = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageProTree.trPro, 0, 10, true)

	pageProTree.Flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pageProTree.flTree, 0, 4, true).
		AddItem(pageProTree.Pages, 0, 7, false)

	pageProTree.descArea = tview.NewTextArea()
	pageProTree.descArea.SetBorderColor(tcell.ColorBlue)
	pageProTree.descArea.SetBorderPadding(1, 1, 1, 1)
	pageProTree.descArea.SetDisabled(true)

	pageProTree.descArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("comment").
		SetTitleAlign(tview.AlignLeft)

	pageProTree.descArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			saveObjDescCtrl()
			pageProTree.flTree.RemoveItem(pageProTree.descArea)
			app.SetFocus(pageProTree.flTree)
			return nil
		}
		return event
	})

	pageProTree.nameArea = tview.NewTextArea()
	pageProTree.nameArea.SetBorderColor(tcell.ColorBlue)
	pageProTree.nameArea.SetBorderPadding(1, 1, 1, 1)
	pageProTree.nameArea.SetDisabled(true)

	pageProTree.nameArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("name").
		SetTitleAlign(tview.AlignLeft)

	pageProTree.nameArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			saveObjNameCtrl()
			pageProTree.trPro.GetCurrentNode().SetText(strings.TrimSpace(pageProTree.nameArea.GetText()))
			pageProTree.flTree.RemoveItem(pageProTree.nameArea)
			app.SetFocus(pageProTree.flTree)
			return nil
		}
		return event
	})

	pageProTree.flTree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyCtrlQ {
			if pageProTree.descArea.GetDisabled() == true {
				if pageProTree.nameArea.GetDisabled() == false {
					hideObjName()
				}
				pageProTree.descArea.SetTitle("comment")
				pageProTree.flTree.AddItem(pageProTree.descArea, 0, 3, false)
				pageProTree.descArea.SetText(getObjDesc(), true)
				app.SetFocus(pageProTree.descArea)
				pageProTree.descArea.SetDisabled(false)
			} else {
				hideObjDesc()
			}

		}
		if event.Key() == tcell.KeyCtrlW {
			if pageProTree.nameArea.GetDisabled() == true {
				if pageProTree.nameArea.GetDisabled() == true {
					hideObjDesc()
				}
				pageProTree.nameArea.SetTitle("name")
				pageProTree.flTree.AddItem(pageProTree.nameArea, 0, 1, false)
				pageProTree.nameArea.SetText(pageProTree.trPro.GetCurrentNode().GetText()+" <mask>", true)
				app.SetFocus(pageProTree.nameArea)
				pageProTree.nameArea.SetDisabled(false)
			} else {
				hideObjName()
			}
		}

		return event
	})

	pagePro.Pages.AddPage("proTree", pageProTree.Flex, true, true)

}

func setTreePro(pos int) {

	queryObject := `select id
						   , name
						   , object_type
						from obj
					   where id_parent is null
						 and id_prj = ` + strconv.Itoa(pagePro.mPosId[pos]) +
		` order by object_type asc`

	log.Println(queryObject)

	objects, err := database.Query(queryObject)
	check(err)

	for objects.Next() {
		var id sql.NullInt64
		var objName sql.NullString
		var objType sql.NullInt16
		err := objects.Scan(&id, &objName, &objType)
		check(err)

		objNode := tview.NewTreeNode(objName.String).
			SetReference(int(id.Int64)).
			SetSelectable(true).
			SetColor(tcell.ColorGreen)

		objNode.SetSelectedFunc(func() {
			objNodeSelectAction(int(id.Int64))
		})

		switch int(objType.Int16) {
		case 0:
			objNode.SetColor(tcell.ColorOrange)
			pageProTree.curFolderID = int(id.Int64)
			setTreeFolderPro(objNode)
		case 1:
			objNode.SetColor(tcell.ColorGrey)
		default:
			objNode.SetColor(tcell.ColorRed)
		}

		pageProTree.rootPro.AddChild(objNode)
	}

	objects.Close()

	if pageSrc.lSrc.GetItemCount() > 0 {
		pageSrc.lSrc.SetCurrentItem(0)
	}
}

func setTreeFolderPro(node *tview.TreeNode) {

	queryObj := `select id
					   , name
					   , object_type
					from obj
				   where id_parent = ` + strconv.Itoa(pageProTree.curFolderID) +
		` order by object_type asc`

	objects, err := database.Query(queryObj)
	check(err)

	for objects.Next() {
		var id sql.NullInt64
		var objName sql.NullString
		var objType sql.NullInt16
		err := objects.Scan(&id, &objName, &objType)
		check(err)

		objNode := tview.NewTreeNode(objName.String).
			SetReference(int(id.Int64)).
			SetSelectable(true).
			SetColor(tcell.ColorGreen)

		switch int(objType.Int16) {
		case 0: // folder
			objNode.SetColor(tcell.ColorOrange)
			pageProTree.curFolderID = int(id.Int64)
			setTreeFolderPro(objNode)
		case 1: // file
			objNode.SetColor(tcell.ColorGrey)
			objNode.SetSelectedFunc(func() { objNodeSelectAction(int(id.Int64)) })
		default:
			objNode.SetColor(tcell.ColorRed)
		}

		node.AddChild(objNode)
	}

	objects.Close()
}

func setObjExec(id int) {
	query := `select exec
				from obj
			   where id = ` + strconv.Itoa(id)

	obj := database.QueryRow(query)
	check(obj.Err())

	var exec sql.NullString
	err := obj.Scan(&exec)
	check(err)

	pageExec.execArea.SetText(exec.String, true)

}

func delObj(idObj int) {

	queryObj := `DELETE FROM obj
			    WHERE id = ` + strconv.Itoa(idObj)

	_, err := database.Exec(queryObj)
	check(err)

	querySrc := `DELETE FROM src
			    WHERE id_file = ` + strconv.Itoa(idObj)

	_, err = database.Exec(querySrc)
	check(err)

	// TODO
	// bad decision but I don't know how to do it better for now
	// cause of lock database in cycle
	for i := 0; i < 10; i++ {
		delAbsParent()
	}

}

func (pageProTree *pageProTreeType) show() {
	pagePro.Pages.SwitchToPage("proTree")
	app.SetFocus(pageProTree.trPro)
}

func saveObjDescCtrl() {

	var query string
	query = "UPDATE obj" + "\n" +
		"SET comment = '" + pageProTree.descArea.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))

	_, err := database.Exec(query)
	check(err)
}

func saveObjNameCtrl() {
	var query string
	query = "UPDATE obj" + "\n" +
		"SET name = '" + pageProTree.nameArea.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))

	_, err := database.Exec(query)
	check(err)
}

func getObjDesc() string {
	query := `select comment
				from obj` +
		` where id = ` + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))

	obj, err := database.Query(query)
	check(err)

	obj.Next()
	var comment sql.NullString
	err = obj.Scan(&comment)

	obj.Close()

	return comment.String
}

func showObjDesc() {

	query := `select comment
				from obj` +
		` where id = ` + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))

	obj, err := database.Query(query)
	check(err)

	obj.Next()
	var comment sql.NullString
	err = obj.Scan(&comment)

	obj.Close()
}

func hideObjDesc() {
	saveObjDescCtrl()
	pageProTree.descArea.SetDisabled(true)
	pageProTree.flTree.RemoveItem(pageProTree.descArea)
	app.SetFocus(pageProTree.flTree)
}

func hideObjName() {
	saveObjNameCtrl()
	pageProTree.trPro.GetCurrentNode().SetText(strings.TrimSpace(pageProTree.nameArea.GetText()))
	pageProTree.nameArea.SetDisabled(true)
	pageProTree.flTree.RemoveItem(pageProTree.nameArea)
	app.SetFocus(pageProTree.flTree)
}

func objNodeSelectAction(id int) {
	if pageProTree.descArea.GetDisabled() == false {
		pageProTree.descArea.SetDisabled(true)
		pageProTree.flTree.RemoveItem(pageProTree.descArea)
	}
	pageSrc.lSrc.Clear()
	setObjDesc()
	setFileSrc(id)
	setObjExec(id)
	removeSrcDesc()
	showObjDesc()
	pageProTree.Pages.SwitchToPage("src")
}
