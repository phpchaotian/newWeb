package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// 用户结构体
type User struct {
	Id       int        `orm:"pk;auto"`
	Name     string     `orm:"size(15);unique"`
	Pwd      string     `orm:"size(32)"`
	Articles []*Article `orm:"rel(m2m)"`
}

// 文章结构体
type Article struct {
	Id          int          `orm:"pk;auto"`
	Title       string       `orm:"size(50)"`
	Content     string       `orm:"size(500)"`
	Time        time.Time    `orm:"type(datetime);auto_now_add"`
	Img         string       `orm:"size(50);null"`
	Count       int          `orm:"default(0)"`
	ArticleType *ArticleType `orm:"rel(fk)"`
	Users       []*User      `orm:"reverse(many)"`
}

// 文章类型结构体
type ArticleType struct {
	Id       int
	Typename string     `orm:"size(15)"`
	Articles []*Article `orm:"reverse(many)"`
}

func init() {
	orm.RegisterDataBase("default", "mysql", "root:3edc#EDC@tcp(127.0.0.1:3306)/newClass?charset=utf8")
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	orm.RunSyncdb("default", false, true)
}
