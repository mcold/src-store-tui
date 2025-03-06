package main

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
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
	db, err := sql.Open("sqlite", "DB.db")
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
											   from obj)`

	log.Println(queryParent)
	_, err := database.Exec(queryParent)
	check(err)

	queryObj := `DELETE FROM obj
			    WHERE id_prj is null
                   or id_prj not in (select id
                        				from prj)`

	log.Println(queryObj)
	_, err = database.Exec(queryObj)
	check(err)

	querySrc := `DELETE FROM src
			    WHERE id_file is null
                   or id_file not in (select id
                        				from obj)`

	log.Println(querySrc)
	_, err = database.Exec(querySrc)
	check(err)
}
