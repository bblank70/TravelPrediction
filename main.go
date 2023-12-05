package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strconv"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"

	// "google.golang.org/protobuf/types/known/structpb" ///this will need to be uncommented if using .Predict vs .RawPredict
	"google.golang.org/genproto/googleapis/api/httpbody"
)

// type submission holds information recieved as a POST from /, index, home
type Submission struct {
	MonthlyIncome                 int
	Age                           int
	Passport                      int
	MaritalStatus_Divorced        int
	MaritalStatus_Married         int
	MaritalStatus_SingleUnmarried int
	PreferredPropertyStar_3       int
	PreferredPropertyStar_4       int
	PreferredPropertyStar_5       int
	ProductPitched_Basic          int
	ProductPitched_Deluxe         int
	ProductPitched_King           int
	ProductPitched_Standard       int
	ProductPitched_SuperDelux     int
}

// this struct holds information to display on /response
type ModelResult struct {
	Request string // the request that was made
	St      string //st is the body of the response
}

// //////////////////////////////////////////// variables // ////////////////////////////////////////////

// Details is the instance of the Submission struct
var Details = Submission{}

// Results is the instance of the ModelResult struct
var Results = ModelResult{}

// These are to build the request to the endpoint
var i string = ("\"instances\"")
var l string = ":[["
var cat string = "{"
var pre string = cat + i + l
var post string = "]]}"
var body string
var St string
var Requestb string

// this is the pointer to the templates
var tpl *template.Template

// //////////////////////////////////////////// functions // ////////////////////////////////////////////

// init instantiates the templates, they must be .tmpl extenstions
func init() {
	tpl = template.Must(template.ParseGlob("templates/*.tmpl"))

}

func main() {

	// these are our paths
	http.HandleFunc("/", index)
	http.HandleFunc("/verify", verifyer)
	http.HandleFunc("/response", responder)
	//this starts the server
	http.ListenAndServe(":8080", nil)

}

// ////////////////////////////////////////////

// index handles the route to the home page
func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "home.tmpl", nil)
}

// verifyer gets the values from the form, extracts and restructures them to make a post request to the vertex AI endpoint
func verifyer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	income, _ := strconv.Atoi(r.FormValue("MonthlyIncome"))
	age, _ := strconv.Atoi(r.FormValue("Age"))
	passport, _ := strconv.Atoi(r.FormValue("Passport"))
	ms_div, _ := strconv.Atoi(r.FormValue("MaritalStatus_Divorced"))
	ms_mar, _ := strconv.Atoi(r.FormValue("MaritalStatus_Married"))
	ms_su, _ := strconv.Atoi(r.FormValue("MaritalStatus_SingleUnmarried"))
	ps3, _ := strconv.Atoi(r.FormValue("PreferredPropertyStar_3"))
	ps4, _ := strconv.Atoi(r.FormValue("PreferredPropertyStar_4"))
	ps5, _ := strconv.Atoi(r.FormValue("PreferredPropertyStar_5"))
	ppb, _ := strconv.Atoi(r.FormValue("Basic"))
	ppd, _ := strconv.Atoi(r.FormValue("Delux"))
	ppk, _ := strconv.Atoi(r.FormValue("King"))
	pps, _ := strconv.Atoi(r.FormValue("Standard"))
	ppsd, _ := strconv.Atoi(r.FormValue("SuperDelux"))

	//create our struct from the form values
	Details = Submission{
		MonthlyIncome:                 income,
		Age:                           age,
		Passport:                      passport,
		MaritalStatus_Divorced:        ms_div,
		MaritalStatus_Married:         ms_mar,
		MaritalStatus_SingleUnmarried: ms_su,
		PreferredPropertyStar_3:       ps3,
		PreferredPropertyStar_4:       ps4,
		PreferredPropertyStar_5:       ps5,
		ProductPitched_Basic:          ppb,
		ProductPitched_Deluxe:         ppd,
		ProductPitched_King:           ppk,
		ProductPitched_Standard:       pps,
		ProductPitched_SuperDelux:     ppsd,
	}

	v := reflect.ValueOf(Details)
	body = ""
	// typeOfdetails := v.Type() //debugging code to get the type of v

	for i := 0; i < v.NumField(); i++ {
		// fmt.Printf("Field: %s\tValue: %v\n", typeOfdetails.Field(i).Name, v.Field(i).Interface()) //debugging code to build the http body.
		body = body + fmt.Sprintf("%v", v.Field(i).Interface()) + ","

	}

	if last := len(body) - 1; last >= 0 && body[last] == ',' {
		body = body[:last]
	}

	Requestb = pre + body + post
	log.Println("The request string was:", Requestb)

	// structure the body of the raw request
	Raw := &httpbody.HttpBody{}
	Raw.Data = []byte(Requestb)

	// indentify the post request using the raw body and the endpoint
	reqs := &aiplatformpb.RawPredictRequest{
		// Note  GCP Project ID and endpoint ID
		Endpoint: "projects/crafty-willow-399720/locations/us-central1/endpoints/6296606224133652480",
		HttpBody: Raw,
	}

	// ////////////////////////////////////////////
	//uncmt here if using .Predict vs .RawPredict

	// // fixes the m protostruct
	// m, err := structpb.NewValue(map[string]interface{}{
	// 	"MonthlyIncome":                 income,
	// 	"Age":                           age,
	// 	"Passport":                      passport,
	// 	"MaritalStatus_Divorced":        ms_div,
	// 	"MaritalStatus_Married":         ms_mar,
	// 	"MaritalStatus_SingleUnmarried": ms_su,
	// 	"PreferredPropertyStar_3":       ps3,
	// 	"PreferredPropertyStar_4":       ps4,
	// 	"PreferredPropertyStar_5":       ps5,
	// 	"ProductPitched_Basic":          ppb,
	// 	"ProductPitched_Deluxe":         ppd,
	// 	"ProductPitched_King":           ppk,
	// 	"ProductPitched_Standard":       pps,
	// 	"ProductPitched_SuperDelux":     ppsd,
	// })

	// reqs := &aiplatformpb.PredictRequest{
	// 	// Replace your-gcp-project to your GCP Project ID
	// 	// Notice the model text-bison@001 at the end of the endpoint
	// 	// If you want to use other model, change here
	// 	Endpoint:  "projects/crafty-willow-399720/locations/us-central1/endpoints/6296606224133652480",
	// 	Instances: []*structpb.Value{m},
	// }

	// if err != nil {
	// 	log.Println("The protobuffer failed to build:", err)
	// }

	// log.Println("The serialized message sent was:", m)

	// resp, err := C.Predict(Ctx, reqs)
	// if err != nil {
	// 	log.Fatalf("Error 4: %v", err)
	// }

	//uncmt if attempt to use the .Predict vs .RawPredict
	// ////////////////////////////////////////////

	// ////////////////////////////////////////////
	// CTX gets the credentials of the application service account
	Ctx := context.Background()
	C, err := aiplatform.NewPredictionClient(Ctx, option.WithEndpoint("us-central1-aiplatform.googleapis.com:443"))

	if err != nil {
		log.Println("Error 1:", err)
	}
	defer C.Close()

	// gets the response using the credentials of the application service account
	resp, err := C.RawPredict(Ctx, reqs)
	if err != nil {
		log.Fatalf("Error 4: %v", err)
	}
	log.Println(resp)

	// ////////////////////////////////////////////
	// // uncomment if attempting to use .Predict vs .RawPredict
	// RespMap := resp.Predictions[0].GetStructValue().AsMap()
	// Resp := resp.Predictions[0].GetStructValue()

	// ////////////////////////////////////////////

	RespString := fmt.Sprintf("%+v", resp)
	log.Println("The Response String was:", resp)

	//stores the response string from Vertex AI (in Results) so we can render it in /response page
	Results = ModelResult{
		Request: Requestb,
		St:      RespString,
	}

	tpl.ExecuteTemplate(w, "verify.tmpl", Details) // we pass the Details Submission to the template to render the fileds to be verified

}

// ////////////////////////////////////////////

// responder renders the /response page with the data present in Results
func responder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "response.tmpl", Results)
}
