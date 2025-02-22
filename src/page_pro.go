package main

import (
	"database/sql"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
)

type pageProType struct {
	lPro        *tview.List
	trPro       *tview.TreeView
	rootPro     *tview.TreeNode
	lSrc        *tview.List
	mPosId      map[int]int
	curFolderID int
}

var pagePro pageProType

func (pagePro *pageProType) build() {

	pagePro.mPosId = make(map[int]int)
	pagePro.curFolderID = 0

	// list
	pagePro.lPro = tview.NewList()

	pagePro.lPro.SetBorder(true)
	pagePro.lPro.SetBorderColor(tcell.ColorBlue)

	setListPro()

	// tree
	pagePro.rootPro = tview.NewTreeNode(".").
		SetColor(tcell.ColorBlue)

	pagePro.trPro = tview.NewTreeView().
		SetCurrentNode(pagePro.rootPro)

	pagePro.trPro.SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pagePro.trPro.SetRoot(pagePro.rootPro)

	//setTreePro()

	pagePro.rootPro.SetExpanded(true)

	pagePro.lSrc = tview.NewList()

	pagePro.lSrc.SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pagePro.lSrc.ShowSecondaryText(true)

	pagePro.lPro.SetSelectedFunc(func(pos int, _ string, _ string, _ rune) {

		//log.Println(pos)
		//log.Println(pagePro.mPosId[pos])
		//log.Println("-------")
		pagePro.rootPro.ClearChildren()
		log.Println("---- mPos start ----")
		for i, k := range pagePro.mPosId {
			log.Println(strconv.Itoa(i) + ": " + strconv.Itoa(k))
		}
		log.Println("---- mPos end ----")

		setTreePro(pos)
	})

	flexPro := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pagePro.lPro, 0, 1, true).
		AddItem(pagePro.trPro, 0, 2, true).
		AddItem(pagePro.lSrc, 0, 5, false)

	flexPro.SetBorder(true).SetBorderColor(tcell.ColorBlue)

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
			pagePro.lSrc.Clear()
			setFileSrc(id)
		})

		pagePro.rootPro.AddChild(fileNode)
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

		pagePro.curFolderID = id
		setTreeFolderPro(folderNode)
		pagePro.rootPro.AddChild(folderNode)
	}
	log.Println("-------------------------------")
}

// cycle fulfill project tree
func setTreeFolderPro(node *tview.TreeNode) {
	log.Println("-------------------------------")
	log.Println("setTreeFolderPro")
	log.Println("--------------------")

	if pagePro.curFolderID == 0 {
		return
	}

	// files
	log.Println("--------- files -----------")
	queryFile := `select id
					   , name
					   , comment
					from file
				   where id_folder = ` + strconv.Itoa(pagePro.curFolderID)

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
			pagePro.lSrc.Clear()
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
					   where id_parent = ` + strconv.Itoa(pagePro.curFolderID)

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

		pagePro.curFolderID = id
		setTreeFolderPro(folderNode)
		node.AddChild(folderNode)
	}
	log.Println("-------------------------------")
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

		pagePro.lSrc.AddItem(line.String, comment.String, rune(0), func() {})

	}

	log.Println("-------------------------------")
}
