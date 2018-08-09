package initial

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego/orm"
	"io"
	"os"
	"strings"

	"github.com/MobileCPX/IFunnyHub/models"
)

func SqlToFile() {
	f, err := os.Create("data4")

	fmt.Println(err)
	o := orm.NewOrm()
	var data []orm.Params
	o.Raw("select * from items").Values(&data)
	for i := range data {
		var line string
		tiltle := data[i]["title"].(string) + "&&"
		line += tiltle
		file := data[i]["file"].(string)
		line += ("//" + strings.Split(strings.Split(file, "\"Url\":\"//")[1], "\",\"VideoLink")[0] + "&&")

		img := data[i]["main_img"].(string)
		if img != "" {
			line += ("//" + strings.Split(strings.Split(img, "\"Url\":\"//")[1], "\",\"VideoLink")[0] + "&&")
		} else {
			line += "null&&"
		}
		time := data[i]["created_at"].(string)
		line += time[:19]
		_, err := f.WriteString(line + "\n")
		fmt.Println(err)
	}
}

func FileToSql() {
	f, err := os.Open("data4")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		line = strings.Replace(line, "\n", "", -1)
		data := strings.Split(line, "&&")
		var item = new(models.Item)
		o := orm.NewOrm()
		item.Title = data[0]
		item.Source = data[1]
		item.Img = data[2]
		if strings.Contains(item.Source, "video") {
			item.Position = 1
		}
		if strings.Contains(item.Source, "game") {
			item.Position = 2
			item.Source = strings.Split(data[1], "img")[0]
			item.Img = ""
		}
		if strings.Contains(item.Source, "picture") {
			item.Position = 3
		}
		if strings.Contains(item.Source, "ringtone") {
			item.Position = 4
		}
		if strings.Contains(item.Source, "wallpaper") {
			continue
		}
		item.Create = strings.Replace(data[3], "T", " ", -1)
		o.Insert(item)
		fmt.Println(data)
	}
}
