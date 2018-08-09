package user

import (
	"github.com/MobileCPX/IFunnyHub/models"
	"github.com/astaxie/beego/orm"
)

func ValidUserByPass(name, pass string) error {
	o := orm.NewOrm()
	var user = new(models.Users)
	err := o.QueryTable("users").Filter("username", name).Filter("password", pass).One(user)
	return err
}
