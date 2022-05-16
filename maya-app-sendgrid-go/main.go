package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type template struct {
	Recipient           string            `json:"recipient"`
	TemplateId          string            `json:"templateId"`
	DynamicTemplateData map[string]string `json:"dynamic_template_data"`
}

func sendTemplate(w http.ResponseWriter, r *http.Request) {
	var newTemplate template

	reqBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		fmt.Fprintf(w, "Error sending the email")
	}
	json.Unmarshal(reqBody, &newTemplate)

	//Sendgrid Send Template
	m := mail.NewV3Mail()

	address := os.Getenv("COMPANY_EMAIL")
	name := os.Getenv("COMPANY_NAME")
	e := mail.NewEmail(name, address)

	m.SetFrom(e)
	m.SetTemplateID(newTemplate.TemplateId)

	p := mail.NewPersonalization()

	tos := mail.NewEmail(newTemplate.DynamicTemplateData["customer_name"], newTemplate.Recipient)
	p.AddTos(tos)

	for k, v := range newTemplate.DynamicTemplateData {
		p.SetDynamicTemplateData(k, v)
	}

	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Headers)
		json.NewEncoder(w).Encode(request)
	}

}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/template", sendTemplate).Methods("POST")
	log.Print("Listening on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
