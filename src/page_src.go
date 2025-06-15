package main

import (
	"database/sql"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/go-vgo/robotgo"
	"github.com/rivo/tview"
	"log"
	"strconv"
	"strings"
)

type pageSrcType struct {
	lSrc         *tview.List
	descArea     *tview.TextArea
	nameArea     *tview.TextArea
	findArea     *tview.TextArea
	statArea     *tview.TextArea
	mPosId       map[int]int
	curPos       int
	findArr      []int
	importantArr []int
	*tview.Flex
}

var pageSrc pageSrcType

func (pageSrc *pageSrcType) build() {
	pageSrc.mPosId = make(map[int]int)
	pageSrc.importantArr = make([]int, 0)

	pageSrc.lSrc = tview.NewList()

	pageSrc.lSrc.ShowSecondaryText(true).
		SetBorderPadding(1, 1, 1, 1)

	pageSrc.lSrc.SetSelectedBackgroundColor(tcell.ColorOrange)

	pageSrc.lSrc.SetTitle("F4").
		SetTitleAlign(tview.AlignLeft)

	pageSrc.descArea = tview.NewTextArea()
	pageSrc.descArea.SetBorderColor(tcell.ColorBlue)
	pageSrc.descArea.SetBorderPadding(1, 1, 1, 1)
	pageSrc.descArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("comment").
		SetTitleAlign(tview.AlignLeft)

	pageSrc.descArea.SetDisabled(true)

	pageSrc.descArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			hideSrcDesc()
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			descText := strings.Trim(pageSrc.descArea.GetText(), " ")
			arr := strings.Split(descText, "\n")
			pageSrc.descArea.SetText(arr[0], true)
			hideSrcDesc()
		}
		return event
	})

	pageSrc.nameArea = tview.NewTextArea()
	pageSrc.nameArea.SetBorderColor(tcell.ColorBlue)
	pageSrc.nameArea.SetBorderPadding(1, 1, 1, 1)
	pageSrc.nameArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("name").
		SetTitleAlign(tview.AlignLeft)

	pageSrc.nameArea.SetDisabled(true)

	pageSrc.nameArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			hideSrcName()
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			nameText := pageSrc.nameArea.GetText()
			arr := strings.Split(nameText, "\n")
			pageSrc.nameArea.SetText(arr[0], true)
			hideSrcName()
		}
		return event
	})

	pageSrc.findArea = tview.NewTextArea()
	pageSrc.findArea.SetBorderColor(tcell.ColorBlue)
	pageSrc.findArea.SetBorderPadding(1, 1, 1, 1)
	pageSrc.findArea.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(tcell.ColorBlue).
		SetTitle("find").
		SetTitleAlign(tview.AlignLeft)

	pageSrc.findArea.SetDisabled(true)

	pageSrc.findArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyEnter {
			searchText := strings.Trim(pageSrc.findArea.GetText(), " ")
			arr := strings.Split(searchText, "\n")
			pageSrc.findArea.SetText(arr[0], true)
			app.SetFocus(pageSrc.findArea)

			pageSrc.findArr = pageSrc.lSrc.FindItems(searchText, searchText, false, true)
			setFoundLine(true)
			pageSrc.findArea.SetText("", true)
			hideSrcFind()
		}

		if event.Key() == tcell.KeyEsc {
			hideSrcFind()
			return nil
		}

		return event
	})

	pageSrc.statArea = tview.NewTextArea()
	pageSrc.statArea.SetBorderColor(tcell.ColorBlue)
	pageSrc.statArea.SetBorderPadding(1, 1, 1, 1)
	pageSrc.statArea.SetBorder(false)

	pageSrc.lSrc.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		pageSrc.nameArea.SetText(s, true)
		pageSrc.descArea.SetText(s2, true)
		pageSrc.curPos = pageSrc.lSrc.GetCurrentItem()

		pageSrc.statArea.SetText(strconv.Itoa(pageSrc.lSrc.GetItemCount())+": "+strconv.Itoa(i+1), true)
	})

	pageSrc.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pageSrc.lSrc, 0, 17, true).
		AddItem(pageSrc.statArea, 0, 1, false)

	pageSrc.Flex.SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	pageSrc.Flex.SetTitle("F4").
		SetTitleAlign(tview.AlignLeft)

	pageSrc.Flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyCtrlQ {
			if pageSrc.descArea.GetDisabled() == true {
				statHide()
				if pageSrc.nameArea.GetDisabled() == false {
					hideSrcName()
				}
				if pageSrc.findArea.GetDisabled() == false {
					hideSrcFind()
				}
				pageSrc.Flex.AddItem(pageSrc.descArea, 0, 2, false)
				pageSrc.descArea.SetText(getSrcDesc(), true)
				app.SetFocus(pageSrc.descArea)
				pageSrc.descArea.SetDisabled(false)
				saveSrcDesc()
			} else {
				hideSrcDesc()
			}
		}

		if event.Key() == tcell.KeyCtrlF {
			pageSrc.findArea.SetText("", true)
			if pageSrc.findArea.GetDisabled() == true {
				statHide()
				if pageSrc.descArea.GetDisabled() == false {
					hideSrcDesc()
				}
				if pageSrc.nameArea.GetDisabled() == false {
					hideSrcName()
				}
				pageSrc.Flex.AddItem(pageSrc.findArea, 0, 2, false)
				app.SetFocus(pageSrc.findArea)
				pageSrc.findArea.SetDisabled(false)
			} else {
				hideSrcFind()
			}

		}

		if event.Key() == tcell.KeyCtrlO {
			insertEmptyAfter()
			pageSrc.curPos = pageSrc.lSrc.GetCurrentItem() + 1
			pageSrc.lSrc.SetCurrentItem(pageSrc.curPos)
			err := robotgo.KeyTap("Enter")
			check(err)
		}

		if event.Key() == tcell.KeyCtrlW {
			if pageSrc.nameArea.GetDisabled() == true {
				statHide()
				if pageSrc.descArea.GetDisabled() == false {
					hideSrcDesc()
				}
				if pageSrc.findArea.GetDisabled() == false {
					hideSrcFind()
				}
				pageSrc.Flex.AddItem(pageSrc.nameArea, 0, 2, false)
				srcName, _ := pageSrc.lSrc.GetItemText(pageSrc.lSrc.GetCurrentItem())
				pageSrc.nameArea.SetText(srcName+" ", false)
				app.SetFocus(pageSrc.nameArea)
				pageSrc.nameArea.SetDisabled(false)
			} else {
				statHide()
				hideSrcName()
			}
		}

		if event.Key() == tcell.KeyDelete {
			delSrc()
			sanitizeFileSrcLines()
			pageSrc.lSrc.Clear()
			setFileSrc(pageProTree.trPro.GetCurrentNode().GetReference().(int))
			if pageSrc.lSrc.GetItemCount() > pageSrc.curPos {
				pageSrc.lSrc.SetCurrentItem(pageSrc.curPos)
			}
		}

		if event.Key() == tcell.KeyCtrlI {
			setNextImportantLine(true)
		}

		if event.Key() == tcell.KeyCtrlE {
			contentArr := make([]string, 0)
			for ind := range pageSrc.lSrc.GetItemCount() {
				line, _ := pageSrc.lSrc.GetItemText(ind)
				contentArr = append(contentArr, line)
			}

			joined := strings.Join(contentArr, "\n")
			err := clipboard.WriteAll(joined)
			check(err)
		}

		if event.Key() == tcell.KeyCtrlN {
			setFoundLine(true)
		}

		if event.Key() == tcell.KeyCtrlP {
			setFoundLine(false)
		}

		if event.Key() == tcell.KeyDown || event.Key() == tcell.KeyUp || event.Key() == tcell.KeyPgDn || event.Key() == tcell.KeyPgUp {
			err := robotgo.KeyTap("Enter")
			check(err)
		}

		if event.Key() == tcell.KeyTab {
			err := robotgo.KeyTap("Enter")
			check(err)
		}

		if event.Modifiers() == tcell.ModShift && event.Key() == tcell.KeyTab {
			err := robotgo.KeyTap("Enter")
			check(err)
		}

		return event
	})

	pageProTree.Pages.AddPage("src", pageSrc.Flex, true, true)
}

func setFileSrc(idFile int) {
	pageSrc.importantArr = make([]int, 0)

	query := `select id
				   , line
				   , comment
				from src
			   where id_file = ` + strconv.Itoa(idFile) +
		` order by num asc`

	lines, err := database.Query(query)
	check(err)

	posNum := 0
	for lines.Next() {
		posNum++
		var id sql.NullInt64

		var line, comment sql.NullString
		err := lines.Scan(&id, &line, &comment)
		check(err)

		pageSrc.mPosId[posNum-1] = int(id.Int64)

		lineName := strings.ReplaceAll(line.String, "\t", "    ")
		n := len(lineName) - len(strings.TrimSpace(line.String)) - 1
		if n < 0 {
			n = 0
		}
		if len(strings.TrimSpace(comment.String)) > 0 {
			pageSrc.importantArr = append(pageSrc.importantArr, posNum)
		}
		pageSrc.lSrc.AddItem(lineName, strings.Repeat(" ", n)+strings.TrimSpace(comment.String), rune(0), func() {})
	}

	lines.Close()
	statShow()
}

func delSrc() {
	query := `DELETE FROM src
			  WHERE id = ` + strconv.Itoa(pageSrc.mPosId[pageSrc.curPos])

	_, err := database.Exec(query)
	check(err)
}

func updUppNum(nums int) {
	log.Println("updUppNum", nums)
	query := `UPDATE src
    			 SET num = num + ` + strconv.Itoa(nums) +
		` where id_file = ` + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int)) +
		` and num >= ` + strconv.Itoa(pageSrc.curPos+2)

	log.Println(query)
	_, err := database.Exec(query)
	check(err)
}

func (pageSrc *pageSrcType) show() {
	pageProTree.Pages.SwitchToPage("src")
	pageSrc.lSrc.Clear()
	setFileSrc(pageProTree.trPro.GetCurrentNode().GetReference().(int))
	if pageSrc.lSrc.GetItemCount() > 0 {
		pageSrc.curPos = 0
		pageSrc.lSrc.SetCurrentItem(pageSrc.curPos)
	}

	app.SetFocus(pageSrc.Flex)
}

func saveSrcDesc() {
	var query string
	query = "UPDATE src" + "\n" +
		"SET comment = '" + pageSrc.descArea.GetText() + "'\n" +
		"WHERE id = " + strconv.Itoa(pageSrc.mPosId[pageSrc.lSrc.GetCurrentItem()])

	_, err := database.Exec(query)
	check(err)
}

func saveSrcName() {
	var query, val string
	val = strings.TrimRight(pageSrc.nameArea.GetText(), " ")
	if len(val) > 0 {
		val = val + " "
	}
	query = "UPDATE src" + "\n" +
		"SET line = '" + val + "'\n" +
		"WHERE id = " + strconv.Itoa(pageSrc.mPosId[pageSrc.lSrc.GetCurrentItem()])

	_, err := database.Exec(query)
	check(err)
}

func getSrcDesc() string {
	query := `select comment
				from src` +
		` where id = ` + strconv.Itoa(pageSrc.mPosId[pageSrc.lSrc.GetCurrentItem()])

	pro, err := database.Query(query)
	check(err)

	pro.Next()
	var comment sql.NullString
	err = pro.Scan(&comment)

	pro.Close()

	return comment.String
}

func hideSrcDesc() {
	pageSrc.descArea.SetDisabled(true)
	curPos := pageSrc.lSrc.GetCurrentItem()
	saveSrcDesc()
	pageSrc.Flex.RemoveItem(pageSrc.descArea)
	statShow()
	pageSrc.show()
	pageSrc.lSrc.SetCurrentItem(curPos)
	app.SetFocus(pageSrc.lSrc)
}

func hideSrcFind() {
	pageSrc.curPos = pageSrc.lSrc.GetCurrentItem()
	pageSrc.findArea.SetDisabled(true)
	pageSrc.lSrc.SetCurrentItem(pageSrc.curPos)
	pageSrc.Flex.RemoveItem(pageSrc.findArea)
	statShow()

	app.SetFocus(pageSrc.lSrc)
}

func hideSrcName() {
	pageSrc.nameArea.SetDisabled(true)

	saveSrcName()
	pageSrc.Flex.RemoveItem(pageSrc.nameArea)
	statShow()
	curPos := pageSrc.lSrc.GetCurrentItem()
	pageSrc.show()
	pageSrc.lSrc.SetCurrentItem(curPos)
	app.SetFocus(pageSrc.lSrc)
}

func importSrc() {
	content, err := clipboard.ReadAll()
	check(err)

	prjID = pagePro.mPosId[pagePro.lPro.GetCurrentItem()]
	curFileID = pageProTree.trPro.GetCurrentNode().GetReference().(int)

	lines := strings.Split(content, "\n")
	updUppNum(len(lines))
	for i, line := range lines {
		escapedLine := sanitizeString(line)
		saveSrc(curFileID, pageSrc.curPos+2+i, escapedLine)
	}
	sanitizeFileSrcLines()

	curPos := pageSrc.lSrc.GetCurrentItem()
	pageSrc.show()
	pageSrc.lSrc.SetCurrentItem(curPos)
}

func sanitizeFileSrcLines() {
	query := `WITH numbered_rows AS (
				  SELECT 
					id,
					ROW_NUMBER() OVER (ORDER BY num ASC) AS new_num
				  FROM src
				  WHERE id_file = ` + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int)) + `
				)
			UPDATE src
			SET num = nr.new_num
			FROM numbered_rows nr
			WHERE src.id = nr.id
			AND src.id_file = ` + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))
	_, err := database.Exec(query)
	check(err)
}

func insertEmptyAfter() {
	prjID = pagePro.mPosId[pagePro.lPro.GetCurrentItem()]
	curFileID = pageProTree.trPro.GetCurrentNode().GetReference().(int)

	updUppNum(1)
	saveSrc(curFileID, pageSrc.curPos+2, "")
	sanitizeFileSrcLines()

	curPos := pageSrc.lSrc.GetCurrentItem()
	pageSrc.show()
	pageSrc.lSrc.SetCurrentItem(curPos)

}

func statHide() {
	pageSrc.Flex.RemoveItem(pageSrc.statArea)
	pageSrc.Flex.SetBorder(false)
	pageSrc.lSrc.SetBorder(true)
	pageSrc.lSrc.SetBorderColor(tcell.ColorBlue)
}

func statShow() {
	if pageSrc.statArea.GetDisabled() == true {
		pageSrc.Flex.AddItem(pageSrc.statArea, 0, 1, false)
		pageSrc.Flex.SetBorderColor(tcell.ColorBlue)
		pageSrc.Flex.SetBorder(true)
		pageSrc.lSrc.SetBorder(false)

		pageSrc.statArea.SetText(strconv.Itoa(pageSrc.lSrc.GetItemCount())+": "+strconv.Itoa(1), true)
		pageSrc.statArea.SetDisabled(false)
	}
}

func setFoundLine(bStraightWay bool) {
	bFound := false
	if bStraightWay {
		for ind := range pageSrc.findArr {
			if pageSrc.findArr[ind] > pageSrc.lSrc.GetCurrentItem() {
				pageSrc.lSrc.SetCurrentItem(pageSrc.findArr[ind])
				err := robotgo.KeyTap("Enter")
				check(err)
				bFound = true
				break
			}
		}
		if !bFound && len(pageSrc.findArr) > 0 {
			pageSrc.lSrc.SetCurrentItem(pageSrc.findArr[0])
		}
	} else {

		for ind := len(pageSrc.findArr) - 1; ind >= 0; ind-- {
			if pageSrc.findArr[ind] < pageSrc.lSrc.GetCurrentItem() {
				pageSrc.lSrc.SetCurrentItem(pageSrc.findArr[ind])
				err := robotgo.KeyTap("Enter")
				check(err)
				bFound = true
				break
			}
		}
		if !bFound && len(pageSrc.findArr) > 0 {
			pageSrc.lSrc.SetCurrentItem(pageSrc.findArr[len(pageSrc.findArr)-1])
		}
	}
}

func setNextImportantLine(bStraightWay bool) {
	bFound := false
	if bStraightWay {
		for ind := range pageSrc.importantArr {
			if pageSrc.importantArr[ind] > pageSrc.lSrc.GetCurrentItem() {
				pageSrc.lSrc.SetCurrentItem(pageSrc.importantArr[ind])
				err := robotgo.KeyTap("Enter")
				check(err)
				bFound = true
				break
			}
		}
		if !bFound && len(pageSrc.importantArr) > 0 {
			pageSrc.lSrc.SetCurrentItem(pageSrc.importantArr[0])
		}
	} else {
		for ind := len(pageSrc.importantArr) - 1; ind >= 0; ind-- {
			if pageSrc.importantArr[ind] < pageSrc.lSrc.GetCurrentItem() {
				pageSrc.lSrc.SetCurrentItem(pageSrc.importantArr[ind])
				err := robotgo.KeyTap("Enter")
				check(err)
				bFound = true
				break
			}
		}
		if !bFound && len(pageSrc.importantArr) > 0 {
			pageSrc.lSrc.SetCurrentItem(pageSrc.importantArr[len(pageSrc.importantArr)-1])
		}
	}
}
