package main

import (
	"database/sql"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type pageProType struct {
	lPro       *tview.List
	descArea   *tview.TextArea
	nameArea   *tview.TextArea
	exportArea *tview.TextArea
	importArea *tview.TextArea
	mPosId     map[int]int
	flexPro    *tview.Flex
	flListPro  *tview.Flex
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

	pagePro.exportArea = tview.NewTextArea()
	pagePro.exportArea.SetBorderColor(tcell.ColorBlue)
	pagePro.exportArea.SetBorderPadding(1, 1, 1, 1)
	pagePro.exportArea.SetDisabled(true)

	pagePro.exportArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("name").
		SetTitleAlign(tview.AlignLeft)

	pagePro.exportArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			hideProExport()
			return nil
		}
		if event.Rune() == 'v' && event.Modifiers() == tcell.ModAlt {
			clipBoardContent, err := clipboard.ReadAll()
			check(err)

			pagePro.exportArea.SetText(pagePro.exportArea.GetText()+clipBoardContent, true)
		}
		return event
	})

	pagePro.importArea = tview.NewTextArea()
	pagePro.importArea.SetBorderColor(tcell.ColorBlue)
	pagePro.importArea.SetBorderPadding(1, 1, 1, 1)
	pagePro.importArea.SetDisabled(true)

	pagePro.importArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("name").
		SetTitleAlign(tview.AlignLeft)

	pagePro.importArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			hideProImport()
			return nil
		}
		if event.Rune() == 'v' && event.Modifiers() == tcell.ModAlt {
			clipBoardContent, err := clipboard.ReadAll()
			check(err)

			pagePro.importArea.SetText(pagePro.importArea.GetText()+clipBoardContent, true)
		}
		return event
	})

	pagePro.flexPro = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pagePro.flListPro, 0, 1, true).
		AddItem(pagePro.Pages, 0, 4, true)

	pagePro.flListPro.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyCtrlQ {
			if pagePro.descArea.GetDisabled() == true {
				if pagePro.nameArea.GetDisabled() == false {
					hideProName()
				}
				if pagePro.exportArea.GetDisabled() == false {
					hideProExport()
				}
				if pagePro.exportArea.GetDisabled() == false {
					hideProImport()
				}
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
				if pagePro.descArea.GetDisabled() == false {
					hideProDesc()
				}
				if pagePro.exportArea.GetDisabled() == false {
					hideProExport()
				}
				if pagePro.exportArea.GetDisabled() == false {
					hideProImport()
				}
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

		if event.Key() == tcell.KeyCtrlS {
			if pagePro.exportArea.GetDisabled() == true {
				if pagePro.descArea.GetDisabled() == false {
					hideProDesc()
				}
				if pagePro.nameArea.GetDisabled() == false {
					hideProName()
				}
				if pagePro.importArea.GetDisabled() == false {
					hideProImport()
				}
				pagePro.exportArea.SetTitle("export")
				pagePro.flListPro.AddItem(pagePro.exportArea, 0, 1, false)
				app.SetFocus(pagePro.exportArea)
				pagePro.exportArea.SetDisabled(false)
			} else {
				hideProExport()
			}

		}

		if event.Key() == tcell.KeyInsert {
			if pagePro.importArea.GetDisabled() == true {
				if pagePro.descArea.GetDisabled() == false {
					hideProDesc()
				}
				if pagePro.nameArea.GetDisabled() == false {
					hideProName()
				}
				if pagePro.importArea.GetDisabled() == false {
					hideProExport()
				}
				pagePro.importArea.SetTitle("import")
				pagePro.flListPro.AddItem(pagePro.importArea, 0, 1, false)
				app.SetFocus(pagePro.importArea)
				pagePro.importArea.SetDisabled(false)
			} else {
				hideProImport()
			}
		}

		return event
	})

	pageMain.pages.AddPage("pro", pagePro.flexPro, true, true)

}

func setListPro() {
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
}

func setProComment() {
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
}

func delPro(idPro int) {

	querySrc := `DELETE FROM src
			    WHERE id_prj = ` + strconv.Itoa(idPro)

	_, err := database.Exec(querySrc)
	check(err)

	queryObj := `DELETE FROM obj
			    WHERE id_prj = ` + strconv.Itoa(idPro)

	_, err = database.Exec(queryObj)

	queryPro := `DELETE FROM prj
			    WHERE id = ` + strconv.Itoa(idPro)

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

	query := "UPDATE prj" + "\n" +
		"SET name = '" + strings.TrimSpace(pagePro.nameArea.GetText()) + "'\n" +
		"WHERE id = " + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])

	_, err := database.Exec(query)
	check(err)
}

func saveProProDesc() {

	query := "UPDATE prj" + "\n" +
		"SET comment = '" + strings.TrimSpace(pagePro.descArea.GetText()) + "'\n" +
		"WHERE id = " + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])

	_, err := database.Exec(query)
	check(err)
}

func getProDesc() string {
	query := `select comment
				from prj` +
		` where id = ` + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])

	pro, err := database.Query(query)
	check(err)

	pro.Next()
	var comment sql.NullString
	err = pro.Scan(&comment)

	pro.Close()

	return comment.String
}

func saveProDescCtrl() {

	var query string
	query = "UPDATE prj" + "\n" +
		"SET comment = '" + pagePro.descArea.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])

	_, err := database.Exec(query)
	check(err)
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

func hideProExport() {
	pagePro.exportArea.SetDisabled(true)
	curPos := pagePro.lPro.GetCurrentItem()

	path := pagePro.exportArea.GetText()
	if len(strings.TrimSpace(path)) > 0 {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.Mkdir(path, 0777)
			check(err)
		}
		downloadPrj(pagePro.lPro.GetCurrentItem(), path)
	}

	pagePro.flListPro.RemoveItem(pagePro.exportArea)
	pagePro.lPro.Clear()
	setListPro()
	pagePro.lPro.SetCurrentItem(curPos)
	app.SetFocus(pagePro.lPro)
}

func hideProImport() {
	log.Println("hideProImport")
	pagePro.importArea.SetDisabled(true)
	curPos := pagePro.lPro.GetCurrentItem()

	path := pagePro.importArea.GetText()

	if len(strings.TrimSpace(path)) > 0 {
		objInfo, err := os.Stat(path)
		check(err)

		mode := objInfo.Mode()
		switch {
		case mode.IsDir():
			importPrj(path)
			//err = importDir(path)
			//check(err)
		default:
		}
	}

	pagePro.flListPro.RemoveItem(pagePro.importArea)
	reloadProTree()
	app.SetFocus(pagePro.lPro)
	if pagePro.lPro.GetItemCount() > curPos {
		pagePro.lPro.SetCurrentItem(curPos)
		setTreePro(curPos)
		setProComment()
	}
}

func downloadPrj(pos int, path string) {

	query := `select name
				from prj
			  where id = ` + strconv.Itoa(pagePro.mPosId[pagePro.lPro.GetCurrentItem()])

	pros, err := database.Query(query)
	check(err)

	pros.Next()
	var prjName sql.NullString
	err = pros.Scan(&prjName)

	objRootPathfile := filepath.Join(path, prjName.String)
	if _, err := os.Stat(objRootPathfile); os.IsNotExist(err) {
		err := os.Mkdir(objRootPathfile, 0777)
		check(err)
	}
	check(err)

	pros.Close()

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
		var objID sql.NullInt64
		var objName sql.NullString
		var objType sql.NullInt16
		err := objects.Scan(&objID, &objName, &objType)
		check(err)

		switch int(objType.Int16) {
		case 0:

			folderPath := filepath.Join(objRootPathfile, objName.String)
			if _, err := os.Stat(folderPath); os.IsNotExist(err) {
				err := os.Mkdir(folderPath, 0777)
				check(err)
			}
			check(err)

			downloadFolderPro(int(objID.Int64), folderPath)
		case 1:

			file, err := os.Create(filepath.Join(objRootPathfile, objName.String))
			if err != nil {
				panic(err)
			}
			downloadFilePro(int(objID.Int64), filepath.Join(objRootPathfile, objName.String))
			defer file.Close()
		}
	}

	objects.Close()
}

func downloadFolderPro(objID int, path string) {

	queryObj := `select id
					   , name
					   , object_type
					from obj
				   where id_parent = ` + strconv.Itoa(objID) +
		` order by object_type asc`

	objects, err := database.Query(queryObj)
	check(err)

	objPath := path

	if _, err := os.Stat(objPath); os.IsNotExist(err) {
		err := os.Mkdir(objPath, 0777)
		check(err)
	}
	check(err)

	for objects.Next() {
		var objID sql.NullInt64
		var objName sql.NullString
		var objType sql.NullInt16
		err := objects.Scan(&objID, &objName, &objType)
		check(err)

		switch int(objType.Int16) {
		case 0:
			if _, err := os.Stat(objPath); os.IsNotExist(err) {
				err := os.Mkdir(objPath, 0777)
				check(err)
			}
			check(err)

			downloadFolderPro(int(objID.Int64), filepath.Join(objPath, objName.String))
		case 1:
			filePath := filepath.Join(objPath, objName.String)

			file, err := os.Create(filePath)
			if err != nil {
				panic(err)
			}
			downloadFilePro(int(objID.Int64), filePath)
			defer file.Close()
		}
	}

	objects.Close()

}

func downloadFilePro(objID int, path string) {

	query := `select line
				from src
			   where id_file = ` + strconv.Itoa(objID) +
		` order by num asc`

	lines, err := database.Query(query)
	check(err)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	check(err)

	for lines.Next() {
		var line sql.NullString
		err := lines.Scan(&line)
		check(err)

		_, err = file.WriteString(line.String + "\n")
	}
	lines.Close()

	defer file.Close()
}
