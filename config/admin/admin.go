package admin

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"

	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/qor/admin"
	"github.com/qor/media/media_library"
	"github.com/qor/media/oss"
	"github.com/qor/oss/s3"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/validations"

	"github.com/NewTrident/iFunnyHub/app/models"
	"github.com/NewTrident/iFunnyHub/config/admin/bindatafs"
	"github.com/NewTrident/iFunnyHub/config/auth"
	"github.com/NewTrident/iFunnyHub/db"
)

var Admin *admin.Admin

func init() {
	Admin = admin.New(&qor.Config{
		DB: db.DB,
	})
	Admin.SetAssetFS(bindatafs.AssetFS.NameSpace("admin"))
	Admin.SetSiteName("ifunny hub")
	Admin.SetAuth(auth.AdminAuthStruct{})

	// Category Management
	category := Admin.AddResource(
		&models.Category{}, &admin.Config{Menu: []string{"Category Management"}, Priority: -2},
	)

	// Tag Management
	tag := Admin.AddResource(
		&models.Tag{}, &admin.Config{Menu: []string{"Tag Management"}, Priority: -2},
	)

	// Add Item
	item := Admin.AddResource(&models.Item{}, &admin.Config{Menu: []string{"Item Management"}})
	itemImagesResource := Admin.AddResource(&models.ItemImage{}, &admin.Config{Menu: []string{"Item Management"}, Priority: -1})

	item.Meta(&admin.Meta{Name: "Category", Config: &admin.SelectOneConfig{AllowBlank: false}})
	item.Meta(&admin.Meta{Name: "Tags", Config: &admin.SelectManyConfig{SelectMode: "bottom_sheet"}})
	item.Meta(
		&admin.Meta{Name: "Seed",
			Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
				s := rand.NewSource(time.Now().Unix())
				r := rand.New(s)
				i := resource.(*models.Item)
				if i.Seed == 0 {
					i.Seed = r.Float64()
				}
			},
		},
	)
	item.Meta(&admin.Meta{Name: "MainIMG", Config: &media_library.MediaBoxConfig{
		RemoteDataResource: itemImagesResource,
		Max:                1,
	}})
	item.Meta(&admin.Meta{Name: "File", Config: &media_library.MediaBoxConfig{
		RemoteDataResource: itemImagesResource,
	}})

	item.Filter(&admin.Filter{
		Name:   "Tags",
		Config: &admin.SelectOneConfig{RemoteDataResource: tag},
	})

	item.Filter(&admin.Filter{
		Name:   "Categorys",
		Config: &admin.SelectOneConfig{RemoteDataResource: category},
	})

	item.UseTheme("grid")

	user := Admin.AddResource(&models.User{}, &admin.Config{Menu: []string{"User Management"}})
	user.Meta(&admin.Meta{Name: "Role", Config: &admin.SelectOneConfig{Collection: []string{"Admin", "Maintainer", "Member"}}})
	user.Meta(&admin.Meta{Name: "Password",
		Type:            "password",
		FormattedValuer: func(interface{}, *qor.Context) interface{} { return "" },
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			values := metaValue.Value.([]string)
			if len(values) > 0 {
				if newPassword := values[0]; newPassword != "" {
					bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
					if err != nil {
						context.DB.AddError(validations.NewError(user, "Password", "Can't encrpt password"))
						return
					}
					u := resource.(*models.User)
					u.Password = string(bcryptPassword)
				}
			}
		},
	})

	oss.Storage = s3.New(
		&s3.Config{
			AccessID:  "AKIAJLC6TBMI3AQB6WLA",
			AccessKey: "/K37fZ1BliOITIvJL0+1R9AkkE7oMQezTT9pNakh",
			Region:    "eu-west-2",
			Bucket:    "ifunny",
			// Endpoint:  "https://s3.amazonaws.com/ifnny-test",
			Endpoint: "//static.ifunnyhub.com",
			ACL:      awss3.BucketCannedACLPublicRead,
		},
	)

	item.SearchAttrs("Item.Title", "Category.Name", "Tag.Name")
	Admin.AddSearchResource(item)
}
