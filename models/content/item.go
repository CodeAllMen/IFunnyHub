package content

import (
	"fmt"
	"strconv"

	"github.com/MobileCPX/IFunnyHub/models"
	"github.com/MobileCPX/IFunnyHub/utils"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// 所有item数量
func GetAllItemNum(position int) (string, int) {
	var filter string
	if position == 0 {
		filter = ""
	} else {
		filter = "where position = " + strconv.Itoa(position)
	}
	o := orm.NewOrm()
	var num int
	o.Raw("select count(1) from item " + filter).QueryRow(&num)
	return utils.NumToString(num), num
}

// 通过positison获取所有item
func GetAllItemlByPosition(position, page int) (int64, []models.Item) {
	o := orm.NewOrm()
	itms := []models.Item{}
	num, err := o.QueryTable("item").Filter("position", position).Limit(20).Offset((page - 1) * 20).All(&itms)
	if err != nil {
		logs.Debug("get item by position error: ", err)
		return num, itms
	}
	return num, itms
}

// 通过position获取五个item
func GetFiveItemlByPosition(position int) (int64, []models.Item) {
	o := orm.NewOrm()
	itms := []models.Item{}
	_, item_num := GetAllItemNum(position)
	fmt.Println(item_num)
	item_start := utils.Rand(item_num)
	num, err := o.QueryTable("item").Filter("position", position).Limit(5).Offset(item_start).All(&itms)
	if err != nil {
		logs.Debug("get item by position error: ", err)
		return num, itms
	}
	return num, itms
}

// 获取单个item
func GetItemlById(id string) models.Item {
	o := orm.NewOrm()
	itms := models.Item{}
	id_int, _ := strconv.Atoi(id)

	o.QueryTable("item").Filter("id", id_int).One(&itms)
	return itms
}

// 获取附近的id
func GetItemlRound(id, cat, action string) string {
	o := orm.NewOrm()
	itms := models.Item{}
	id_int, _ := strconv.Atoi(id)

	var position int
	switch cat {
	case "video":
		position = 1
	case "game":
		position = 2
	case "picture":
		position = 3
	case "ringtone":
		position = 4
	}

	switch action {
	case "last":
		o.QueryTable("item").Filter("position", position).Filter("id__lt", id_int).OrderBy("-id").Limit(1).One(&itms)
		if itms.Id == 0 {
			o.QueryTable("item").Filter("position", position).OrderBy("-id").Limit(1).One(&itms)
		}
	case "next":
		o.QueryTable("item").Filter("position", position).Filter("id__gt", id_int).OrderBy("id").Limit(1).One(&itms)
		if itms.Id == 0 {
			o.QueryTable("item").Filter("position", position).OrderBy("id").Limit(1).One(&itms)
		}
	default:
		o.QueryTable("item").Filter("id", id_int).One(&itms)
	}
	return strconv.Itoa(itms.Id)
}

// 获取四个随机item
func GetRandomItem(position int, id string) []models.Item {
	o := orm.NewOrm()
	itms := []models.Item{}
	_, item_num := GetAllItemNum(position)
	item_start := utils.Rand(item_num)
	o.QueryTable("item").Filter("position", position).Limit(4).Offset(item_start).All(&itms)
	return itms
}
