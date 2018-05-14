package main

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/sourcegraph/go-ses"
)

func main() {
	emailSes := ses.Config{
		Endpoint:        "https://email.us-east-1.amazonaws.com",
		AccessKeyID:     "AKIAJL4EBWN735YF2A5Q",
		SecretAccessKey: "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",
	}
	var buffer bytes.Buffer
	absPath, _ := filepath.Abs("../views/email.tmpl")
	tmpl, _ := template.ParseFiles(absPath)
	tmpl.Execute(&buffer, map[string]interface{}{
		"resetPasswordLink": fmt.Sprintf(
			"http://www.ifunnyhub.com/reset_password?email=%s&token=%s",
			"admin@ifunnyhub.com",
			"jkljljlkj",
		),
	})

	res, err := emailSes.SendEmailHTML(
		"godone@jpecomm.com",
		"cca053@gmail.com",
		"hello, world!",
	)
	if err == nil {
		fmt.Printf("Sent email: %s...\n", res)
	} else {
		fmt.Printf("Error sending email: %s\n", err)
	}
}
