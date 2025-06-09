package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	prjID     int
	curDirID  int
	curFileID int
	stack     Stack[int]
)

type Stack[T any] struct {
	elements []T
}

func (s *Stack[T]) Push(item T) {
	s.elements = append(s.elements, item)
}

func (s *Stack[T]) Pop() (T, error) {
	var zero T
	if len(s.elements) == 0 {
		return zero, nil
	}

	item := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return item, nil
}

func importPrj(path string) {

	prjName := filepath.Base(path)
	savePrj(prjName)
	var err error
	prjID, err = getLastProjectID()
	check(err)

	err = importDir(path)
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}
}

func importDir(path string) error {

	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			if curDirID != 0 {
				stack.Push(curDirID)
			}

			saveObj(entry.Name(), prjID, curDirID, 0)
			curDirID, err = getLastObjectID()
			check(err)

			err := importDir(fullPath)
			if err != nil {
				return err
			}

			curDirID, err = stack.Pop()
			check(err)
		} else {
			err = importFile(entry.Name(), fullPath, 1)
			check(err)
		}
	}
	return nil
}

func importObj(path string) error {
	curPrjPos := pagePro.lPro.GetCurrentItem()
	curObjNode := pageProTree.trPro.GetCurrentNode()
	prjID = pagePro.mPosId[pagePro.lPro.GetCurrentItem()]
	idDir := getCurDirID()

	if idDir != 0 {
		curDirID = idDir
	}

	objInfo, err := os.Stat(path)
	check(err)

	mode := objInfo.Mode()
	switch {
	case mode.IsDir():
		saveObj(filepath.Base(path), prjID, curDirID, 0)
		curDirID, err = getLastObjectID()
		check(err)

		err = importDir(path)
		check(err)
	case mode.IsRegular():
		err = importFile(filepath.Base(path), path, 1)
		check(err)
	}

	reloadProTree()
	pagePro.lPro.SetCurrentItem(curPrjPos)
	pageProTree.trPro.SetCurrentNode(curObjNode)

	return nil
}

func getCurDirID() int {
	curNode := pageProTree.trPro.GetCurrentNode()
	if curNode == nil {
		return 0
	}

	query := `select o.id
				   , o.object_type
				   , o.id_parent
				from obj o` +
		` where id = ` + strconv.Itoa(pageProTree.trPro.GetCurrentNode().GetReference().(int))

	objs, err := database.Query(query)
	check(err)

	objs.Next()
	var id, idParent sql.NullInt64
	var objType sql.NullInt16
	err = objs.Scan(&id, &objType, &idParent)
	check(err)

	objs.Close()

	switch int(objType.Int16) {
	case 0:
		return int(id.Int64)
	case 1:
		if int(idParent.Int64) > 0 {
			return int(idParent.Int64)
		} else {
			return 0
		}
	}
	return 0
}

func processFile(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		escapedLine := sanitizeString(line)
		saveSrc(curFileID, i+1, escapedLine)
	}
	return nil
}

func sanitizeString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func importFile(name string, path string, objType int) error {

	saveObj(name, prjID, curDirID, objType)
	var err error
	curFileID, err = getLastObjectID()
	if err != nil {
		return err
	}

	err = processFile(path)
	if err != nil {
		return err
	}

	return err
}
