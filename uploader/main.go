package main

import (
	"fmt"

	"github.com/NewTrident/iFunnyHub/app/models"
	"github.com/NewTrident/iFunnyHub/db"
)

func main() {
	items := []models.Item{}
	db.DB.Raw("SELECT * FROM items  WHERE (manual_upload = true) ORDER BY created_at asc,position desc").Scan(&items)
	// db.DB.Raw("SELECT * FROM items ORDER BY created_at asc,position desc").Scan(&items)

	ID := 100000
	for _, item := range items {
		ID++
		db.DB.Exec("update items set id = ? where id = ?", ID, item.ID)
		if db.DB.Error != nil {
			fmt.Println("fuck!!!!!!!!!!!!!!!!: ", db.DB.Error)
			db.DB.Rollback()
			return
		}
		db.DB.Exec("update item_tags set item_id = ? where item_id = ?", ID, item.ID)
		if db.DB.Error != nil {
			fmt.Println("fuck!!!!!!!!!!!!!!!!: ", db.DB.Error)
			db.DB.Rollback()
			return
		}
	}
}
