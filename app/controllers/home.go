package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/NewTrident/iFunnyHub/app/models"
	"github.com/NewTrident/iFunnyHub/config"
	"github.com/NewTrident/iFunnyHub/config/utils"
)

// func HomeIndex(w http.ResponseWriter, req *http.Request) {
func HomeIndex(c *gin.Context) {
	req := c.Request
	w := c.Writer

	// offset := req.URL.Query().Get("offset")
	offset := 0
	limit := config.PerPage
	pageQuery := req.URL.Query().Get("page")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		page = 1
	}
	offset = page*config.PerPage - config.PerPage

	tags, ok := req.URL.Query()["tags"]

	var totalCount int

	items := []models.Item{}

	tx := utils.GetDB(req)
	if ok {
		tx.Raw(`
		SELECT
		    *
		FROM
		    items
		INNER JOIN
		    item_tags ON item_id = items.id
		WHERE items.deleted_at IS NULL AND tag_id IN (?)
		ORDER BY items.created_at desc
		LIMIT
		    ?
		OFFSET
		    ?
		`, tags, limit, offset).Scan(&items)
		// tx.Model(models.Item{}).Count(&totalCount).Order("created_at desc").Offset(offset).Limit(limit).Find(&items)
	} else {
		tx.Model(models.Item{}).Count(&totalCount).Order("created_at desc").Offset(offset).Limit(limit).Find(&items)
	}

	config.View.Execute("home_index", map[string]interface{}{
		"items":       items,
		"total_count": totalCount,
		"per_page":    config.PerPage,
		"page":        page,
	}, req, w)
}
