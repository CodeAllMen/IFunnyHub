package models

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	// "github.com/qor/media/oss"

	"github.com/qor/media/media_library"
	"github.com/qor/sorting"
	"github.com/qor/validations"
)

type MediaBox struct {
	media_library.MediaBox
}

func (m MediaBox) URL(styles ...string) string {
	for _, file := range m.Files {
		return file.URL(styles...)
	}
	return ""
}

// Item ...
type Item struct {
	gorm.Model
	sorting.SortingDESC

	Title string
	Tags  []Tag `l10n:"sync" gorm:"many2many:item_tags;"`
	// File    media_library.MediaLibraryStorage `sql:"type:varchar(4096)"`
	File    MediaBox //`sql:"type:varchar(4096)"`
	MainIMG MediaBox //ItemMainImg            `sql:"type:varchar(4096)"`
	Link    string
	Seed    float64

	Category   Category `l10n:"sync"`
	CategoryID uint

	ManualUpload bool
}

func (item Item) PreViewImage() string {
	if item.ManualUpload {
		return item.MainIMG.URL()
	}
	if len(item.File.Files) < 1 {
		return item.MainIMG.URL("@qor_preview")
	}
	if item.File.Files[0].IsImage() {
		return item.File.URL("@qor_preview")
	}
	return item.MainIMG.URL("@qor_preview")
}

func (item Item) Image() string {
	if len(item.File.Files) < 1 {
		return item.MainIMG.URL()
	}
	if item.File.Files[0].IsImage() {
		return item.File.Files[0].URL()
	}
	return item.MainIMG.URL()
}

func (item Item) Validate(db *gorm.DB) {
	if strings.TrimSpace(item.Title) == "" {
		db.AddError(validations.NewError(item, "Name", "Name can not be empty"))
	}

	if item.Category.ID == 0 {
		db.AddError(validations.NewError(item, "Category", "Category can not be empty"))
	}

	if len(item.File.Files) == 0 {
		db.AddError(validations.NewError(item, "File", "File can not be empty"))
	}
}

func (item Item) GetCategory() *Category {
	c := CategoryService.GetByID(item.CategoryID)
	return c
}

func (item Item) AfterCreate(tx *gorm.DB) error {
	tx.Exec("Update items SET tsv_content = to_tsvector('english', coalesce(title, '')) WHERE ID = ?", item.ID)
	return nil
}

func (item Item) AfterUpdate(tx *gorm.DB) error {
	tx.Exec("Update items SET tsv_content = to_tsvector('english', coalesce(title, '')) WHERE ID = ?", item.ID)
	return nil
}

type Storage struct {
	media_library.MediaLibraryStorage
}

func (s Storage) URL(styles ...string) string {
	if len(styles) == 0 && !isVideoFormat(s.Url) {
		styles = append(styles, "original")
	}
	if s.Url != "" && len(styles) > 0 {
		ext := path.Ext(s.Url)
		return fmt.Sprintf("%v.%v%v", strings.TrimSuffix(s.Url, ext), styles[0], ext)
	}
	return s.Url
}

type ItemImage struct {
	gorm.Model

	File         Storage `sql:"type:varchar(4096)" media_library:"url:/system/{{class}}/{{primary_key}}/{{column}}.{{extension}}`
	SelectedType string
}

func (itemImage *ItemImage) SetSelectedType(typ string) {
	itemImage.SelectedType = typ
}

func (itemImage *ItemImage) GetSelectedType() string {
	return itemImage.SelectedType
}

func (itemImage *ItemImage) ScanMediaOptions(mediaOption media_library.MediaOption) error {
	if bytes, err := json.Marshal(mediaOption); err == nil {
		return itemImage.File.Scan(bytes)
	} else {
		return err
	}
}

func (itemImage *ItemImage) GetMediaOption() (mediaOption media_library.MediaOption) {
	mediaOption.Video = itemImage.File.Video
	mediaOption.FileName = itemImage.File.FileName
	mediaOption.URL = itemImage.File.URL()
	mediaOption.OriginalURL = itemImage.File.URL("original")
	mediaOption.CropOptions = itemImage.File.CropOptions
	mediaOption.Sizes = itemImage.File.GetSizes()
	mediaOption.Description = itemImage.File.Description
	return
}

func isVideoFormat(name string) bool {
	formats := []string{".mp4", ".m4p", ".m4v", ".m4v", ".mov", ".mpeg", ".webm", ".avi", ".ogg", ".ogv"}

	ext := strings.ToLower(regexp.MustCompile(`(\?.*?$)`).ReplaceAllString(filepath.Ext(name), ""))

	for _, format := range formats {
		if format == ext {
			return true
		}
	}

	return false
}
