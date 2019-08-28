package main

import (
	// "database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

type BlogLanguage struct {
	Id   int
	Name string
	Code string
}

var db *sqlx.DB
var err error

func init() {
	db, err = sqlx.Open("mysql", "root:tjazzh203@tcp(127.0.0.1)/blog?charset=utf8")
	checkErr(err)
}

func main() {
	getLanguage(9)
}

func getLanguage(id int) {
	language := BlogLanguage{}
	err = db.QueryRowx("select id,name,code from blog_language where id=?", id).StructScan(&language)
	checkErr(err)
	fmt.Println(language)
	fmt.Printf("id=%d  name=%s code=%s \n", language.Id, language.Name, language.Code)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func close() {
	db.Close()
}
