package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qor/session"
	"github.com/qor/session/manager"

	"github.com/NewTrident/iFunnyHub/app/models"
	"github.com/NewTrident/iFunnyHub/config"
	"github.com/NewTrident/iFunnyHub/config/utils"
)

func ItemDetail(w http.ResponseWriter, req *http.Request) {
	var (
		item models.Item
		tx   = utils.GetDB(req)
	)

	v := req.Context().Value("params")
	action := req.URL.Query().Get("action")

	if params, ok := v.(gin.Params); ok {
		if c, ok := params.Get("category"); ok {
			if id, ok := params.Get("id"); ok {

				category := models.CategoryService.GetByCode(c)

				switch action {
				case "next":
					tx.Select("id").Where("id > ? AND category_id = ?", id, category.ID).Order("created_at ASC").Limit(1).First(&item)

					if item.ID == 0 {
						manager.SessionManager.Flash(w, req, session.Message{
							Text: "This is already the last one. randomly recommend one for you",
							Type: "error",
						})
					}

					http.Redirect(w, req, fmt.Sprintf("/category/%s/item/%d", c, item.ID), http.StatusFound)
					return
				case "last":
					tx.Select("id").Where("id < ? AND category_id = ?", id, category.ID).Order("created_at DESC").Limit(1).First(&item)
					http.Redirect(w, req, fmt.Sprintf("/category/%s/item/%d", c, item.ID), http.StatusFound)
					return
				}

				if !tx.Where("id = ?", id).First(&item).RecordNotFound() {
					config.View.Execute(
						fmt.Sprintf("item/%s", c),
						map[string]interface{}{"item": item}, req, w)
					return
				}
			}
		}
	}
	http.Redirect(w, req, "/", http.StatusFound)
}
