package controllers

import (
	"github.com/MobileCPX/IFunnyHub/models/content"
	"github.com/MobileCPX/IFunnyHub/models/page"
	"github.com/astaxie/beego"
)

type Index struct {
	beego.Controller
}

func (c *Index) Get() {
	data := page.GetItemsIndex()
	num, _ := content.GetAllItemNum(0)
	if c.GetSession("userLogin") != nil {
		c.Data["Valid"] = true
		c.Data["UserName"] = c.GetSession("userLogin")
	}
	c.Data["Num"] = num
	c.Data["Video"] = data["Video"]
	c.Data["Ringtone"] = data["Ringtone"]
	c.Data["Picture"] = data["Picture"]
	c.Data["Game"] = data["Game"]
	c.TplName = "category/index.html"
	c.Layout = "index.html"
}
