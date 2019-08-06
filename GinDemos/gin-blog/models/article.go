package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
	//CoverImageUrl string `json:"cover_image_url"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreateOn", time.Now().Unix())
	return nil
}

func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}

func ExistArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id=?", id).First(&article)
	if article.ID > 0 {
		return true
	}
	return false
}

func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)
	return
}

// func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
// 	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)
// 	return
// }

func GetArticles(pageNum int, pageSize int, maps interface{}) ([]*Article, error) {

	// var articles []Article
	// err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	// if err != nil && err != gorm.ErrRecordNotFound {
	// 	return nil, err
	// }

	var articles []*Article
	err := db.Where(maps).Find(&articles).Offset(pageNum).Limit(pageSize).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func GetArticle(id int) (article Article) {
	db.Where("id=?", id).First(&article)
	db.Model(&article).Related(&article.Tag)

	return
}

func EditArticle(id int, data interface{}) bool {
	db.Model(&Article{}).Where("id=?", id).Updates(data)
	return true
}

func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	})
	return true
}

func DeleteArticle(id int) bool {
	db.Where("id=?", id).Delete(Article{})
	return true
}
