package models

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/validations"
)

// Tag ...
type Tag struct {
	gorm.Model
	Name string

	Categorys []Category `l10n:"sync" gorm:"many2many:item_tags;"`
	// l10n.LocaleCreatable
}

func (tag Tag) Validate(db *gorm.DB) {
	if strings.TrimSpace(tag.Name) == "" {
		db.AddError(validations.NewError(tag, "Name", "Name can not be empty"))
	}
}
