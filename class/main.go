package main

import (
	_ "class/models"
	_ "class/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.AddFuncMap("ShowPrePage", HandlerPrePage)
	beego.AddFuncMap("ShowNextPage", HandlerNextPage)
	beego.Run()
}

// 页码的上一页
func HandlerPrePage(data int) int {
	pageIndex := data - 1
	return pageIndex
}

// 页码的下一页
func HandlerNextPage(data int) int {
	pageIndex := data + 1
	return pageIndex
}
