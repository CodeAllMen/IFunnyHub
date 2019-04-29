package routers

import (
	"github.com/MobileCPX/IFunnyHub/controllers"
	"github.com/MobileCPX/IFunnyHub/controllers/lancio_it"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.Index{})
	beego.Router("/category/?:cate", &controllers.Category{})
	beego.Router("/category/?:cate/item/?:id", &controllers.Item{})
	beego.Router("/login", &controllers.Login{})
	beego.Router("/logout", &controllers.Logout{})
	beego.Router("/user/?:method/", &controllers.User{})

	// Lancio 意大利服务
	beego.Router("/ifunny/?:lancio_it", &lancio_it.HelpHomePageControllers{})

	// H3G条款
	beego.Router("/TC", &lancio_it.HelpHomePageControllers{}, "Get:TCPage")
	beego.Router("/INFO", &lancio_it.HelpHomePageControllers{}, "Get:InfoPage")
	beego.Router("/TERMINALI", &lancio_it.HelpHomePageControllers{}, "Get:TerminaPage")

	beego.Router("/success/unsub", &lancio_it.LancioITUnsubController{}, "Get:UnsubSuccessPage")
	beego.Router("/unsub/wifi", &lancio_it.LancioITUnsubController{}, "Get:UnsubMsisdn")

	beego.Router("/sub/req", &lancio_it.LancioITUnsubController{}, "Get:StartSub")
}
