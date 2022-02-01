package dbclient

import "github.com/jmoiron/sqlx"

var DBClient *sqlx.DB

func InitialiseDBConnection() {
	db, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/core")
	if err != nil {
		panic(err.Error())
	}

	DBClient = db
}