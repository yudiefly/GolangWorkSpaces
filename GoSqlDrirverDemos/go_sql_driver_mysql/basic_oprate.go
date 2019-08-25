package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type blogLanguage struct {
	id   int
	name string
	code string
}

var db *sql.DB
var err error

func init() {
	db, err = sql.Open("mysql", "root:tjazzh203@tcp(127.0.0.1)/blog?charset=utf8")
	if err != nil {
		panic(err)
	}
}

func main() {
	// fmt.Println("Insert a row to blog_language table...")
	// insert(blogLanguage{
	// 	name: "zzh203-1",
	// 	code: "yuidefly-1",
	// })
	// fmt.Println("Query datas from blog_languages...")
	// query()
	// update(blogLanguage{
	// 	id:   8,
	// 	name: "朱宗海",
	// 	code: "ZZH",
	// })
	// fmt.Println("Update datas from blog_languages and query...")
	query()
	fmt.Println("Query datas from blog_languages by id...")
	GetLanguage(9)

	//关闭数据库链接
	close()
}

func insert(language blogLanguage) {
	//Insert
	stmt, err := db.Prepare("INSERT INTO blog_language(name, code) VALUES (?, ?)")
	checkErr(err)
	_, err = stmt.Exec(language.name, language.code)
	checkErr(err)

}

func update(language blogLanguage) {
	stmt, err := db.Prepare("UPDATE blog_language set name=?,code=? where id=?")
	checkErr(err)

	res, err := stmt.Exec(language.name, language.code, language.id)
	checkErr(err)

	affect, err := res.RowsAffected()
	fmt.Println(affect)

}

func delete(id int) {
	stmt, err := db.Prepare("delete from blog_language where id=?")
	checkErr(err)

	res, err := stmt.Exec(id)
	checkErr(err)
	fmt.Println(res)
}

func GetLanguage(id int) {
	var uid int
	var name string
	var code string
	err = db.QueryRow("select * from blog_language where id=?", id).Scan(&uid, &name, &code)
	checkErr(err)

	fmt.Printf("id=%d  name=%s code=%s \n", uid, name, code)
}

func query() {
	rows, err := db.Query("select * from blog_language")
	checkErr(err)

	for rows.Next() {
		var uid int
		var name string
		var code string
		each_err := rows.Scan(&uid, &name, &code)
		checkErr(each_err)
		fmt.Printf("id=%d  name=%s code=%s \n", uid, name, code)
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func close() {
	db.Close()
}
