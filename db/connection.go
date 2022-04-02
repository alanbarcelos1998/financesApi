package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func Connection() {
	var erroAbertura error

	var connectDb string = fmt.Sprintf("%s:%s@tcp(%s:3307)/%s", "root", "root", "127.0.0.1", "finances")

	Db, erroAbertura = sql.Open("mysql", connectDb)

	if erroAbertura != nil {
		log.Fatal(erroAbertura.Error())
	}

	erroPing := Db.Ping()

	if erroPing != nil {
		log.Fatal(erroPing.Error())
	}

}
