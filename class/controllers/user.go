package controllers

import (
	"class/models"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 注册结构体
type RegController struct {
	beego.Controller
}

// 显示注册页面
func (this *RegController) ShowReg() {
	this.TplName = "user/register.html"
}

// 注册提交处理
func (this *RegController) HandlerReg() {
	// 获取浏览器的数据
	name := this.GetString("userName")
	pwd := this.GetString("password")
	if name == "" || pwd == "" {
		beego.Info("用户名和密码不能为空")
		return
	}
	// 获取orm对象
	o := orm.NewOrm()
	user := models.User{}
	user.Name = name
	user.Pwd = pwd
	// 执行插入
	_, err := o.Insert(&user)
	if err != nil {
		beego.Info("插入失败")
		return
	}
	// 跳转页面
	this.Redirect("/login", 302)
}

// 登陆结构体
type LoginController struct {
	beego.Controller
}

// 显示登陆页面
func (this *LoginController) ShowLogin() {
	userName := this.Ctx.GetCookie("userName")
	if userName != "" {
		this.Data["userName"] = userName
		this.Data["check"] = "checked"
	}
	this.TplName = "user/login.html"
}

// 登陆提交处理
func (this *LoginController) HandlerLogin() {
	// 获取浏览器的数据
	name := this.GetString("userName")
	pwd := this.GetString("password")
	if name == "" || pwd == "" {
		beego.Info("用户名和密码不能为空")
		return
	}
	// 获取orm对象
	o := orm.NewOrm()
	user := models.User{}
	user.Name = name
	// 执行查询
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Info("查询错误")
		return
	}
	// 判断密码是否一致
	if user.Pwd != pwd {
		beego.Info("密码错误")
		return
	}
	// 记住用户名
	check := this.GetString("remember")
	if check == "on" {
		this.Ctx.SetCookie("userName", name, time.Second*3600)
	} else {
		this.Ctx.SetCookie("userName", "sss", -1)
	}

	this.SetSession("userName", name)

	// 跳转页面
	this.Redirect("/Article/article", 302)
}
