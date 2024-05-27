package main

import (
	"database/sql"
	"log"
	"sync"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

func initDB() {
	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", "./users.db")
		if err != nil {
			log.Fatal(err)
		}

		// Создание таблицы пользователей, если не существует
		createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL,
		friends TEXT NOT NULL
		);
		`
		_, err = db.Exec(createTableQuery)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func closeDB() {
	if db != nil {
		db.Close()
	}
}