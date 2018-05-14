package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bluele/gforms"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/qor/admin"
	"github.com/qor/auth"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/auth/authority"
	"github.com/qor/auth/claims"
	"github.com/qor/auth/providers/password/encryptor"
	"github.com/qor/auth/providers/password/encryptor/bcrypt_encryptor"
	"github.com/qor/auth_themes/clean"
	"github.com/qor/qor"
	"github.com/qor/roles"
	"github.com/qor/session"
	"github.com/qor/session/manager"

	"github.com/NewTrident/iFunnyHub/app/forms"
	"github.com/NewTrident/iFunnyHub/app/models"
	"github.com/NewTrident/iFunnyHub/config"
	"github.com/NewTrident/iFunnyHub/db"
)

type PasswordProvider struct {
	Encryptor encryptor.Interface
}

func (p *PasswordProvider) GetName() string {
	return "password"
}

func (p *PasswordProvider) Login(c *auth.Context) {
	c.Auth.LoginHandler(c, func(c *auth.Context) (*claims.Claims, error) {
		var (
			authInfo    auth_identity.Basic
			req         = c.Request
			tx          = c.Auth.GetDB(req)
			provider, _ = c.Provider.(*PasswordProvider)
		)

		f := req.Context().Value("form")
		form := f.(*forms.LoginForm)

		authInfo.Provider = provider.GetName()
		authInfo.UID = form.EmailOrPhone

		if tx.Model(c.Auth.AuthIdentityModel).Where(authInfo).Scan(&authInfo).RecordNotFound() {
			return nil, auth.ErrInvalidAccount
		}

		if err := provider.Encryptor.Compare(authInfo.EncryptedPassword, form.Password); err == nil {
			return authInfo.ToClaims(), err
		}

		return nil, auth.ErrInvalidPassword
	})
}

func (p *PasswordProvider) Logout(c *auth.Context) {
	c.Auth.LogoutHandler(c)
}
func (p *PasswordProvider) ConfigAuth(c *auth.Auth) {
	return
}
func (p *PasswordProvider) Register(c *auth.Context) {
	return
	// c.Auth.RegisterHandler(c, func(c *auth.Context) (*claims.Claims, error){
	// })
}
func (p *PasswordProvider) Callback(c *auth.Context) {
	return
}
func (p *PasswordProvider) ServeHTTP(c *auth.Context) {
	return
}

var (
	// Auth initialize Auth for Authentication
	AdminAuth = clean.New(&auth.Config{
		SessionStorer: &auth.SessionStorer{
			SessionName:    "_auth_session",
			SessionManager: manager.SessionManager,
			SigningMethod:  jwt.SigningMethodHS256,
		},
		ViewPaths: []string{"app/views/admin_auth"},
		URLPrefix: "adminauth",
		DB:        db.DB,
		UserModel: models.User{},
		RegisterHandler: func(c *auth.Context, authorize func(c *auth.Context) (*claims.Claims, error)) {
			return
		},
		LoginHandler: func(c *auth.Context, authorize func(c *auth.Context) (*claims.Claims, error)) {
			var (
				user = models.User{}
				req  = c.Request
				w    = c.Writer
				tx   = c.Auth.GetDB(req)
			)

			req.ParseForm()
			user.Email = strings.TrimSpace(req.Form.Get("login"))
			if tx.Where("Email = ?", user.Email).Find(&user).RecordNotFound() {
				return
			}
			if user.Role != "Admin" {
				return
			}

			cla, err := authorize(c)

			if err == nil && cla != nil {
				c.Auth.Login(w, req, cla)
				http.Redirect(w, req, "/admin", http.StatusSeeOther)
			}

			c.Auth.Config.Render.Execute(
				"auth/login",
				map[string]interface{}{}, req, w,
			)
			return
		},
	})

	// Auth initialize Auth for Authentication
	GuestAuth = auth.New(&auth.Config{
		SessionStorer: &auth.SessionStorer{
			SessionName:    "_guest_auth_session",
			SessionManager: manager.SessionManager,
			SigningMethod:  jwt.SigningMethodHS256,
		},
		ViewPaths:  []string{"app/views/guest_auth"},
		Render:     config.View,
		DB:         db.DB,
		UserModel:  models.User{},
		Redirector: auth.Redirector{RedirectBack: config.RedirectBack},
		LoginHandler: func(c *auth.Context, authorize func(c *auth.Context) (*claims.Claims, error)) {
			var (
				req = c.Request
				w   = c.Writer
			)

			if req.Method == "POST" {
				form := forms.MakeLoginForm(req)
				if !form.IsValid() {

					for _, err := range form.Errors() {
						manager.SessionManager.Flash(w, req, session.Message{
							Text: fmt.Sprintf("%v", err),
							Type: "error",
						})
					}
					c.Auth.Config.Render.Execute("auth/login",
						map[string]interface{}{}, req, w)
					return
				}
				form.Load(form)
				ctx := context.WithValue(req.Context(), "form", form)
				c.Request = req.WithContext(ctx)

				cla, err := authorize(c)
				if err != nil {
					manager.SessionManager.Flash(w, req, session.Message{
						Text: fmt.Sprintf("%v", err),
						Type: "error",
					})
					c.Auth.Config.Render.Execute("auth/login",
						map[string]interface{}{"error": err}, req, w)
					return
				}

				err = c.Auth.Login(c.Writer, c.Request, cla)
				if err != nil {
					manager.SessionManager.Flash(w, req, session.Message{
						Text: fmt.Sprintf("%v", err),
						Type: "error",
					})
					c.Auth.Config.Render.Execute("auth/login",
						map[string]interface{}{}, req, w)
					return
				}

				manager.SessionManager.Flash(w, req, session.Message{
					Text: "login successful",
					Type: "success",
				})

				user, ok := c.Auth.GetCurrentUser(req).(*models.User)
				if !ok {
					manager.SessionManager.Flash(w, req, session.Message{
						Text: "Get user Fail",
						Type: "error",
					})
					c.Auth.Config.Render.Execute("auth/login",
						map[string]interface{}{}, req, w)
					return
				}
				if user.First {
					http.Redirect(w, req, "/setting/email", http.StatusSeeOther)
					return
				}

				next := req.URL.Query().Get("next")
				if next != "" {
					http.Redirect(w, req, next, http.StatusSeeOther)
					return
				}

				http.Redirect(w, req, "/", http.StatusSeeOther)
				return
			}

			c.Auth.Config.Render.Execute(
				"auth/login",
				map[string]interface{}{"errors": gforms.Errors{"x": []string{"x"}}}, req, w,
			)
			return
		},
	})

	// Authority initialize Authority for Authorization
	Authority = authority.New(&authority.Config{
		Auth: AdminAuth,
	})

	// Authority initialize Authority for Authorization
	GuestAuthority = authority.New(&authority.Config{
		Auth: GuestAuth,
		AccessDeniedHandler: func(w http.ResponseWriter, req *http.Request) { // redirect to home page by default
			n := req.Context().Value("next")
			next, _ := n.(string)

			_, ok := GuestAuth.GetCurrentUser(req).(*models.User)

			if !ok {
				manager.SessionManager.Flash(w, req, session.Message{
					Text: "Please login first...",
					Type: "error",
				})
				http.Redirect(w, req, fmt.Sprintf("/auth/login?next=%s", next), http.StatusSeeOther)
			}

			manager.SessionManager.Flash(w, req, session.Message{
				Text: "Please subscribe first...",
				Type: "error",
			})
			http.Redirect(w, req, fmt.Sprintf("/?next=%s", next), http.StatusSeeOther)
		},
	})
)

const (
	ADMIN          = "admin"
	PICTURE        = "picture"
	VIDEO          = "video"
	GAME           = "game"
	WALLPAPER      = "wallpaper"
	RINGTONE       = "ringtone"
	ResouceCheck   = "resouce_check"
	FristUserCheck = "user"
)

func init() {
	GuestAuth.RegisterProvider(&PasswordProvider{
		Encryptor: bcrypt_encryptor.New(&bcrypt_encryptor.Config{}),
	})
	Authority.Register("logged_in_half_hour", authority.Rule{TimeoutSinceLastLogin: time.Minute * 30})

	roles.Register("admin", func(req *http.Request, currentUser interface{}) bool {
		return currentUser != nil && currentUser.(*models.User).Role == "Admin"
	})

	// roles.Register(ADMIN, func(req *http.Request, currentUser interface{}) bool {
	// 	return currentUser != nil && currentUser.(*models.User).Role == "Admin"
	// })
	roles.Register(FristUserCheck, func(req *http.Request, currentUser interface{}) bool {
		user, ok := currentUser.(*models.User)
		if ok {
			return user.First || user.Role == "Admin"
		}
		return false
	})

	roles.Register(ResouceCheck, func(req *http.Request, currentUser interface{}) bool {
		p := req.Context().Value("params")

		if currentUser == nil {
			goto FALSE
		}

		if params, ok := p.(gin.Params); ok {
			v, _ := params.Get("category")

			switch v {
			case PICTURE:
				if currentUser.(*models.User).SubPicture {
					return true
				}
			case VIDEO:
				if currentUser.(*models.User).SubVideo {
					return true
				}
			case WALLPAPER:
				if currentUser.(*models.User).SubWallpaper {
					return true
				}
			case RINGTONE:
				if currentUser.(*models.User).SubRingtone {
					return true
				}
			case GAME:
				if currentUser.(*models.User).SubGame {
					return true
				}
			default:
				goto FALSE
			}

		}

	FALSE:
		return false
	})
}

type AdminAuthStruct struct{}

func (AdminAuthStruct) LoginURL(*admin.Context) string {
	return "/adminauth/login"
}

func (AdminAuthStruct) LogoutURL(*admin.Context) string {
	return "/adminauth/logout"
}

func (AdminAuthStruct) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	currentUser, _ := AdminAuth.GetCurrentUser(c.Request).(qor.CurrentUser)
	return currentUser
}

type GuestAuthStruct struct{}

func (GuestAuthStruct) LoginURL(*admin.Context) string {
	return "/auth/login"
}

func (GuestAuthStruct) LogoutURL(*admin.Context) string {
	return "/auth/logout"
}

func (GuestAuthStruct) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	currentUser, _ := GuestAuth.GetCurrentUser(c.Request).(qor.CurrentUser)
	return currentUser
}
