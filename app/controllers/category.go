package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/NewTrident/iFunnyHub/app/models"
	"github.com/NewTrident/iFunnyHub/config"
	"github.com/NewTrident/iFunnyHub/config/utils"
)

func GetByCategory(w http.ResponseWriter, req *http.Request) {
	var (
		category models.Category
		items    []models.Item
		tx       = utils.GetDB(req)
	)
	offset := 0
	limit := config.PerPage
	pageQuery := req.URL.Query().Get("page")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		page = 1
	}
	offset = page*config.PerPage - config.PerPage

	tags, tags_ok := req.URL.Query()["tags"]

	var totalCount int

	v := req.Context().Value("params")

	if params, ok := v.(gin.Params); ok {
		if c, ok := params.Get("category"); ok {
			if !tx.Where("code = ?", c).First(&category).RecordNotFound() {
				if tags_ok {
					tx.Model(&models.Item{}).Joins(
						"JOIN item_tags ON item_tags.item_id = items.id AND item_tags.tag_id in (?)", tags,
					).Where("category_id = ?", category.ID).Count(&totalCount)

					tx.Where(
						&models.Item{CategoryID: category.ID},
					).Joins(
						"JOIN item_tags ON item_tags.item_id = items.id AND item_tags.tag_id in (?)", tags,
					).Order("created_at desc").Offset(offset).Limit(limit).Find(&items)

					config.View.Execute(
						fmt.Sprintf("category/%s", c),
						map[string]interface{}{
							"per_page":    config.PerPage,
							"total_count": totalCount,
							"items":       items,
							"page":        page,
						}, req, w)
					return
				}

				tx.Model(&models.Item{}).Where("category_id = ?", category.ID).Count(&totalCount)

				tx.Where(
					&models.Item{CategoryID: category.ID},
				).Order("created_at desc").Offset(offset).Limit(limit).Find(&items)

				config.View.Execute(
					fmt.Sprintf("category/%s", c),
					map[string]interface{}{
						"per_page":    config.PerPage,
						"total_count": totalCount,
						"items":       items,
						"page":        page,
					}, req, w)
				return
			}
		}
	}

	http.Redirect(w, req, "/", http.StatusFound)
}
