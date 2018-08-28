package models

import (
	"github.com/astaxie/beego/orm"
)

type Categories struct {
	Id        int    `orm:"column(id);pk;auto"`
	UpdatedAt string `orm:"column(updated_at);null"`
	Name      string `orm:"column(name);null"`
	Code      string `orm:"column(code);null"`
}

type Item struct {
	Id       int    `orm:"column(id);pk;auto"`
	Title    string `orm:"column(title);null"`
	Img      string `orm:"column(img);null"`
	Source   string `orm:"column(source);null"`
	Position int    `orm:"column(position);null"`
	Playnum  int    `orm:"column(playnum);null"`
	Create   string `orm:"column(create);null"`
	Like     int    `orm:"column(like);null"`
	Dislike  int    `orm:"column(dislike);null"`
	Url      string `orm:"column(url);null"`
}

type Users struct {
	Id           int    `orm:"column(id);pk;auto"`
	CreatedAt    string `orm:"column(created_at);null"`
	UpdatedAt    string `orm:"column(updated_at);null"`
	DeletedAt    string `orm:"column(deleted_at);null"`
	Password     string `orm:"column(password);null"`
	Username     string `orm:"column(username);null"`
	SubPicture   bool   `orm:"column(sub_picture);null"`
	SubVideo     bool   `orm:"column(sub_video);null"`
	SubGame      bool   `orm:"column(sub_game);null"`
	SubWallpaper bool   `orm:"column(sub_wallpaper);null"`
	SubRingtone  bool   `orm:"column(sub_ringtone);null"`
	SubSport     bool   `orm:"column(sub_sport);null"`
	SubAdtools   bool   `orm:"column(sub_adtools);null"`
	SubAdgames   bool   `orm:"column(sub_adgames);null"`
}

func init() {
	orm.RegisterModel(new(Users), new(Categories), new(Item))
}
