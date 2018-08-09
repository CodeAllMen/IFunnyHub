package controllers

import (
	"github.com/MobileCPX/IFunnyHub/models/content"
	"github.com/MobileCPX/IFunnyHub/models/page"
	"github.com/astaxie/beego"
)

type Item struct {
	beego.Controller
}

func (c *Item) Get() {

	if c.GetSession("userLogin") != nil {
		c.Data["Valid"] = true
		c.Data["UserName"] = c.GetSession("userLogin")
	}

	category := c.Ctx.Input.Param(":cate")

	if c.GetString("action") != "" {
		itemid := content.GetItemlRound(c.Ctx.Input.Param(":id"), category, c.GetString("action"))
		c.Redirect("/category/"+category+"/item/"+itemid, 302)
	}

	switch category {
	case "video":
		ItemVideo(c)
	case "game":
		ItemGame(c)
	case "picture":
		ItemPicture(c)
	case "ringtone":
		ItemRingtone(c)
	}
	c.Data["last"] = "/category/" + category + "/item/" + c.Ctx.Input.Param(":id") + "?action=last"
	c.Data["next"] = "/category/" + category + "/item/" + c.Ctx.Input.Param(":id") + "?action=next"
}

func ItemVideo(c *Item) {
	data := page.VideoContent(c.Ctx.Input.Param(":id"))
	c.Data["Src"] = data["Src"]
	c.Data["Title"] = data["Title"]
	c.Data["Img"] = data["Img"]
	c.Data["Like"] = data["Like"]
	c.Data["Dislike"] = data["Dislike"]
	c.Data["Time"] = data["Time"]
	c.Data["YouMaylike"] = data["YouMaylike"]
	c.TplName = "item/video.html"
	c.Layout = "index.html"
}

func ItemGame(c *Item) {
	data := page.GameContent(c.Ctx.Input.Param(":id"))
	c.Data["Src"] = data["Src"]
	c.Data["Title"] = data["Title"]
	c.Data["Like"] = data["Like"]
	c.Data["Dislike"] = data["Dislike"]
	c.Data["Time"] = data["Time"]
	c.Data["YouMaylike"] = data["YouMaylike"]
	c.TplName = "item/game.html"
	c.Layout = "index.html"
}

func ItemPicture(c *Item) {
	data := page.PictureContent(c.Ctx.Input.Param(":id"))
	c.Data["Id"] = data["Id"]
	c.Data["Src"] = data["Src"]
	c.Data["Title"] = data["Title"]
	c.Data["Like"] = data["Like"]
	c.Data["Dislike"] = data["Dislike"]
	c.Data["Time"] = data["Time"]
	c.Data["YouMaylike"] = data["YouMaylike"]
	c.TplName = "item/picture.html"
	c.Layout = "index.html"
}

func ItemRingtone(c *Item) {
	data := page.RingtoneContent(c.Ctx.Input.Param(":id"))
	c.Data["Src"] = data["Src"]
	c.Data["Title"] = data["Title"]
	c.Data["Like"] = data["Like"]
	c.Data["Dislike"] = data["Dislike"]
	c.Data["Time"] = data["Time"]
	c.Data["YouMaylike"] = data["YouMaylike"]
	c.TplName = "item/ringtone.html"
	c.Layout = "index.html"
}
