package models

import (
	"strings"

	"golang.org/x/sync/syncmap"

	"github.com/jinzhu/gorm"
	"github.com/qor/sorting"
	"github.com/qor/validations"

	"github.com/NewTrident/iFunnyHub/db"
)

// CategoryService ...
var CategoryService *categoryService

type categoryService struct {
	Map *syncmap.Map
}

func (c categoryService) All() []*Category {
	var categories []*Category
	c.Map.Range(
		func(key, value interface{}) bool {
			_, ok := key.(string)
			if ok {
				categories = append(categories, value.(*Category))
			}
			return true
		},
	)
	return categories

}

func (c categoryService) GetByID(id uint) *Category {

	if v, found := c.Map.Load(id); found {
		return v.(*Category)
	}

	return nil
}

func (c categoryService) GetByCode(code string) *Category {

	if v, found := c.Map.Load(code); found {
		return v.(*Category)
	}

	return nil
}

// InitcategoryServic ...
func InitcategoryService() {
	var categories []Category
	db.DB.Find(&categories)

	Map := new(syncmap.Map)

	for i := 0; i < len(categories); i++ {
		c := categories[i]
		Map.Store(c.ID, &c)
		Map.Store(c.Code, &c)
	}
	CategoryService = &categoryService{
		Map: Map,
	}
}

type Category struct {
	gorm.Model
	sorting.Sorting

	Name string
	Code string
}

func (category Category) Validate(db *gorm.DB) {
	if strings.TrimSpace(category.Name) == "" {
		db.AddError(validations.NewError(category, "Name", "Name can not be empty"))
	}

	if strings.TrimSpace(category.Code) == "" {
		db.AddError(validations.NewError(category, "Code", "Code can not be empty"))
	}
}
