package main

import (
	"log"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID         int `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

type Language struct {
	ID   int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name string `gorm:"name"`
	Code string `gorm:"code"`
}

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func main() {
	db, err := gorm.Open("mysql", "root:tjazzh203@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	//本来是按照模型的名字来定位表名，该行则指定每个模型的名字前面加上”blog_”前缀后才是表名
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defualtTableName string) string {
		return "blog_" + defualtTableName
	}

	defer db.Close()

	db.SingularTable(true)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	//db.Exec()方法执行，下面的代码是错误的
	//db.Exec("INSERT INTO blog_tag(name,created_on,created_by,modified_on,modified_by,state) VALUES ('test-one-in', 'testadmin','testadmin','0','testadmin','1')", nil)

	tt := db.Create(&Tag{
		Name:       "test-one-2",
		CreatedBy:  "testadmin-2",
		ModifiedBy: "testadmin-2",
		State:      0,
	}).Error

	log.Println(tt)

	db.Create(&Language{
		Code: "c005",
		Name: "NM-0005",
	})

}
