package initial

import (
	"bufio"
	"fmt"
	"github.com/MobileCPX/IFunnyHub/models"
	"github.com/astaxie/beego/orm"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
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

func GameToSql() {
	f, err := os.Open("game")
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
		data := strings.Split(line, "####")
		var item = new(models.Item)
		o := orm.NewOrm()
		item.Title = data[1]
		item.Source = "//static.ifunnyhub.com/game/" + data[2] + "/"

		item.Position = 2
		item.Img = ""
		item.Url = data[0]

		now := time.Now()
		rand_min := rand.Int63n(525600)
		per := strconv.FormatInt(rand_min, 10)
		m, _ := time.ParseDuration("-" + per + "m")
		m1 := now.Add(m)
		starttime := m1.Format("2006-01-02 15:04:05")

		item.Create = starttime
		item.Playnum = int(rand.Int63n(65432))
		item.Like = int(rand.Int63n(int64(item.Playnum)))
		item.Dislike = int(rand.Int63n(int64(item.Like)))
		o.Insert(item)
		fmt.Println(data)
	}
}

func FileToSql() {
	f, err := os.Open("game")
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
		item.Playnum = int(rand.Int63n(65432))
		item.Like = int(rand.Int63n(int64(item.Playnum)))
		item.Dislike = int(rand.Int63n(int64(item.Like)))
		o.Insert(item)
		fmt.Println(data)
	}
}
