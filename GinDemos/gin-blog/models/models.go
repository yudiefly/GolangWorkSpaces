package models

import (
	"fmt"
	"log"

	"gin-blog/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:id`
	CreateOn   int `json:created_on`
	ModifiedOn int `json:modified_on`
}

func init() {
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database:'%v", err)
	}
	dbType = sec.GetKey("TYPE").String()
	dbName = sec.GetKey("NAME").String()
	user = sec.GetKey("USER")
	password = sec.Key("PASSWORD").String()
	host = sec.GetKey("HOST").String()
	tablePrefix = sec.GetKey("TABLE_PREFIX")

	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=ture&loc=Local", user, password, host, dbName))

	if err != nil {
		log.Println(err)
	}
	//获取带数据库表名的gorm.DefaultTableNameHandler值
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defualtTableName string) string {
		return tablePrefix + defualtTableName
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}
func CloseDB() {
	defer db.Close()
}
