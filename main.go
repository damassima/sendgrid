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

	// priority : flag > .env (file) > env

	f := &Params{}
	parse_flags(f)
	err_read := godotenv.Load()
	if err_read != nil {
		log.Fatalf("error: %v", err_read)
	}

	p := &Params{}
	merge_params(p, f)

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

func parse_flags(f *Params) {
	flag.StringVar(&f.sendgridUsername, "sendgrid_username", "", "usage sendgrid_username")
	flag.StringVar(&f.sendgridPassword, "sendgrid_password", "", "usage sendgrid_password")
	flag.StringVar(&f.tos, "tos", "", "usage tos")
	flag.StringVar(&f.recipients, "recipients", "", "usage recipients")
	flag.StringVar(&f.ccs, "ccs", "", "usage ccs")
	flag.StringVar(&f.ccRecipients, "cc_recipients", "", "usage cc_recipients")
	flag.StringVar(&f.bccs, "bccs", "", "usage bccs")
	flag.StringVar(&f.bccRecipients, "bcc_recipients", "", "usage bcc_recipients")
	flag.StringVar(&f.from, "from", "", "usage from")
	flag.StringVar(&f.fromName, "from_name", "", "usage from_name")
	flag.StringVar(&f.replyTo, "reply_to", "", "usage reply_to")
	flag.StringVar(&f.subject, "subject", "", "usage subject")
	flag.StringVar(&f.text, "text", "", "usage text")
	flag.StringVar(&f.html, "html", "", "usage html")
	flag.StringVar(&f.attachmentFilePath, "attachment_file_path", "", "usage attachment_file_path")
	flag.Parse()
}

func merge_params(p *Params, f *Params) {
	p.sendgridUsername = coalesce(f.sendgridUsername, os.Getenv("SENDGRID_USERNAME"))
	p.sendgridPassword = coalesce(f.sendgridPassword, os.Getenv("SENDGRID_PASSWORD"))
	p.tos = coalesce(f.tos, os.Getenv("TOS"))
	p.recipients = coalesce(f.recipients, os.Getenv("RECIPIENTS"))
	p.ccs = coalesce(f.ccs, os.Getenv("CCS"))
	p.ccRecipients = coalesce(f.ccRecipients, os.Getenv("CC_RECIPIENTS"))
	p.bccs = coalesce(f.bccs, os.Getenv("BCCS"))
	p.bccRecipients = coalesce(f.bccRecipients, os.Getenv("BCC_RECIPIENTS"))
	p.from = coalesce(f.from, os.Getenv("FROM"))
	p.fromName = coalesce(f.fromName, os.Getenv("FROM_NAME"))
	p.replyTo = coalesce(f.replyTo, os.Getenv("REPLY_TO"))
	p.subject = coalesce(f.subject, os.Getenv("SUBJECT"))
	p.text = coalesce(f.text, os.Getenv("TEXT"))
	p.html = coalesce(f.html, os.Getenv("HTML"))
	p.attachmentFilePath = coalesce(f.attachmentFilePath, os.Getenv("ATTACHMENT_FILE_PATH"))
}

func coalesce(v1 string, v2 string) string {
	ret := ""
	if v1 == "" {
		ret = v2
	} else {
		ret = v1
	}
	return ret
}
