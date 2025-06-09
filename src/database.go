package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"runtime"
)

type databaseType struct {
	*sql.DB
	Path             string
	ConnectionString string
}

var database databaseType

func (database *databaseType) buildConnectionString() {
	database.ConnectionString = database.Path
}

func (database *databaseType) Connect() error {
	db, err := sql.Open("sqlite", "DBs"+string(os.PathSeparator)+os.Args[1]+".db")
	check(err)

	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	database.DB = db
	return nil
}

func check(err interface{}) {
	if err != nil {
		_, fileName, lineNo, _ := runtime.Caller(1) // Получаем информацию о вызывающем файле
		log.Printf("%s: %d\n", filepath.Base(fileName), lineNo)
		log.Println(err)

		panic(err)
	}
}

func delAbsParent() {
	queryParent := `DELETE FROM obj
					 WHERE id_parent not in (select id
											   from obj) and id_parent != 0`
	_, err := database.Exec(queryParent)
	check(err)

	queryObj := `DELETE FROM obj
			    WHERE id_prj is null
                   or id_prj not in (select id
                        				from prj)`
	_, err = database.Exec(queryObj)
	check(err)

	querySrc := `DELETE FROM src
			    WHERE id_file is null
                   or id_file not in (select id
                        				from obj)`
	_, err = database.Exec(querySrc)
	check(err)
}

func savePrj(prjName string) {
	// TODO: make function return ID: clear call getLastProjectID
	query := fmt.Sprintf("INSERT INTO PRJ (id_item, name) VALUES (1, '%s');\n", prjName)

	_, err := database.Exec(query)
	check(err)
}

func saveObj(objName string, prjID int, parentID int, objType int) {

	query := fmt.Sprintf("INSERT INTO OBJ (id_prj, id_parent, name, object_type)"+
		"VALUES (%d, %d, '%s', %d);\n", prjID, parentID, objName, objType)

	_, err := database.Exec(query)
	check(err)
}

func saveSrc(idFile int, num int, line string) {
	query := fmt.Sprintf("INSERT INTO SRC (id_prj, id_file, num, line) VALUES ( %d, %d, %d, '%s');\n",
		prjID, idFile, num, line)

	_, err := database.Exec(query)
	check(err)
}

func getLastProjectID() (int, error) {
	query := `select max(id)
				from prj`

	var maxID int
	err := database.QueryRow(query).Scan(&maxID)
	check(err)

	return maxID, nil
}

func getLastObjectID() (int, error) {
	query := `select max(id)
				from obj`

	var maxID int
	err := database.QueryRow(query).Scan(&maxID)
	check(err)

	return maxID, nil
}
