package controllers

import (
	"fmt"
	"strconv"

	"github.com/MobileCPX/IFunnyHub/models/content"
	"github.com/MobileCPX/IFunnyHub/models/page"
	"github.com/astaxie/beego"
)

type Category struct {
	beego.Controller
}

func (c *Category) Get() {

	if c.GetSession("userLogin") != nil {
		c.Data["Valid"] = true
		c.Data["UserName"] = c.GetSession("userLogin")
	}
	var pagenum int
	pagenum, err := strconv.Atoi(c.GetString("page"))
	if err != nil {
		pagenum = 1
	}
	category := c.Ctx.Input.Param(":cate")
	var position int
	switch category {
	case "video":
		position = 1
	case "game":
		position = 2
	case "picture":
		position = 3
	case "ringtone":
		position = 4
	}
	data := page.GetEcahItems(position, pagenum)
	num_str, num_int := content.GetAllItemNum(position) //资源总数
	c.Data["Num"] = num_str
	c.Data["Data"] = data

	var pagetotal int
	fmt.Println(num_int, "=============")
	if num_int%20 == 0 {
		pagetotal = num_int / 20
	} else {
		pagetotal = num_int/20 + 1
	}

	if pagenum == 1 {
		c.Data["pre"] = false
	} else {
		c.Data["pre"] = true
	}
	if pagenum == pagetotal {
		c.Data["next"] = false
	} else {
		c.Data["next"] = true
	}
	c.Data["prenum"] = pagenum - 1
	c.Data["nextnum"] = pagenum + 1
	c.Data["page"] = page.GetPageNumList(pagenum, pagetotal)
	c.TplName = "category/" + category + ".html"
	c.Layout = "index.html"
}
