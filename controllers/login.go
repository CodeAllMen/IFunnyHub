package controllers

import (
	"github.com/astaxie/beego"
)

type Login struct {
	beego.Controller
}

func (c *Login) Get() {
	if c.GetSession("userLogin") != nil {
		c.Redirect("/", 302)
	} else {
		c.Data["action"] = "/user/login/?next=" + c.GetString("url")
		c.TplName = "login.html"
		c.Layout = "index.html"
	}
}
