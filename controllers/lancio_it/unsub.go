package lancio_it

import (
	"github.com/astaxie/beego"
)

type LancioITUnsubController struct {
	beego.Controller
}

func (c *LancioITUnsubController) StartSub() {
	// urlParm := c.Ctx.Input.URI()
	// if len(strings.Split(urlParm, "?")) == 2 {
	// 	c.Redirect("https://li.lpaosub.com/sub/req?"+strings.Split(urlParm, "?")[1], 302)
	// } else {
	// 	c.Redirect("https://li.lpaosub.com/sub/req", 302)
	// }
	c.Redirect("https://google.com", 302)
}

func (c *LancioITUnsubController) UnsubSuccessPage() {

	c.TplName = "lancio_it/unsub/unsub_success.html"
}

func (c *LancioITUnsubController) UnsubMsisdn() {
	errCode := c.GetString("errorCode")
	c.Data["phone"] = c.GetString("phone")
	c.Data["display"] = "none"
	if errCode == "NotFind" {
		c.Data["display"] = ""
	}

	c.TplName = "lancio_it/unsub/msisdn.html"
}
