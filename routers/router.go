package routers

import (
	"github.com/MobileCPX/IFunnyHub/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.Index{})
	beego.Router("/category/?:cate", &controllers.Category{})
	beego.Router("/category/?:cate/item/?:id", &controllers.Item{})
	beego.Router("/login", &controllers.Login{})
	beego.Router("/logout", &controllers.Logout{})
	beego.Router("/user/?:method/", &controllers.User{})
}
