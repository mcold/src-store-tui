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
