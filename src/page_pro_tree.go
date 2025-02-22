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

	pageProTree.curFolderID = 0

	pageProTree.rootPro = tview.NewTreeNode(".").
		SetColor(tcell.ColorBlue)

	pageProTree.trPro = tview.NewTreeView().
		SetCurrentNode(pageProTree.rootPro)

	pageProTree.trPro.SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pageProTree.trPro.SetRoot(pageProTree.rootPro)

	pageProTree.rootPro.SetExpanded(true)

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
	log.Println("--------- files -----------")

	// files
	queryFile := `select id
				   , name
				   , comment
				from file
			   where id_folder is null
				 and id_prj = ` + strconv.Itoa(pagePro.mPosId[pos])

	log.Println(queryFile)

	files, err := database.Query(queryFile)
	check(err)

	for files.Next() {
		id := -1
		var fileName, fileComment sql.NullString
		err := files.Scan(&id, &fileName, &fileComment)
		check(err)

		log.Println(id)
		log.Println(fileName.String)
		log.Println(fileComment.String)

		fileNode := tview.NewTreeNode(fileName.String).
			SetReference(id).
			SetSelectable(true).
			SetColor(tcell.ColorGreen)

		fileNode.SetSelectedFunc(func() {
			pageSrc.lSrc.Clear()
			setFileSrc(id)
			pageProTree.Pages.SwitchToPage("src")
		})

		pageProTree.rootPro.AddChild(fileNode)
	}

	// folders
	log.Println("--------- folders -----------")
	log.Println("current item in list: " + strconv.Itoa(pagePro.lPro.GetCurrentItem()))
	log.Println("mpos[pagePro.lPro.GetCurrentItem()]: " + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()]))
	queryFolder := `select id
				   , name
				   , comment
				from folder
			   where id_parent is null
                 and id_prj = ` + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])

	log.Println(queryFolder)

	folders, err := database.Query(queryFolder)
	check(err)

	for folders.Next() {
		id := -1
		var folderName, folderComment sql.NullString
		err := folders.Scan(&id, &folderName, &folderComment)
		check(err)

		log.Println(id)
		log.Println(folderName.String)
		log.Println(folderComment.String)

		folderNode := tview.NewTreeNode(folderName.String).SetReference(id).SetSelectable(true).
			SetColor(tcell.ColorGreen)

		pageProTree.curFolderID = id
		setTreeFolderPro(folderNode)
		pageProTree.rootPro.AddChild(folderNode)
	}
	log.Println("-------------------------------")
}

// cycle fulfill project tree
func setTreeFolderPro(node *tview.TreeNode) {
	log.Println("-------------------------------")
	log.Println("setTreeFolderPro")
	log.Println("--------------------")

	if pageProTree.curFolderID == 0 {
		return
	}

	// files
	log.Println("--------- files -----------")
	queryFile := `select id
					   , name
					   , comment
					from file
				   where id_folder = ` + strconv.Itoa(pageProTree.curFolderID)

	log.Println(queryFile)

	files, err := database.Query(queryFile)
	check(err)

	for files.Next() {
		id := -1
		var fileName, fileComment sql.NullString
		err := files.Scan(&id, &fileName, &fileComment)
		check(err)

		fileNode := tview.NewTreeNode(fileName.String).
			SetReference(id).
			SetSelectable(true).
			SetColor(tcell.ColorGreen)

		fileNode.SetSelectedFunc(func() {
			pageSrc.lSrc.Clear()
			setFileSrc(id)
		})

		node.AddChild(fileNode)
	}

	log.Println("--------- folders -----------")
	// folders
	queryFolder := `select id
						   , name
						   , comment
						from folder
					   where id_parent = ` + strconv.Itoa(pageProTree.curFolderID)

	log.Println(queryFolder)

	folders, err := database.Query(queryFolder)
	check(err)

	for folders.Next() {
		id := -1
		var folderName, folderComment sql.NullString
		err := folders.Scan(&id, &folderName, &folderComment)
		check(err)

		folderNode := tview.NewTreeNode(folderName.String).SetReference(id).SetSelectable(true).
			SetColor(tcell.ColorGreen)

		pageProTree.curFolderID = id
		setTreeFolderPro(folderNode)
		node.AddChild(folderNode)
	}
	log.Println("-------------------------------")
}
