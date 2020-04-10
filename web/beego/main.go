package main

import "github.com/astaxie/beego/orm"

type User struct {
	Id int
	Name string  `orm: "size(100)"`
}

func init()  {
	//设置默认数据库
	orm.RegisterDataBase("default","mysql","root:131206@")
}
