package config

import (
	"html/template"
	"os"

	"github.com/jinzhu/configor"
	"github.com/microcosm-cc/bluemonday"
	"github.com/qor/auth/providers/github"
	"github.com/qor/auth/providers/google"
	"github.com/qor/redirect_back"
	"github.com/qor/render"
	"github.com/qor/session/manager"
	"github.com/sourcegraph/go-ses"
)

var Config = struct {
	Port uint `default:"8157" env:"PORT"`
	DB   struct {
		Name     string `env:"DBName" default:"qor_example"`
		Adapter  string `env:"DBAdapter" default:"mysql"`
		Host     string `env:"DBHost" default:"localhost"`
		Port     string `env:"DBPort" default:"3306"`
		User     string `env:"DBUser"`
		Password string `env:"DBPassword"`
	}
	Github       github.Config
	Google       google.Config
	ChallengeKey string `default:"TIOEU(AIRG;98S7TZythdfkljjEFGZ[D0F8R"`
	Prefix       string `default:"//ifunny.s3.eu-west-2.amazonaws.com"`
}{}

var (
	PerPage      = 30
	Root         = os.Getenv("GOPATH") + "/src/github.com/NewTrident/iFunnyHub"
	View         *render.Render
	RedirectBack = redirect_back.New(&redirect_back.Config{
		SessionManager:  manager.SessionManager,
		IgnoredPrefixes: []string{"/auth", "/adminauth"},
	})
	MailSes ses.Config
)

func init() {
	if err := configor.Load(&Config, "config/database.yml", "config/smtp.yml", "config/application.yml"); err != nil {
		panic(err)
	}

	View = render.New(&render.Config{
		DefaultLayout: "ifunny",
	})

	htmlSanitizer := bluemonday.UGCPolicy()
	View.RegisterFuncMap("raw", func(str string) template.HTML {
		return template.HTML(htmlSanitizer.Sanitize(str))
	})
	MailSes = ses.Config{
		Endpoint:        "https://email.us-east-1.amazonaws.com",
		AccessKeyID:     "AKIAJL4EBWN735YF2A5Q",
		SecretAccessKey: "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",
	}

	// dialer := gomail.NewDialer(Config.SMTP.Host, Config.SMTP.Port, Config.SMTP.User, Config.SMTP.Password)
	// sender, err := dialer.Dial()

	// Mailer = mailer.New(&mailer.Config{
	// 	Sender: gomailer.New(&gomailer.Config{Sender: sender}),
	// })
}
