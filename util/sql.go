package util

import (
	"database/sql"
	"log"
)

func connectToDb() *sql.DB {
	db, err := sql.Open("sqlite3", "mydatabase.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func GetQuery(query string) *sql.Rows {
	db := connectToDb()
	defer db.Close()

	results, err := db.Query(query)
	if err != nil {
		panic("Error occurred while executing query")
	}

	return results
}

func GetShmeckles(userId string) {

}
