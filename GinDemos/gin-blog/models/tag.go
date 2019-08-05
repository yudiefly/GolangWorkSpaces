package models

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

//获取Tag列表(带分页功能）
func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var (
		tags []Tag
		err  error
	)

	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Find(&tags).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = db.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

func GetTagTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Tag{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

/*
几点说明：
1、我们创建了一个Tag struct{}，用于Gorm的使用。并给予了附属属性json，这样子在c.JSON的时候就会自动转换格式，非常的便利

2、可能会有的初学者看到return，而后面没有跟着变量，会不理解；其实你可以看到在函数末端，我们已经显示声明了返回值，这个变量在函数体内也可以直接使用，因为他在一开始就被声明了

3、有人会疑惑db是哪里来的；因为在同个models包下，因此db *gorm.DB是可以直接使用的
*/

//判断标签是否存在
func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name=?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

func ExistTagById(id int) bool {
	var tag Tag
	err := db.Select("id").Where("id=? AND deleted_on=?", id, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if tag.ID > 0 {
		return true
	}
	return false
}

//新增标签
func AddTag(name string, state int, createdBy string) bool {
	err := db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}).Error

	if err != nil {
		log.Println(err)
	}

	return true
}

func DeleteTag(id int) bool {
	db.Where("id=?", id).Delete(&Tag{})
	return true
}

func EditTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id=?", id).Updates(data)
	return true
}

/*
  这属于gorm的Callbacks，可以将回调方法定义为模型结构的指针，在创建、更新、查询、删除时将被调用，如果任何回调返回错误，gorm将停止未来操作并回滚所有更改。

  gorm所支持的回调方法：
  创建：BeforeSave、BeforeCreate、AfterCreate、AfterSave
  更新：BeforeSave、BeforeUpdate、AfterUpdate、AfterSave
  删除：BeforeDelete、AfterDelete
  查询：AfterFind
*/

func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}
