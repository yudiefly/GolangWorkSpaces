package models

type Tag struct {
	Model
	Name       string `json:name`
	CreateBy   string `json:created_by`
	ModifiedBy string `json:modified_by`
	State      int    `json:state`
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
	return
}

func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)
	return
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

//新增标签
func AddTag(name string, state int, createBy string) bool {
	db.Create(&Tag{
		Name:     name,
		State:    state,
		CreateBy: createBy,
	})
	return true
}
