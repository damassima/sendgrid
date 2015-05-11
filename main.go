package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"log"
	"net/mail"
	"os"
	"path"
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

	email := sendgrid.NewMail()

	for _, to := range strings.Split(p.tos, ",") {
		email.AddTo(to)
	}
	recipients, _ := mail.ParseAddressList(p.recipients)
	email.AddRecipients(recipients)
	for _, cc := range strings.Split(p.ccs, ",") {
		email.AddCc(cc)
	}
	ccRecipients, _ := mail.ParseAddressList(p.ccRecipients)
	email.AddCcRecipients(ccRecipients)
	for _, bcc := range strings.Split(p.bccs, ",") {
		email.AddBcc(bcc)
	}
	bccRecipients, _ := mail.ParseAddressList(p.bccRecipients)
	email.AddBccRecipients(bccRecipients)

	email.SetFrom(p.from)
	email.SetFromName(p.fromName)
	email.SetReplyTo(p.replyTo)
	email.SetSubject(p.subject)
	email.SetText(p.text)
	email.SetHTML(p.html)
	if filepath := p.attachmentFilePath; filepath != "" {
		file, err := os.OpenFile(filepath, os.O_RDONLY, 0600)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		email.AddAttachment(path.Base(filepath), file)
		defer file.Close()
	}

	sg := sendgrid.NewSendGridClient(p.sendgridUsername, p.sendgridPassword)
	if r := sg.Send(email); r == nil {
		fmt.Printf("Email sent to %v\n", p.tos)
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
