package controllers

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"github.com/qor/session"
	"github.com/qor/session/manager"

	"github.com/NewTrident/iFunnyHub/app/forms"
	"github.com/NewTrident/iFunnyHub/app/models"
	"github.com/NewTrident/iFunnyHub/config"
	"github.com/NewTrident/iFunnyHub/config/utils"
	"github.com/NewTrident/iFunnyHub/db"
)

const cHalfHour = time.Duration(30) * time.Minute
const cNumbers = "1234567890"

func randSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = cNumbers[rand.Intn(len(cNumbers))]
	}
	return string(b)
}

// UserSub ...
func UserSub(c *gin.Context) {
	phone := c.PostForm("phone")
	action := c.PostForm("action")
	sub := c.PostForm("sub")
	challenge := c.PostForm("c")
	timestamp := c.PostForm("t")

	key := fmt.Sprintf("%s^%s^%s^%s^%s", config.Config.ChallengeKey, phone, sub, timestamp, action)

	h := md5.New()
	io.WriteString(h, key)

	keyHash := fmt.Sprintf("%x", h.Sum(nil))
	if challenge == keyHash {
		i, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			c.JSON(400, map[string]string{"err": "parse timestamp error"})
			return
		}
		tm := time.Unix(i, 0)

		if time.Now().UTC().Sub(tm) >= cHalfHour {
			c.JSON(400, map[string]string{"err": "time out"})
			return
		}

		user := &models.User{}
		db.DB.Where("phone = ?", phone).Find(user)
		pwd := ""
		if action == "unsub" {
			if user.ID == 0 {
				c.JSON(400, map[string]string{"err": "can't unsub"})
				return
			}

			err = user.UnSub(sub)
			if err != nil {
				c.JSON(400, map[string]string{"err": fmt.Sprintf("%v", err)})
				return
			}
		} else if action == "sub" {
			if user.ID == 0 {
				user.Phone = phone
				pwd = randSeq(6)
				user.Create(db.DB, pwd)
			}
			err = user.Sub(sub)
			if err != nil {
				c.JSON(400, map[string]string{"err": fmt.Sprintf("%v", err)})
				return
			}
		}

		db.DB.Save(user)
		c.JSON(200, map[string]string{"user": phone, "password": pwd, "sub": sub, "action": action})
		return
	}
	c.JSON(400, map[string]string{"err": "challenge error"})
	return
}

// UserEmailSetting ...
func UserEmailSetting(w http.ResponseWriter, req *http.Request) {
	var (
		user = utils.GetCurrentGuestUser(req)
		tx   = utils.GetDB(req)
	)
	if req.Method == "POST" {
		form := forms.MakeEmailForm(req)
		if form.IsValid() {
			form.Load(form)
			user.Email = form.Email
			user.First = false
			tx.Save(user)

			http.Redirect(w, req, "/", http.StatusSeeOther)
		}
		for _, err := range form.Errors() {
			manager.SessionManager.Flash(w, req, session.Message{
				Text: fmt.Sprintf("%v", err),
				Type: "error",
			})
		}

	}

	config.View.Execute(
		"set_email",
		map[string]interface{}{}, req, w,
	)
	return

}

func UserForgetPassword(c *gin.Context) {
	req := c.Request
	w := c.Writer

	if req.Method == "POST" {
		form := forms.MakeEmailForm(req)
		if form.IsValid() {
			form.Load(form)

			user := models.User{}
			if !db.DB.Where("email = ?", form.Email).Find(&user).RecordNotFound() {
				user.ConfirmToken = uuid.NewRandom().String()
				db.DB.Save(&user)
				if db.DB.Error != nil {
					manager.SessionManager.Flash(w, req, session.Message{
						Text: fmt.Sprintf("%v", "generated user confirm token error"),
						Type: "error",
					})
					goto RETURN
				}

				var buffer bytes.Buffer
				absPath, _ := filepath.Abs("app/views/email.tmpl")
				tmpl, _ := template.ParseFiles(absPath)
				tmpl.Execute(&buffer, map[string]interface{}{
					"resetPasswordLink": fmt.Sprintf(
						"http://www.ifunnyhub.com/reset_password?email=%s&token=%s",
						user.Email,
						user.ConfirmToken,
					),
				})

				_, err := config.MailSes.SendEmailHTML("admin@ifunnyhub.com", user.Email, "Retrieve your password on ifunnyhub", fmt.Sprintf(
					"http://www.ifunnyhub.com/reset_password?email=%s&token=%s",
					user.Email,
					user.ConfirmToken,
				), buffer.String())
				if err != nil {
					manager.SessionManager.Flash(w, req, session.Message{
						Text: fmt.Sprintf("%v", "sending email error"),
						Type: "error",
					})
					goto RETURN
				}
				manager.SessionManager.Flash(w, req, session.Message{
					Text: fmt.Sprintf("%v", "success"),
					Type: "success",
				})
				c.Redirect(http.StatusSeeOther, "/")
				return
			}
			manager.SessionManager.Flash(w, req, session.Message{
				Text: fmt.Sprintf("%v", "wrong email address"),
				Type: "error",
			})
			goto RETURN
		}
		for _, err := range form.Errors() {
			manager.SessionManager.Flash(w, req, session.Message{
				Text: fmt.Sprintf("%v", err),
				Type: "error",
			})
		}
	}
RETURN:
	config.View.Execute(
		"forget_password",
		map[string]interface{}{}, req, w,
	)
	return
}

// UserResetPassword ...
func UserResetPassword(c *gin.Context) {
	req := c.Request
	w := c.Writer

	if req.Method == "POST" {
		email, ok := c.GetQuery("email")
		if !ok {
			manager.SessionManager.Flash(w, req, session.Message{
				Text: "reset password error",
				Type: "error",
			})
			c.Redirect(http.StatusSeeOther, "/")
			return
		}

		token, ok := c.GetQuery("token")
		if !ok {
			manager.SessionManager.Flash(w, req, session.Message{
				Text: "reset password error",
				Type: "error",
			})
			c.Redirect(http.StatusSeeOther, "/")
			return
		}

		user := models.User{}
		db.DB.Where("email = ?", email).Find(&user)
		if db.DB.Error != nil {
			manager.SessionManager.Flash(w, req, session.Message{
				Text: "reset password error",
				Type: "error",
			})
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
		if user.ConfirmToken == "" {
			manager.SessionManager.Flash(w, req, session.Message{
				Text: "reset password error",
				Type: "error",
			})
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
		if user.ConfirmToken == token {
			form := forms.MakeResetPasswordForm(req)
			if form.IsValid() {
				form.Load(form)
				user.SetPassword(form.Password)
				user.ConfirmToken = ""
				db.DB.Save(&user)
				manager.SessionManager.Flash(w, req, session.Message{
					Text: "reset password success",
					Type: "success",
				})
				c.Redirect(http.StatusSeeOther, "/")
				return
			}
			for _, err := range form.Errors() {
				manager.SessionManager.Flash(w, req, session.Message{
					Text: fmt.Sprintf("%v", err),
					Type: "error",
				})
			}
			config.View.Execute(
				"reset_password",
				map[string]interface{}{}, req, w,
			)
			return
		}
		manager.SessionManager.Flash(w, req, session.Message{
			Text: "reset password error",
			Type: "error",
		})
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	config.View.Execute(
		"reset_password",
		map[string]interface{}{}, req, w,
	)
	return
}
