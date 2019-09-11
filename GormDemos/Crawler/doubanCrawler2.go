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
