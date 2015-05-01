package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"log"
	"os"
	"strings"
)

type Params struct {
	sendgridUsername   string
	sendgridPassword   string
	tos                string
	recipients         string
	ccs                string
	ccRecipients       string
	bccs               string
	bccRecipients      string
	from               string
	fromName           string
	replyTo            string
	subject            string
	text               string
	html               string
	attachmentFilePath string
}

func main() {

	p := &Params{}
	parse_flags(p)
	err_read := godotenv.Load()
	if err_read != nil {
		log.Fatalf("error: %v", err_read)
	}
	SENDGRID_USERNAME := os.Getenv("SENDGRID_USERNAME")
	SENDGRID_PASSWORD := os.Getenv("SENDGRID_PASSWORD")

	email := sendgrid.NewMail()
	for _, to := range strings.Split(p.tos, ",") {
		email.AddTo(to)
	}
	email.SetFrom(os.Getenv("FROM"))
	email.SetFromName(os.Getenv("FROM_NAME"))
	email.SetSubject(os.Getenv("SUBJECT"))
	email.SetText(os.Getenv("TEXT"))
	// file, _ := os.OpenFile("./gif.gif", os.O_RDONLY, 0600)
	// email.AddAttachment("gif.gif", file)
	// defer file.Close()

	sg := sendgrid.NewSendGridClient(SENDGRID_USERNAME, SENDGRID_PASSWORD)
	if r := sg.Send(email); r == nil {
		fmt.Println("Email sent!")
	} else {
		fmt.Println(r)
	}
}

func parse_flags(p *Params) {
	flag.StringVar(&p.sendgridUsername, "sendgrid_username", "", "usage sendgrid_username")
	flag.StringVar(&p.sendgridPassword, "sendgrid_password", "", "usage sendgrid_password")
	flag.StringVar(&p.tos, "tos", "", "usage tos")
	flag.StringVar(&p.recipients, "recipients", "", "usage recipients")
	flag.StringVar(&p.ccs, "ccs", "", "usage ccs")
	flag.StringVar(&p.ccRecipients, "cc_recipients", "", "usage cc_recipients")
	flag.StringVar(&p.bccs, "bccs", "", "usage bccs")
	flag.StringVar(&p.bccRecipients, "bcc_recipients", "", "usage bcc_recipients")
	flag.StringVar(&p.from, "from", "", "usage from")
	flag.StringVar(&p.fromName, "from_name", "", "usage from_name")
	flag.StringVar(&p.replyTo, "reply_to", "", "usage reply_to")
	flag.StringVar(&p.subject, "subject", "", "usage subject")
	flag.StringVar(&p.text, "text", "", "usage text")
	flag.StringVar(&p.html, "html", "", "usage html")
	flag.StringVar(&p.attachmentFilePath, "attachment_file_path", "", "usage attachment_file_path")
	flag.Parse()
}
