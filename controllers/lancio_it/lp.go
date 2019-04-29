package lancio_it

import (
	"github.com/astaxie/beego"
)

type HelpHomePageControllers struct {
	beego.Controller
}

func (c *HelpHomePageControllers) Get() {
	page := c.Ctx.Input.Param(":lancio_it")
	switch page {
	case "home":
		c.TplName = "lancio_it/lp/help_index.html"
	case "unsub":
		// 检查退订状态，如果status 为空，表示第一次尝试退订，status == "wifi",表示wifi情况下退订
		// status := c.GetString("status")
		// subID := ""
		// if c.Ctx.GetCookie("user") != "" {
		// 	subID = c.Ctx.GetCookie("user")
		// 	fmt.Println("退订  subID", subID)
		// }

		// if status == "" {
		// 	c.Redirect("https://li.lpaosub.com/unsub/218?s="+subID, 302)
		// } else if status == "wifi" {
		// 	errorCode := c.GetString("errorCode")
		// 	if errorCode == "NotFindMsisdn" {
		// 		c.Data["NotFindMsisdn"] = true
		// 		c.Data["Msisdn"] = c.GetString("msisdn")
		// 	}
		// 	c.TplName = "lancio_it/unsub/msisdn_unsub.html"
		// }
		c.Redirect("https://google.com", 302)

	case "compatible":
		c.TplName = "lancio_it/lp/COMPATIBILI.html"
	case "tc":
		c.TplName = "lancio_it/lp/tnc.html"
	case "privity":
		c.TplName = "lancio_it/lp/privity.html"
	case "contatta":
		c.TplName = "lancio_it/lp/contatta.html"
	case "unsub_success":
		c.TplName = "lancio_it/lp/unsub_success.html"
	default:
		c.TplName = "404.html"
	}

}

func (c *HelpHomePageControllers) TCPage() {
	c.TplName = "lancio_it/lp/tnc.html"
}

func (c *HelpHomePageControllers) InfoPage() {
	c.TplName = "lancio_it/lp/privity.html"
}

func (c *HelpHomePageControllers) TerminaPage() {
	c.TplName = "lancio_it/lp/contatta.html"
}
