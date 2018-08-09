package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strings"

	_ "github.com/MobileCPX/IFunnyHub/initial"
	_ "github.com/MobileCPX/IFunnyHub/routers"
)

var FilterUser = func(ctx *context.Context) {
	_, ok := ctx.Input.Session("userLogin").(string)
	if !ok && !strings.Contains(ctx.Request.RequestURI, "/login") && ctx.Request.RequestURI != "/" {
		ctx.Redirect(302, "/login?url="+ctx.Request.RequestURI)
	}
}

func main() {
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
	beego.Run()
}
