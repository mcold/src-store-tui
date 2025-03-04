package main

import (
	"database/sql"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
)

type pageProTreeType struct {
	trPro       *tview.TreeView
	rootPro     *tview.TreeNode
	curFolderID int
	*tview.Flex
	*tview.Pages
}

var pageProTree pageProTreeType

func (pageProTree *pageProTreeType) build() {

	pageProTree.Pages = tview.NewPages()
	pageSrc.build()
	pageObjDesc.build()
	pageSrcDesc.build()
	pageExec.build()

	pageProTree.curFolderID = 0

	pageProTree.rootPro = tview.NewTreeNode(".").
		SetColor(tcell.ColorBlue)

	pageProTree.trPro = tview.NewTreeView().
		SetCurrentNode(pageProTree.rootPro)

	pageProTree.trPro.SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pageProTree.trPro.SetRoot(pageProTree.rootPro)

	pageProTree.rootPro.SetExpanded(true)

	pageObjDesc.Flex.SetFocusFunc(func() {
		app.SetFocus(pageObjDesc.descArea)
	})

	pageSrcDesc.Flex.SetFocusFunc(func() {
		app.SetFocus(pageSrcDesc.descArea)
	})

	pageProTree.Flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pageProTree.trPro, 0, 1, true).
		AddItem(pageProTree.Pages, 0, 3, false)

	//pageProTree.Flex.SetBorder(true).SetBorderColor(tcell.ColorBlue).

	pageProTree.Flex.SetTitle("F3/F4").
		SetTitleAlign(tview.AlignLeft)

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
			setFileSrc(int(id.Int64))
			setObjExec(int(id.Int64))
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
		//pageProTree.mObjType[strconv.Itoa(int(id.Int64))+fileName.String] = "file"
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

	obj := database.QueryRow(query)
	check(obj.Err())

	var exec, output sql.NullString
	log.Println(output)
	err := obj.Scan(&exec, &output)
	check(err)

	pageExec.execArea.SetText(exec.String, true)
	pageExec.outArea.SetText(output.String, true)

}
