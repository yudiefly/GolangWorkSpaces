package models

import (
	"fmt"
	"log"
	"time"

	"gin-blog/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

func init() {

	/*
		marked by zzh203 2019-8-7 14:32
		// var (
		// 	err                                               error
		// 	dbType, dbName, user, password, host, tablePrefix string
		// )

		// sec, err := setting.Cfg.GetSection("database")

		// if err != nil {
		// 	log.Fatal(2, "Fail to get section 'database:'%v", err)
		// }

		// dbType = sec.Key("TYPE").String()

		// dbName = sec.Key("NAME").String()
		// user = sec.Key("USER").String()
		// password = sec.Key("PASSWORD").String()
		// host = sec.Key("HOST").String()

		// tablePrefix = sec.Key("TABLE_PREFIX").String()

		// connStr3 := "root:tjazzh203@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True&loc=Local"
		// connStr2 := user + ":" + password + "@tcp(" + host + ")/" + dbName + "?charset=utf8&parseTime=True&loc=Local"
	*/

	var err error

	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", setting.DatabaseSetting.User, setting.DatabaseSetting.Password, setting.DatabaseSetting.Host, setting.DatabaseSetting.Name)

	// fmt.Println(connStr2)
	// fmt.Println(connStr3)
	// fmt.Println(connStr == connStr3)

	fmt.Println(connStr)

	db, err = gorm.Open(setting.DatabaseSetting.Type, connStr) //gorm.Open(dbType, connStr)

	if err != nil {
		log.Println(err)
	}
	//获取带数据库表名的gorm.DefaultTableNameHandler值
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defualtTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defualtTableName
	}

	db.SingularTable(true)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
}

func CloseDB() {
	defer db.Close()
}

/*
	相关知识点备注：
	检查是否有含有错误（db.Error）
	scope.FieldByName 通过 scope.Fields() 获取所有字段，判断当前是否包含所需字段
*/
// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifyTime` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

/*
	相关知识点备注：
	【scope.Get("gorm:delete_option")】 	检查是否手动指定了delete_option
	【scope.FieldByName("DeletedOn")】 	获取我们约定的删除字段，若存在则 UPDATE 软删除，若不存在则 DELETE 硬删除
	【scope.QuotedTableName() 】			返回引用的表名，这个方法 GORM 会根据自身逻辑对表名进行一些处理
	【scope.CombinedConditionSql() 】	返回组合好的条件SQL，看一下方法原型很明了
	【scope.AddToVars】					该方法可以添加值作为SQL的参数，也可用于防范SQL注入
*/

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
