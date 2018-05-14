package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qor/middlewares"
	qorUtils "github.com/qor/qor/utils"
	"github.com/qor/render"
	"github.com/qor/session"
	"github.com/qor/session/manager"
	"github.com/qor/sorting"

	"github.com/NewTrident/iFunnyHub/app/controllers"
	"github.com/NewTrident/iFunnyHub/app/models"
	"github.com/NewTrident/iFunnyHub/config"
	"github.com/NewTrident/iFunnyHub/config/admin"
	"github.com/NewTrident/iFunnyHub/config/admin/bindatafs"
	"github.com/NewTrident/iFunnyHub/config/auth"
	"github.com/NewTrident/iFunnyHub/config/utils"
	"github.com/NewTrident/iFunnyHub/db"
)

func ginParamWrap(handle http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request
		w := c.Writer
		ctx := context.WithValue(req.Context(), "params", c.Params)

		// use for Authorize
		ctx = context.WithValue(
			ctx,
			"next",
			url.PathEscape(fmt.Sprintf("%s%s", req.URL.Path, req.URL.Query().Encode())),
		)

		handle.ServeHTTP(w, req.WithContext(ctx))
	}
}

func main() {
	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	compileTemplate := cmdLine.Bool("compile-templates", false, "Compile Templates")
	cmdLine.Parse(os.Args[1:])

	// router
	// router := chi.NewRouter()
	router := gin.New()

	// Set Transaction
	router.Use(func(c *gin.Context) {
		req := c.Request
		var (
			tx = db.DB
		)

		ctx := context.WithValue(req.Context(), qorUtils.ContextDBName, tx)
		// next.ServeHTTP(w, req.WithContext(ctx))
		c.Request = req.WithContext(ctx)

		c.Next()
	})

	router.GET("/", controllers.HomeIndex)
	router.GET("/search", controllers.Search)
	router.GET("/category/:category", ginParamWrap(
		auth.GuestAuthority.Authorize(
			auth.ResouceCheck,
		)(http.HandlerFunc(controllers.GetByCategory))),
	)
	router.GET("/category/:category/item/:id", ginParamWrap(
		auth.GuestAuthority.Authorize(
			auth.ResouceCheck,
		)(http.HandlerFunc(controllers.ItemDetail))),
	)
	router.POST("/user_sub", controllers.UserSub)
	router.GET("/forget_password", controllers.UserForgetPassword)
	router.POST("/forget_password", controllers.UserForgetPassword)
	router.GET("/reset_password", controllers.UserResetPassword)
	router.POST("/reset_password", controllers.UserResetPassword)

	router.Any("/setting/email", ginParamWrap(
		auth.GuestAuthority.Authorize(
			auth.FristUserCheck,
		)(http.HandlerFunc(controllers.UserEmailSetting))))

	router.Static("/static", "./static")

	// fs := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()
	mux.Handle("/adminauth/", auth.AdminAuth.NewServeMux())
	mux.Handle("/auth/", auth.GuestAuth.NewServeMux())
	// mux.Handle("/static/", fs)
	mux.Handle("/", router)
	admin.Admin.MountTo("/admin", mux)

	models.InitcategoryService()

	config.View.FuncMapMaker = func(render *render.Render, req *http.Request, w http.ResponseWriter) template.FuncMap {
		funcMap := template.FuncMap{}

		funcMap["Sprintf"] = fmt.Sprintf

		funcMap["addURLQuery"] = utils.AddURLQuery

		funcMap["req"] = func() *http.Request {
			return req
		}

		funcMap["urlValues"] = func() url.Values {
			value := req.URL.Query()
			value.Del("page")

			return value
		}

		funcMap["toLower"] = strings.ToLower

		funcMap["timeFormat"] = func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		}

		funcMap["get_categories"] = func() (categories []*models.Category) {
			//utils.GetDB(req).Find(&categories)
			sorter := sorting.SortableCollection{
				PrimaryKeys: []string{"2", "3", "4", "5", "6"},
			}
			categories = models.CategoryService.All()
			sorter.Sort(categories)
			return
		}

		funcMap["is_category_select"] = func(cat string) bool {
			v := req.Context().Value("params")
			if params, ok := v.(gin.Params); ok {
				if c, ok := params.Get("category"); ok {
					return cat == c
				}
			}
			return false
		}

		funcMap["render_tags"] = func(req *http.Request) template.HTML {
			var buffer bytes.Buffer
			tags := []models.Tag{}

			v := req.Context().Value("params")
			queried := false
			if params, ok := v.(gin.Params); ok {
				if c, ok := params.Get("category"); ok {
					utils.GetDB(req).Raw(`
						WITH seleced_tags AS (
						  SELECT
							item_tags.tag_id
						  FROM
							item_tags 
						  INNER JOIN
							items ON items.id = item_tags.item_id
						  INNER JOIN 
							categories ON items.category_id = categories.id
						  WHERE
							categories.Code = ?
						)
						SELECT
						  *
						FROM
						  tags 
						WHERE id in (select tag_id from seleced_tags)`, c).Scan(&tags)
					queried = true
				}
			}
			if !queried {
				utils.GetDB(req).Find(&tags)
			}

			funcMap := template.FuncMap{
				"addURLQuery": utils.AddURLQuery,
				"string":      fmt.Sprintf,
			}

			absPath, _ := filepath.Abs("app/views/tags.tmpl")
			temp, err := template.New("tags.tmpl").Funcs(
				funcMap,
			).ParseFiles(absPath)
			if err != nil {
				log.Print(err)
				return template.HTML("")
			}

			err = temp.Execute(&buffer, map[string]interface{}{
				"tags": tags,
				"req":  req,
			})

			if err != nil {
				log.Print(err)
				return template.HTML("")
			}

			return template.HTML(buffer.String())
		}

		funcMap["flashes"] = func() []session.Message {
			return manager.SessionManager.Flashes(w, req)

		}

		funcMap["render_item"] = func(item models.Item) template.HTML {
			var buffer bytes.Buffer

			absPath, _ := filepath.Abs("app/views/item.tmpl")
			tmpl, _ := template.ParseFiles(absPath)
			tmpl.Execute(&buffer, map[string]interface{}{
				"item":     item,
				"category": item.GetCategory(),
			})
			return template.HTML(buffer.String())
		}

		funcMap["current_user"] = func() *models.User {
			return utils.GetCurrentGuestUser(req)
		}

		funcMap["render_pagenation"] = func(total, perPage, page int) utils.Pagenation {
			return utils.Pagenation{
				Total:   total,
				PerPage: perPage,
				Page:    page,
			}
		}

		funcMap["tags_string"] = func(item models.Item) string {
			var tags []models.Tag
			var buffer bytes.Buffer

			length := len(tags)
			utils.GetDB(req).Raw(`
			SELECT
			  *
			FROM
			  tags
			INNER jOIN
			  item_tags ON tags.id = item_tags.tag_id
			WHERE
			  item_tags.item_id = ?
			`, item.ID).Scan(&tags)

			for i, v := range tags {
				buffer.WriteString(v.Name)

				if length != i+1 {
					buffer.WriteString(", ")
				}
			}

			return buffer.String()
		}

		funcMap["you_may_like"] = func() template.HTML {
			var buffer bytes.Buffer
			items := []models.Item{}
			utils.GetDB(req).Order("random()").Limit(4).Find(&items)
			absPath, _ := filepath.Abs("app/views/like.tmpl")
			tmpl, _ := template.ParseFiles(absPath)
			tmpl.Execute(&buffer, map[string]interface{}{
				"items": items,
			})
			return template.HTML(buffer.String())
		}
		return funcMap
	}

	if *compileTemplate {
		bindatafs.AssetFS.Compile()
	} else {
		fmt.Printf("Listening on: %v\n", config.Config.Port)
		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.Config.Port), middlewares.Apply(mux)); err != nil {
			panic(err)
		}
	}

}
