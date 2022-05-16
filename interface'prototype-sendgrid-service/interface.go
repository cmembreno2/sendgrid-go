package main

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Template struct {
	Recipient           string
	TemplateId          string
	DynamicTemplateData map[string]string
}

type EmailService interface {
	SendEmail()
}

func (t Template) SendEMail() {
	fmt.Println(t.TemplateId)                           //debug
	fmt.Println(t.DynamicTemplateData["customer_name"]) //debug
	fmt.Println(t.Recipient)                            //debug
	m := mail.NewV3Mail()

	address := os.Getenv("COMPANY_EMAIL")
	name := os.Getenv("COMPANY_NAME")
	e := mail.NewEmail(name, address)

	m.SetFrom(e)
	m.SetTemplateID(t.TemplateId)

	p := mail.NewPersonalization()

	tos := mail.NewEmail(t.DynamicTemplateData["customer_name"], t.Recipient)
	p.AddTos(tos)

	for k, v := range t.DynamicTemplateData {
		fmt.Println(k, v) //debug
		p.SetDynamicTemplateData(k, v)
	}

	m.AddPersonalizations(p)

	req := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	req.Method = "POST"
	var Body = mail.GetRequestBody(m)
	req.Body = Body
	response, err := sendgrid.API(req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Headers)
	}

}

func main() {
	t := Template{
		Recipient:  "cmembreno@getmaya.com",
		TemplateId: "d-2d71c010e51d4cff868b9ce15922468b",
		DynamicTemplateData: map[string]string{
			"subject":       "Thanks for joining Maya Financial!",
			"customer_name": "Carlos Membreno"}}

	t.SendEMail()
}
