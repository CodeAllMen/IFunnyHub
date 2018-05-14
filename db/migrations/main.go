package main

import (
	"github.com/NewTrident/iFunnyHub/app/models"
	"github.com/NewTrident/iFunnyHub/db"
	"github.com/qor/auth/auth_identity"
)

func main() {
	// func main() {
	AutoMigrate(&models.User{})

	AutoMigrate(&models.Item{})
	AutoMigrate(&models.Category{})

	AutoMigrate(&auth_identity.AuthIdentity{})
}

func AutoMigrate(values ...interface{}) {
	for _, value := range values {
		db.DB.AutoMigrate(value)
	}
}
