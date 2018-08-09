package controllers

import (
	"github.com/MobileCPX/IFunnyHub/models/user"
	"github.com/astaxie/beego"
)

type User struct {
	beego.Controller
}

func (c *User) Post() {
	method := c.Ctx.Input.Param(":method")
	switch method {
	case "login":
		ValidUser(c)
	}
}

func ValidUser(c *User) {
	nexturl := c.GetString("next")
	name := c.GetString("EmailOrPhone")
	pass := c.GetString("Password")
	if err := user.ValidUserByPass(name, pass); err == nil {
		c.SetSession("userLogin", name)
		if nexturl == "" {
			c.Redirect("/", 302)
		} else {
			c.Redirect(nexturl, 302)
		}
	} else {
		c.Data["action"] = "/user/login/?next=" + nexturl
		c.Data["Alert"] = true
		c.Data["Text"] = "invalid password"
		c.Data["Type"] = "error"
		c.TplName = "login.html"
		c.Layout = "index.html"
	}
}

type Logout struct {
	beego.Controller
}

func (c *Logout) Get() {
	c.DelSession("userLogin")
	c.Redirect("/", 302)
}
