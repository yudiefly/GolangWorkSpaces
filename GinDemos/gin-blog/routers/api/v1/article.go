package v1

import (
	//"log"
	"net/http"

	//"gin-blog/models"
	"gin-blog/pkg/e"
	//"gin-blog/pkg/logging"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/boombuler/barcode/qr"

	"gin-blog/pkg/app"
	"gin-blog/pkg/qrcode"
	"gin-blog/service/article_service"
	"gin-blog/service/tag_service"
)

// @Summary 获取单个文章
// @Produce  json
// @Parameter  id query int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"系统错误"}"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	/*
		id, _ := com.StrTo(c.Param("id")).Int()

		valid := validation.Validation{}
		valid.Min(id, 1, "id").Message("ID必须大于0")

		code := e.INVALID_PARAMS
		var data interface{}
		if !valid.HasErrors() {
			if models.ExistArticleByID(id) {
				data = models.GetArticle(id)
				code = e.SUCCESS
			} else {
				code = e.ERROR_NOT_EXIST_ARTICLE
			}
		} else {
			for _, err := range valid.Errors {
				log.Println(err.Key, err.Message)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
	*/

	//抽象出article_service层后，代码更为清晰了
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, article)
}

// @Summary 获取多个文章
// @Produce  json
// @Parameter  tag_id query int false "TagID"
// @Parameter  state query int false "State"
// @Parameter  created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{"list":[],"TotalCount":8},"msg":"ok"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"系统错误"}"
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	/*
		data := make(map[string]interface{})
		maps := make(map[string]interface{})
		valid := validation.Validation{}

		var state int = -1
		if arg := c.Query("state"); arg != "" {
			state, _ = com.StrTo(arg).Int()
			maps["state"] = state

			valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
		}

		var tagId int = -1
		if arg := c.Query("tag_id"); arg != "" {
			tagId, _ = com.StrTo(arg).Int()
			maps["tag_id"] = tagId

			valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
		}
		code := e.INVALID_PARAMS
		if !valid.HasErrors() {
			code = e.SUCCESS

			//list, err := models.GetArticles(util.GetPage(c), setting.PageSize, maps)
			list, err := models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
			if err != nil {
				//log.Fatal(err)
				logging.Fatal(err)
			} else {
				data["list"] = list
			}

			data["total"] = models.GetArticleTotal(maps)

		} else {
			for _, err := range valid.Errors {
				//log.Println(err.Key, err.Message)
				logging.Info(err.Key, err.Message)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
	*/

	//抽象出article_service层，让代码更为清晰
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state")
	}

	tagId := -1
	if arg := c.PostForm("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	total, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["list"] = articles
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

////新增文章的对象
// type AddArticleForm struct {
// 	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
// 	Title         string `form:"title" valid:"Required;MaxSize(100)"`
// 	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
// 	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
// 	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
// 	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
// 	State         int    `form:"state" valid:"Range(0,1)"`
// }

type AddArticleForm struct {
	TagID         int    `form:"tag_id" json:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" json:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" json:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" json:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by" json:"created_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" json:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" json:"state" valid:"Range(0,1)"`
}

// @Summary 新增文章
// @Produce  json
// @Parameter  tag_id query int true "TagID"
// @Parameter  title query string true "Title"
// @Parameter  desc query string true "Desc"
// @Parameter  content query string true "Content"
// @Parameter  created_by query string true "CreatedBy"
// @Parameter  state query int true "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"系统错误"}"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	/*
		tagId, _ := com.StrTo(c.Query("tag_id")).Int()
		title := c.Query("title")
		desc := c.Query("content")
		content := c.Query("content")
		createdBy := c.Query("created_by")
		state, _ := com.StrTo(c.DefaultQuery("state", "0")).Int()
		coverImageUrl := c.Query("cover_image_url")

		valid := validation.Validation{}
		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
		valid.Required(title, "title").Message("标题不能为空")
		valid.Required(desc, "desc").Message("简述不能为空")
		valid.Required(content, "content").Message("内容不能为空")
		valid.Required(createdBy, "created_by").Message("创建人不能为空")
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

		code := e.INVALID_PARAMS
		if !valid.HasErrors() {
			if models.ExistTagById(tagId) {
				data := make(map[string]interface{})
				data["tag_id"] = tagId
				data["title"] = title
				data["desc"] = desc
				data["content"] = content
				data["created_by"] = createdBy
				data["state"] = state
				data["cover_image_url"] = coverImageUrl
				models.AddArticle(data)
				code = e.SUCCESS
			} else {
				code = e.ERROR_NOT_EXIST_TAG
			}
		} else {
			for _, err := range valid.Errors {
				log.Println(err.Key, err.Message)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})

	*/

	//抽象出article_service层，让代码更为清晰
	var (
		appG = app.Gin{C: c}
		form AddArticleForm
	)

	//httpCode, errCode := app.BindAndValid(c, &form)
	httpCode, errCode := app.BindPostJsonAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
	}
	if err := articleService.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//修改文章的对象
type EditArticleForm struct {
	ID            int    `form:"id" valid:"Required;Min(1)"`
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// @Summary 修改文章
// @Produce  json
// @Parameter  id path int true "ID"
// @Parameter  tag_id query string false "TagID"
// @Parameter  title query string false "Title"
// @Parameter  desc query string false "Desc"
// @Parameter  content query string false "Content"
// @Parameter  modified_by query string true "ModifiedBy"
// @Parameter  state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"系统错误"}"
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	/*
		valid := validation.Validation{}

		id, _ := com.StrTo(c.Param("id")).Int()
		tagId, _ := com.StrTo(c.Query("tag_id")).Int()
		title := c.Query("title")
		desc := c.Query("desc")
		content := c.Query("content")
		modifiedBy := c.Query("modified_by")
		coverImageUrl := c.Query("cover_image_url")

		var state int = -1
		if arg := c.Query("state"); arg != "" {
			state, _ = com.StrTo(arg).Int()
			valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
		}

		valid.Min(id, 1, "id").Message("ID必须大于0")
		valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
		valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
		valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
		valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
		valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

		code := e.INVALID_PARAMS
		if !valid.HasErrors() {
			if models.ExistArticleByID(id) {
				if models.ExistTagById(tagId) {
					data := make(map[string]interface{})
					if tagId > 0 {
						data["tag_id"] = tagId
					}
					if title != "" {
						data["title"] = title
					}
					if desc != "" {
						data["desc"] = desc
					}
					if content != "" {
						data["content"] = content
					}

					data["modified_by"] = modifiedBy

					data["cover_image_url"] = coverImageUrl

					models.EditArticle(id, data)
					code = e.SUCCESS
				} else {
					code = e.ERROR_NOT_EXIST_TAG
				}
			} else {
				code = e.ERROR_NOT_EXIST_ARTICLE
			}
		} else {
			for _, err := range valid.Errors {
				//log.Println(err.Key, err.Message)
				logging.Info(err.Key, err.Message)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]string),
		})
	*/

	//抽象出article_service层，让代码更为清晰
	var (
		appG = app.Gin{C: c}
		form = EditArticleForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	articleService := article_service.Article{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		ModifiedBy:    form.ModifiedBy,
		State:         form.State,
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = articleService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 删除文章
// @Produce  json
// @Parameter  id query int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"系统错误"}"
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	/*
		id, _ := com.StrTo(c.Param("id")).Int()

		valid := validation.Validation{}
		valid.Min(id, 1, "id").Message("ID必须大于0")

		code := e.INVALID_PARAMS
		if !valid.HasErrors() {
			if models.ExistArticleByID(id) {
				models.DeleteArticle(id)
				code = e.SUCCESS
			} else {
				code = e.ERROR_NOT_EXIST_ARTICLE
			}
		} else {
			for _, err := range valid.Errors {
				log.Println(err.Key, err.Message)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]string),
		})
	*/

	//抽象出article_service层，让代码更清晰
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

const (
	QRCODE_URL = "https://github.com/yudiefly/GolangWorkSpaces/"
)

func GetnerateArticlePoster(c *gin.Context) {
	/*
		appG := app.Gin{C: c}
		qrc := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
		path := qrcode.GetQrCodeFullPath()
		_, _, err := qrc.Encode(path)
		if err != nil {
			appG.Response(http.StatusOK, e.ERROR, nil)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	*/
	appG := app.Gin{C: c}
	article := &article_service.Article{}
	qr := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
	posterName := article_service.GetPosterFlag() + "_" + qrcode.GetQrCodeFileName(qr.URL) + qr.GetQrCodeExt()
	articlePoster := article_service.NewArticlePoster(posterName, article, qr)
	articlePosterBgService := article_service.NewArticlePosterBg(
		"bg.jpg",
		articlePoster,
		&article_service.Rect{
			X0: 0,
			Y0: 0,
			X1: 550,
			Y1: 700,
		},
		&article_service.Pt{
			X: 125,
			Y: 298,
		},
	)

	_, filePath, err := articlePosterBgService.Generate()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GEN_ARTICLE_POSTER_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"poster_url":      qrcode.GetQrCodeFullUrl(posterName),
		"poster_save_url": filePath + posterName,
	})
}
