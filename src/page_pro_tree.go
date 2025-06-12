package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageProTreeType struct {
	trPro       *tview.TreeView
	rootPro     *tview.TreeNode
	descArea    *tview.TextArea
	nameArea    *tview.TextArea
	exportArea  *tview.TextArea
	importArea  *tview.TextArea
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
			//setTreePro(0)
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
		AddItem(pageProTree.flTree, 0, 2, true).
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

	pageProTree.exportArea = tview.NewTextArea()
	pageProTree.exportArea.SetBorderColor(tcell.ColorBlue)
	pageProTree.exportArea.SetBorderPadding(1, 1, 1, 1)
	pageProTree.exportArea.SetDisabled(true)

	pageProTree.exportArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("export").
		SetTitleAlign(tview.AlignLeft)

	pageProTree.exportArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			hideObjExport()
			return nil
		}
		if event.Rune() == 'v' && event.Modifiers() == tcell.ModAlt {
			clipBoardContent, err := clipboard.ReadAll()
			check(err)

			pageProTree.exportArea.SetText(pageProTree.exportArea.GetText()+clipBoardContent, true)
		}
		return event
	})

	pageProTree.importArea = tview.NewTextArea()
	pageProTree.importArea.SetBorderColor(tcell.ColorBlue)
	pageProTree.importArea.SetBorderPadding(1, 1, 1, 1)
	pageProTree.importArea.SetDisabled(true)

	pageProTree.importArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("import").
		SetTitleAlign(tview.AlignLeft)

	pageProTree.importArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			hideObjImport()
			return nil
		}
		if event.Rune() == 'v' && event.Modifiers() == tcell.ModAlt {
			clipBoardContent, err := clipboard.ReadAll()
			check(err)

			pageProTree.importArea.SetText(pageProTree.importArea.GetText()+clipBoardContent, true)
		}
		return event
	})

	pageProTree.flTree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyCtrlQ {
			if pageProTree.descArea.GetDisabled() == true {
				if pageProTree.nameArea.GetDisabled() == false {
					hideObjName()
				}
				if pageProTree.exportArea.GetDisabled() == false {
					hideObjExport()
				}
				if pageProTree.importArea.GetDisabled() == false {
					hideObjImport()
				}
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
				if pageProTree.exportArea.GetDisabled() == false {
					hideObjExport()
				}
				if pageProTree.importArea.GetDisabled() == false {
					hideObjImport()
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

		if event.Key() == tcell.KeyCtrlS {
			if pageProTree.exportArea.GetDisabled() == true {
				if pageProTree.descArea.GetDisabled() == false {
					hideObjDesc()
				}
				if pageProTree.nameArea.GetDisabled() == false {
					hideObjName()
				}
				if pageProTree.importArea.GetDisabled() == false {
					hideObjImport()
				}
				pageProTree.exportArea.SetTitle("save")
				pageProTree.flTree.AddItem(pageProTree.exportArea, 0, 1, false)
				app.SetFocus(pageProTree.exportArea)
				pageProTree.exportArea.SetDisabled(false)
			} else {
				hideObjExport()
			}

		}

		if event.Key() == tcell.KeyInsert {
			if pageProTree.importArea.GetDisabled() == true {
				if pageProTree.descArea.GetDisabled() == false {
					hideObjDesc()
				}
				if pageProTree.nameArea.GetDisabled() == false {
					hideObjName()
				}
				if pageProTree.exportArea.GetDisabled() == false {
					hideObjExport()
				}
				pageProTree.importArea.SetTitle("import")
				pageProTree.flTree.AddItem(pageProTree.importArea, 0, 1, false)
				app.SetFocus(pageProTree.importArea)
				pageProTree.importArea.SetDisabled(false)
			} else {
				hideObjImport()
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
					   where (id_parent is null or id_parent = 0)
						 and id_prj = ` + strconv.Itoa(pagePro.mPosId[pos]) +
		` order by object_type asc`

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
	showObjDesc()
	pageProTree.Pages.SwitchToPage("src")
}

func hideObjExport() {
	pageProTree.exportArea.SetDisabled(true)

	path := pageProTree.exportArea.GetText()
	if len(strings.TrimSpace(path)) > 0 {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.Mkdir(path, 0777)
			check(err)
		}
		downloadObj(pageProTree.trPro.GetCurrentNode().GetReference().(int), path)
	}

	pageProTree.flTree.RemoveItem(pageProTree.exportArea)
	app.SetFocus(pageProTree.flTree)
}

func hideObjImport() {
	pageProTree.importArea.SetDisabled(true)

	path := pageProTree.importArea.GetText()

	if len(strings.TrimSpace(path)) > 0 {
		_, err := os.Stat(path)
		check(err)

		importObj(path)
	}

	pageProTree.flTree.RemoveItem(pageProTree.importArea)
	app.SetFocus(pageProTree.flTree)
}

func downloadObj(objID int, path string) {

	query := `select name
				   , object_type
		from obj
	  where id = ` + strconv.Itoa(objID)

	objs, err := database.Query(query)
	check(err)

	objs.Next()
	var objName sql.NullString
	var objType sql.NullInt16
	err = objs.Scan(&objName, &objType)
	check(err)

	objs.Close()

	objPath := filepath.Join(path, objName.String)
	switch int(objType.Int16) {
	case 0:

		err := os.Mkdir(objPath, 0777)
		check(err)

		downloadFolderPro(objID, objPath)
	case 1:

		file, err := os.Create(objPath)
		if err != nil {
			panic(err)
		}
		downloadFilePro(objID, objPath)
		defer file.Close()
	}
}
