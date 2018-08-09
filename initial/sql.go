package initial

import (
	"fmt"
	"time"

	_ "github.com/MobileCPX/IFunnyHub/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

func init() {
	user := beego.AppConfig.String("psqluser")
	passwd := beego.AppConfig.String("psqlpass")
	host := beego.AppConfig.String("psqlurls")
	port, err := beego.AppConfig.Int("psqlport")
	dbname := beego.AppConfig.String("psqldb")
	if nil != err {
		port = 5432
	}
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
	orm.RegisterDriver("postgres", orm.DRPostgres) // 注册驱动
	orm.RegisterDataBase("default",
		"postgres",
		fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
			user, passwd, dbname, host, port))
	orm.DefaultRowsLimit = -1
	orm.SetMaxIdleConns("default", 50)
	orm.SetMaxOpenConns("default", 1000)
	orm.DefaultTimeLoc = time.Local
	orm.RunSyncdb("default", false, true)
	orm.DefaultRowsLimit = -1

	// SqlToFile()
	// FileToSql()
}
