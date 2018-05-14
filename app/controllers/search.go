package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/NewTrident/iFunnyHub/app/models"
	"github.com/NewTrident/iFunnyHub/config"
	"github.com/NewTrident/iFunnyHub/config/utils"
)

func Search(c *gin.Context) {
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

	query, ok := req.URL.Query()["q"]
	if !ok {
		c.Redirect(301, "")
	}
	var totalCount int

	items := []models.Item{}

	tx := utils.GetDB(req)
	tx.Model(&models.Item{}).Where("tsv_content @@ plainto_tsquery(?)", query).Count(&totalCount)
	tx.Raw(`
		SELECT
		    *
		FROM
		    items
		WHERE
		    tsv_content @@ plainto_tsquery(?)
		ORDER BY items.created_at desc
		LIMIT
		    ?
		OFFSET
		    ?
		`,
		query,
		limit,
		offset,
	).Scan(&items)

	config.View.Execute("home_index", map[string]interface{}{
		"items":       items,
		"total_count": totalCount,
		"per_page":    config.PerPage,
		"page":        page,
	}, req, w)
}
