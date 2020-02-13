package controllers

import (
	"class/models"
	"math"
	"path"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
)

// 文章结构体
type ArticleController struct {
	beego.Controller
}

// 处理下拉框改变法的请求
func (this *ArticleController) HandlerSelect() {
	// 1.接受数据
	typeName := this.GetString("select")
	// 2.处理数据
	if typeName == "" {
		beego.Info("下拉框传递数据失败")
		return
	}
	// 3.查询数据
	o := orm.NewOrm()
	var article []models.Article
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__Typename", typeName).All(&article)
	// beego.Info(article)
}

// 显示文章列表
func (this *ArticleController) ShowArticle() {

	// 获取orm对象
	o := orm.NewOrm()
	var article []models.Article
	ps := o.QueryTable("Article")

	// 获取页码
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	count, _ := ps.RelatedSel("ArticleType").Count()           // 总几条
	pageSize := 3                                              // 显示6条记录
	pageCount := math.Ceil(float64(count) / float64(pageSize)) // 总几页
	start := pageSize * (pageIndex - 1)
	ps.Limit(pageSize, start).All(&article)

	// 解决上一页和下一页
	FirstPage, EndPage := false, false
	if pageIndex == 1 {
		FirstPage = true
	}
	if pageIndex == int(pageCount) {
		EndPage = true
	}

	// 显示文章类型下拉框
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types

	// 1.接受数据
	typeName := this.GetString("select")
	var articlewithType []models.Article
	// 2.处理数据
	if typeName == "" {
		beego.Info("下拉框传递数据失败")
		ps.Limit(pageSize, start).RelatedSel("ArticleType").All(&articlewithType)
	} else {
		ps.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__Typename", typeName).All(&articlewithType)
	}

	// 跳转页面
	this.Data["count"] = count         // 总几条
	this.Data["pageCount"] = pageCount // 总几页
	this.Data["pageIndex"] = pageIndex // 当前几页
	this.Data["FirstPage"] = FirstPage
	this.Data["typeName"] = typeName
	this.Data["EndPage"] = EndPage
	this.Data["article"] = articlewithType

	this.Layout = "layout.html"
	this.TplName = "index.html"
}

// 显示文章添加页面
func (this *ArticleController) ShowAddArticle() {
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types
	this.TplName = "articles/add.html"
}

// 文章添加提交处理
func (this *ArticleController) HandlerAddArticle() {
	// 获取浏览器的数据
	artiName := this.GetString("articleName")
	artiContent := this.GetString("content")
	f, h, err := this.GetFile("uploadname")
	defer f.Close()
	// 判断文件格式
	Ext := path.Ext(h.Filename)
	if Ext != ".png" && Ext != ".jpg" && Ext != ".jpeg" {
		beego.Info("文件格式错误")
		return
	}
	// 文件大小
	if h.Size > 5000000 {
		beego.Info("图片太大，不允许上传")
		return
	}
	// 保存路径
	fileName := time.Now().Format("2006-01-02 15:04:05")
	this.SaveToFile("uploadname", "./static/img/"+fileName+Ext)
	if err != nil {
		beego.Info("上传失败")
		return
	}
	// 获取orm对象
	o := orm.NewOrm()
	article := models.Article{}
	article.Title = artiName
	article.Content = artiContent
	article.Img = "/static/img/" + fileName + Ext
	// 获取到
	typeName := this.GetString("select")
	if typeName == "" {
		beego.Info("下拉框数据错误")
		return
	}
	var artiType models.ArticleType
	artiType.Typename = typeName
	err = o.Read(&artiType, "Typename")
	if err != nil {
		beego.Info("获取数据类型")
		return
	}
	article.ArticleType = &artiType
	// 执行插入
	_, err = o.Insert(&article)
	if err != nil {
		beego.Info("插入失败")
		return
	}
	// 跳转页面
	this.Redirect("/Article/article", 302)
}

// 显示文章类型页面
func (this *ArticleController) ShowAddType() {
	// 获取orm对象
	o := orm.NewOrm()
	var article []models.ArticleType
	o.QueryTable("ArticleType").All(&article)
	this.Data["article"] = article
	this.TplName = "articles/addType.html"
}

// 文章类型添加提交处理
func (this *ArticleController) HandlerAddType() {
	// 获取浏览器的数据
	typename := this.GetString("typeName")
	if typename == "" {
		beego.Info("添加失败")
		return
	}
	// 获取orm对象
	o := orm.NewOrm()
	article := models.ArticleType{}
	article.Typename = typename
	// 执行插入
	_, err := o.Insert(&article)
	if err != nil {
		beego.Info("插入失败")
		return
	}
	// 跳转页面
	this.Redirect("/Article/addType", 302)
}

// 文章类型删除提交出处理
func (this *ArticleController) HandlerDelType() {
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	article := models.ArticleType{Id: id}
	_, err := o.Delete(&article)
	if err != nil {
		beego.Info("删除失败")
		return
	}
	// 跳转页面
	this.Redirect("/Article/addType", 302)
}

// 显示文章内容页面
func (this *ArticleController) ShowArtiContent() {
	// 获取浏览器的Id
	id, _ := this.GetInt("id")
	// 获取orm对象
	o := orm.NewOrm()
	article := models.Article{Id: id}
	// 执行查询
	err := o.Read(&article)
	if err != nil {
		beego.Info("查询失败")
		return
	}
	// 阅读量
	article.Count += 1
	o.Update(&article)
	// 显示内容
	this.Data["article"] = article
	this.Layout = "layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["contentHead"] = "articles/head.html"
	this.TplName = "articles/content.html"
}

// 文章删除处理
func (this *ArticleController) HandlerDelete() {
	// 获取浏览器的id
	id, _ := this.GetInt("id")
	// 获取orm对象
	o := orm.NewOrm()
	article := models.Article{Id: id}
	// 执行删除
	_, err := o.Delete(&article)
	if err != nil {
		beego.Info("删除失败")
		return
	}
	// 跳转页面
	this.Redirect("/Article/article", 302)
}

// 显示文章编辑页面
func (this *ArticleController) ShowArtiUpdate() {
	// 获取浏览器的id
	id, _ := this.GetInt("id")
	// 获取orm对象
	o := orm.NewOrm()
	article := models.Article{Id: id}
	// 执行查询
	err := o.Read(&article)
	if err != nil {
		beego.Info("查询失败")
		return
	}
	// 显示内容
	this.Data["article"] = article
	this.TplName = "articles/update.html"

}

// 文章编辑提交处理
func (this *ArticleController) HandlerArtiUpdate() {
	// 获取浏览器的id
	id, _ := this.GetInt("id")
	artiName := this.GetString("articleName")
	artiContent := this.GetString("content")
	f, h, err := this.GetFile("uploadname")
	defer f.Close()
	// 判断文件格式
	Ext := path.Ext(h.Filename)
	if Ext != ".png" && Ext != ".jpg" && Ext != ".jpeg" {
		beego.Info("文件格式错误")
		return
	}
	// 文件大小
	if h.Size > 5000000 {
		beego.Info("图片太大，不允许上传")
		return
	}
	// 保存路径
	fileName := time.Now().Format("2006-01-02 15:04:05")
	this.SaveToFile("uploadname", "./static/img/"+fileName+Ext)
	if err != nil {
		beego.Info("不允许上传文件")
		return
	}
	// 获取orm对象
	o := orm.NewOrm()
	article := models.Article{}
	article.Id = id
	article.Title = artiName
	article.Content = artiContent
	article.Img = "/static/img/" + fileName + Ext
	// 执行更新
	_, err = o.Update(&article)
	if err != nil {
		beego.Info("更新失败")
		return
	}
	// 跳转页面
	this.Redirect("/Article/article", 302)
}

// 退出登录
func (this *ArticleController) Logout() {
	this.DelSession("userName")
	this.Redirect("/login", 302)
}
