//// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/NewTrident/iFunnyHub/app/models"
	"github.com/NewTrident/iFunnyHub/db"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/help"
	"github.com/qor/notification"
	"github.com/qor/notification/channels/database"
	"github.com/qor/publish2"
	"github.com/qor/qor"
)

/* How to run this script
   $ go run db/seeds/main.go db/seeds/seeds.go
*/

/* How to upload file
 * $ brew install s3cmd
 * $ s3cmd --configure (Refer https://github.com/theplant/qor-example)
 * $ s3cmd put local_file_path s3://qor3/
 */

var (
	DraftDB       = db.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff)
	AdminUser     *models.User
	AdminUserAuth *auth_identity.AuthIdentity
	Notification  = notification.New(&notification.Config{})
	Tables        = []interface{}{
		&auth_identity.AuthIdentity{},
		&models.User{},
		&models.Category{}, &models.Tag{},
		&models.Item{},
		&models.ItemImage{},
		&help.QorHelpEntry{},
	}
)

func main() {
	Notification.RegisterChannel(database.New(&database.Config{}))
	TruncateTables(Tables...)
	createRecords()
}

func createRecords() {
	fmt.Println("Start create sample data...")

	createAdminUsers()
	fmt.Println("--> Created admin users.")

	createHelps()
	fmt.Println("--> Created helps.")

	fmt.Println("--> Done!")
}

func createAdminUsers() {
	AdminUser = &models.User{}
	AdminUser.Email = "admin@admin.com"
	AdminUser.Role = "Admin"
	// AdminUser.Password = "$2a$10$a8AXd1q6J1lL.JQZfzXUY.pznG1tms8o.PK.tYD.Tkdfc3q7UrNX." // Password: testing
	AdminUser.Password = "$2a$10$r1C56JqTwGbyKNj.oSwLZOw0.pkb5pwd6L8ZvK4pRHRsEqOKUEtUy" // Password: 123123

	AdminUserAuth := &auth_identity.AuthIdentity{}
	AdminUserAuth.Provider = "password"
	AdminUserAuth.UID = "admin@admin.com"
	AdminUserAuth.EncryptedPassword = "$2a$10$a8AXd1q6J1lL.JQZfzXUY.pznG1tms8o.PK.tYD.Tkdfc3q7UrNX." // Password: 123123
	t := time.Now()
	AdminUserAuth.ConfirmedAt = &t

	DraftDB.Create(AdminUser)
	DraftDB.Create(AdminUserAuth)

	// Send welcome notification
	Notification.Send(&notification.Message{
		From:        AdminUser,
		To:          AdminUser,
		Title:       "Welcome To QOR Admin",
		Body:        "Welcome To QOR Admin",
		MessageType: "info",
	}, &qor.Context{DB: DraftDB})
}

func createHelps() {
	helps := map[string][]string{
		"How to setup a microsite":           []string{"micro_sites"},
		"How to create a user":               []string{"users"},
		"How to create an admin user":        []string{"users"},
		"How to handle abandoned order":      []string{"abandoned_orders", "orders"},
		"How to cancel a order":              []string{"orders"},
		"How to create a order":              []string{"orders"},
		"How to upload product images":       []string{"products", "product_images"},
		"How to create a product":            []string{"products"},
		"How to create a discounted product": []string{"products"},
		"How to create a store":              []string{"stores"},
		"How shop setting works":             []string{"shop_settings"},
		"How to setup seo settings":          []string{"seo_settings"},
		"How to setup seo for blog":          []string{"seo_settings"},
		"How to setup seo for product":       []string{"seo_settings"},
		"How to setup seo for microsites":    []string{"micro_sites", "seo_settings"},
		"How to setup promotions":            []string{"promotions"},
		"How to publish a promotion":         []string{"schedules", "promotions"},
		"How to create a publish event":      []string{"schedules", "scheduled_events"},
		"How to publish a product":           []string{"schedules", "products"},
		"How to publish a microsite":         []string{"schedules", "micro_sites"},
		"How to create a scheduled data":     []string{"schedules"},
		"How to take something offline":      []string{"schedules"},
	}

	for key, value := range helps {
		helpEntry := help.QorHelpEntry{
			Title: key,
			Body:  "Content of " + key,
			Categories: help.Categories{
				Categories: value,
			},
		}
		DraftDB.Create(&helpEntry)
	}
}

func TruncateTables(tables ...interface{}) {
	for _, table := range tables {
		if err := DraftDB.DropTableIfExists(table).Error; err != nil {
			panic(err)
		}

		DraftDB.AutoMigrate(table)
	}
}
