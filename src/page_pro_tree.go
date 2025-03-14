package main

import (
	"database/sql"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
	"strings"
)

type pageProTreeType struct {
	trPro       *tview.TreeView
	rootPro     *tview.TreeNode
	descArea    *tview.TextArea
	curFolderID int
	bDesc       bool
	bName       bool
	flTree      *tview.Flex
	*tview.Flex
	*tview.Pages
}

var pageProTree pageProTreeType

func (pageProTree *pageProTreeType) build() {

	pageSrc.bDesc = false
	pageSrc.bName = false

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

	pageProTree.trPro.SetTitle("F3/F4").
		SetTitleAlign(tview.AlignLeft)

	pageProTree.trPro.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		log.Println(event.Key())
		if event.Key() == tcell.KeyDelete {
			log.Println("DELETE")
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
	pageProTree.descArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("comment").
		SetTitleAlign(tview.AlignLeft)

	pageProTree.descArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			removeObjDesc()
			return nil
		}
		return event
	})

	pageProTree.flTree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyCtrlQ {
			log.Println("pageProTree.flTree.SetInputCapture")
			log.Println("KeyCtrlQ")
			log.Println(pageProTree.descArea.HasFocus())
			log.Println(pageProTree.bDesc)

			if pageProTree.descArea.HasFocus() == false {
				if pageProTree.bDesc == false {
					pageProTree.bDesc = true
					pageProTree.flTree.AddItem(pageProTree.descArea, 0, 1, false)
					app.SetFocus(pageProTree.descArea)
				}
				app.SetFocus(pageProTree.descArea)
			} else {
				if pageProTree.bDesc {
					saveObjDescCtrl()
				} else {
					if pageSrc.bName {
						saveObjNameCtrl()
					}
					pageProTree.bDesc = true
					pageProTree.descArea.SetTitle("comment")

					pageProTree.flTree.AddItem(pageProTree.descArea, 0, 1, false)
					app.SetFocus(pageProTree.descArea)
				}
			}
		}
		if event.Key() == tcell.KeyCtrlW {
			if pageSrc.bName {
				saveObjNameCtrl()
			} else {
				if pageSrc.bDesc {
					saveObjDescCtrl()
				}
				pageProTree.bName = true
				pageProTree.descArea.SetTitle("name")
				setSrcLine()
				pageProTree.flTree.AddItem(pageProTree.descArea, 0, 1, true)
				app.SetFocus(pageProTree.descArea)
			}

		}

		return event
	})

	pagePro.Pages.AddPage("proTree", pageProTree.Flex, true, true)

}

// fulfill project tree
func setTreePro(pos int) {
	log.Println("-------------------------------")
	log.Println("setTreePro")

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

		log.Println(id)
		log.Println(objName.String)
		log.Println(objName.String)

		objNode := tview.NewTreeNode(objName.String).
			SetReference(int(id.Int64)).
			SetSelectable(true).
			SetColor(tcell.ColorGreen)

		objNode.SetSelectedFunc(func() {
			pageSrc.lSrc.Clear()
			setObjDesc()
			setFileSrc(int(id.Int64))
			setObjExec(int(id.Int64))
			removeSrcDesc()
			if pageProTree.bName == false {
				showObjDesc()
			}
			pageProTree.Pages.SwitchToPage("src")
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

	log.Println("-------------------------------")
}

// cycle fulfill project tree
func setTreeFolderPro(node *tview.TreeNode) {
	log.Println("-------------------------------")
	log.Println("setTreeFolderPro")
	log.Println("--------------------")

	queryObj := `select id
					   , name
					   , object_type
					from obj
				   where id_parent = ` + strconv.Itoa(pageProTree.curFolderID) +
		` order by object_type asc`

	log.Println(queryObj)

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
			objNode.SetSelectedFunc(func() {
				pageSrc.lSrc.Clear()
				setFileSrc(int(id.Int64))
			})
		default:
			objNode.SetColor(tcell.ColorRed)
		}

		node.AddChild(objNode)
	}

	objects.Close()

	log.Println("-------------------------------")
}

func setObjExec(id int) {
	log.Println("-------------------------------")
	log.Println("setObjExec")
	log.Println("--------------------")
	query := `select exec
				   , output
				from obj
			   where id = ` + strconv.Itoa(id)

	log.Println(query)

	obj := database.QueryRow(query)
	check(obj.Err())

	var exec, output sql.NullString
	log.Println(output)
	err := obj.Scan(&exec, &output)
	check(err)

	pageExec.execArea.SetText(exec.String, true)
	pageExec.outArea.SetText(output.String, true)

}

func delObj(idObj int) {
	log.Println("-------------------------------")
	log.Println("delObj")
	log.Println("--------------------")

	queryObj := `DELETE FROM obj
			    WHERE id = ` + strconv.Itoa(idObj)

	log.Println(queryObj)
	_, err := database.Exec(queryObj)
	check(err)

	querySrc := `DELETE FROM src
			    WHERE id_file = ` + strconv.Itoa(idObj)

	log.Println(querySrc)
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

func saveObj() {
	log.Println("-------------------------------")
	log.Println("saveObj")
	log.Println("---------------------")

	var query string
	if pageProTree.bDesc {
		query = "UPDATE obj" + "\n" +
			"SET comment = '" + pageProTree.descArea.GetText() + "'\n" +
			"WHERE id = " + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))
	} else if pageProTree.bName {
		query = "UPDATE obj" + "\n" +
			"SET name = '" + strings.TrimSpace(pageProTree.descArea.GetText()) + "'\n" +
			"WHERE id = " + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))
	}
	log.Println(query)

	_, err := database.Exec(query)
	check(err)

	log.Println("-------------------------------")

}

func saveObjDescCtrl() {
	saveObj()
	pageProTree.bDesc = false
	pageProTree.flTree.RemoveItem(pageProTree.descArea)
	app.SetFocus(pageProTree.trPro)
}

func saveObjNameCtrl() {
	saveObj()
	pageProTree.bName = false
	pageProTree.flTree.RemoveItem(pageProTree.descArea)
}

func removeObjDesc() {
	if pageProTree.bName {
		saveObjNameCtrl()
		pageProTree.bName = false
		return
	}
	if pageProTree.bDesc {
		saveObjDescCtrl()
		pageProTree.bDesc = false
		return
	}
}

func showObjDesc() {
	log.Println("-------------------------------")
	log.Println("showObjDesc")
	query := `select comment
				from obj` +
		` where id = ` + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))

	log.Println(query)

	obj, err := database.Query(query)
	check(err)

	obj.Next()
	var comment sql.NullString
	err = obj.Scan(&comment)
	if len(comment.String) > 0 {
		pageProTree.descArea.SetText(comment.String, true)
		pageProTree.descArea.SetTitle("comment")
		pageProTree.flTree.AddItem(pageProTree.descArea, 0, 1, false)
		pageProTree.bDesc = false
	} else {
		pageProTree.flTree.RemoveItem(pageProTree.descArea)
		pageProTree.bDesc = false
	}

	obj.Close()

	log.Println("-------------------------------")

}
