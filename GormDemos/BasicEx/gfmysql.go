package main

import (
	"log"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID         int `gorm:"primary_key" json:id`
	CreateOn   int `json:created_on`
	ModifiedOn int `json:modified_on`
}

type BlogTag struct {
	Model
	Name       string `json:name`
	CreateBy   string `json:created_by`
	ModifiedBy string `json:modified_by`
	State      int    `json:state`
}

func main() {
	db, err := gorm.Open("mysql", "root:tjazzh203@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	myTag := BlogTag{
		Name:       "test-one",
		CreateBy:   "testadmin",
		ModifiedBy: "testadmin",
		State:      1,
	}

	db.SingularTable(true)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.Exec("INSERT INTO blog_tag(name,created_on,created_by,modified_on,modified_by,state) VALUES ('test-one', 'testadmin','testadmin','0','testadmin','1');", nil)
	db.Create(&myTag)

}
