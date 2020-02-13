package routers

import (
	"class/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/Article/*", beego.BeforeRouter, FilterFunc)
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.RegController{}, "get:ShowReg;post:HandlerReg")                             // 注册
	beego.Router("/login", &controllers.LoginController{}, "get:ShowLogin;post:HandlerLogin")                          // 登陆
	beego.Router("/Article/article", &controllers.ArticleController{}, "get:ShowArticle;post:HandlerSelect")           // 文章列表
	beego.Router("/Article/addArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandlerAddArticle") // 文章添加
	beego.Router("/Article/addType", &controllers.ArticleController{}, "get:ShowAddType;post:HandlerAddType")          // 文章类型
	beego.Router("/Article/delType", &controllers.ArticleController{}, "get:HandlerDelType")                           // 文章类型删除
	beego.Router("/Article/artiContent", &controllers.ArticleController{}, "get:ShowArtiContent")                      // 文章内容
	beego.Router("/Article/artiDelete", &controllers.ArticleController{}, "get:HandlerDelete")                         // 文章删除
	beego.Router("/Article/artiUpdate", &controllers.ArticleController{}, "get:ShowArtiUpdate;post:HandlerArtiUpdate") // 文章编辑
	beego.Router("/Article/logout", &controllers.ArticleController{}, "get:Logout")                                    // 退出登录
}

var FilterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
	}
}
