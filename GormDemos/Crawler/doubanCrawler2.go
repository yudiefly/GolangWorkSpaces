package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Movie struct {
	ID    int    `gorm:"AUTO_INCREMENT"`
	Title string `gorm:"type:varchar(100);unique_index"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := gorm.Open("mysql", "root:tjazzh203@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True&loc=Local")

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defualtTableName string) string {
		return "blog_" + defualtTableName
	}

	defer db.Close()

	db.SingularTable(true)

	checkError(err)

	var movie Movie
	db.First(&movie, 1291552) // 查询ID为1291552的记录，只支持主键
	log.Println(movie)
	log.Println(movie.ID, movie.Title) // 可以获得对应属性

	var movies []Movie
	db.Order("id").Limit(13).Find(&movies) // 按ID升序排，取前3个记录赋值给movies
	log.Println(movies)
	log.Println("ID      Title")
	for _, emv := range movies {
		log.Println(emv.ID, emv.Title)
	}

	db.Order("id desc").Limit(10).Offset(1).Find(&movies) // Order也支持desc选择降序, offset表示对结果集从第2个记录开始
	log.Println(movies)

	db.Select("title").First(&movies, "title=?", "风之谷") // 用Select可以限定只返回那些字段，First也支持条件
	log.Println(movie)

	var counts int64
	db.Where("id=?", 1291552).Or("title=?", "风之谷").Find(&movies).Count(&counts)
	log.Println(counts)

}
