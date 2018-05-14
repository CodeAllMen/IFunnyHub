package models

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"

	"github.com/qor/auth/auth_identity"
)

type User struct {
	gorm.Model

	Email        string `form:"email"`
	Password     string
	Role         string
	Phone        string
	SubPicture   bool
	SubVideo     bool
	SubGame      bool
	SubWallpaper bool
	SubRingtone  bool
	First        bool
	ConfirmToken string
}

func (user User) DisplayName() string {
	return user.Email
}

func (user *User) Sub(sub string) error {
	switch sub {
	case "picture":
		user.SubPicture = true
	case "video":
		user.SubVideo = true
	case "game":
		user.SubGame = true
	case "wallpaper":
		user.SubWallpaper = true
	case "ringtone":
		user.SubRingtone = true
	default:
		return errors.New("error sub category")
	}

	return nil
}

func (user *User) UnSub(sub string) error {
	switch sub {
	case "picture":
		user.SubPicture = false
	case "video":
		user.SubVideo = false
	case "game":
		user.SubGame = false
	case "wallpaper":
		user.SubWallpaper = false
	case "ringtone":
		user.SubRingtone = false
	default:
		return errors.New("error sub category")
	}

	return nil
}

func (user *User) SetPassword(pwd string) error {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(bcryptPassword)
	return nil
}

func (user *User) Create(tx *gorm.DB, pwd string) error {
	if user.ID == 0 {
		user.SetPassword(pwd)
		user.First = true
		tx.Create(user)
		return nil
	}
	return nil
}

func (user User) AfterCreate(tx *gorm.DB) error {
	if user.Email != "" {
		authIdentity := &auth_identity.AuthIdentity{}
		authIdentity.Provider = "password"
		authIdentity.UID = user.Email
		authIdentity.UserID = fmt.Sprintf("%v", user.ID)
		authIdentity.EncryptedPassword = user.Password
		t := time.Now()
		authIdentity.ConfirmedAt = &t
		tx.Create(authIdentity)
	}

	if user.Phone != "" {
		authIdentity := &auth_identity.AuthIdentity{}
		authIdentity.Provider = "password"
		authIdentity.UID = user.Phone
		authIdentity.EncryptedPassword = user.Password
		t := time.Now()
		authIdentity.ConfirmedAt = &t
		tx.Create(authIdentity)
	}

	return nil
}

func (user User) AfterUpdate(tx *gorm.DB) error {
	if user.Email != "" {
		authIdentity := &auth_identity.AuthIdentity{}
		tx.Where("uid = ?", user.Email).Find(authIdentity)
		if authIdentity.ID == 0 {
			authIdentity.Provider = "password"
			authIdentity.UID = user.Email
			authIdentity.UserID = fmt.Sprintf("%v", user.ID)
			authIdentity.EncryptedPassword = user.Password
			t := time.Now()
			authIdentity.ConfirmedAt = &t
			tx.Create(authIdentity)
		} else {
			authIdentity.Provider = "password"
			authIdentity.UID = user.Email
			authIdentity.UserID = fmt.Sprintf("%v", user.ID)
			authIdentity.EncryptedPassword = user.Password
			tx.Save(authIdentity)
		}
	}

	if user.Phone != "" {
		authIdentity := &auth_identity.AuthIdentity{}
		tx.Where("uid = ?", user.Phone).Find(authIdentity)
		if authIdentity.ID == 0 {
			authIdentity.Provider = "password"
			authIdentity.UID = user.Phone
			authIdentity.UserID = fmt.Sprintf("%v", user.ID)
			authIdentity.EncryptedPassword = user.Password
			t := time.Now()
			authIdentity.ConfirmedAt = &t
			tx.Create(authIdentity)
		} else {
			authIdentity.Provider = "password"
			authIdentity.UID = user.Phone
			authIdentity.UserID = fmt.Sprintf("%v", user.ID)
			authIdentity.EncryptedPassword = user.Password
			tx.Save(authIdentity)
		}
	}
	return nil
}

func (user User) AvailableLocales() []string {
	return []string{"en-US", "zh-CN"}
}
