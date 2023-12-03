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
	"google.golang.org/protobuf/types/known/structpb"
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

///////////////////////////////////////////// variables

// Details is the instance of the Submission struct
var Details = Submission{}

// Results is the instance of the ModelResult struct
var Results = ModelResult{}

// this is a demo URL for JSON testing
// var posturl = "https://jsonplaceholder.typicode.com/posts"

// this is the GCP endpoint
var posturl = "https://us-central1-aiplatform.googleapis.com/v1/projects/crafty-willow-399720/locations/us-central1/endpoints/3122105048511807488:predict"

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

// //////////////////////////////////////////// init instantiates the templates, they must be .tmpl extenstions
func init() {
	tpl = template.Must(template.ParseGlob("templates/*.tmpl"))

}

// ////////////////////////////////////////////
func main() {

	// these are our paths
	http.HandleFunc("/", index)
	http.HandleFunc("/verify", verifyer)
	http.HandleFunc("/response", responder)
	//this starts the server
	http.ListenAndServe(":8080", nil)

}

///////////////////////////////////

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "home.tmpl", nil)
}

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

	//create our struct
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
	// typeOfdetails := v.Type()

	for i := 0; i < v.NumField(); i++ {
		// fmt.Printf("Field: %s\tValue: %v\n", typeOfdetails.Field(i).Name, v.Field(i).Interface())
		body = body + fmt.Sprintf("%v", v.Field(i).Interface()) + ","
	}

	if last := len(body) - 1; last >= 0 && body[last] == ',' {
		body = body[:last]
	}

	Requestb = pre + body + post
	fmt.Println("The request string was:", Requestb)
	// resp, err := http.Post(posturl, "application/x-www-form-urlencoded", bytes.NewBuffer(payload))

	/////////////////////from  https://medium.com/google-cloud/generative-ai-app-development-using-vertex-ai-and-golang-cf315c7fa4e1

	Ctx := context.Background()
	C, err := aiplatform.NewPredictionClient(Ctx)
	// C, err := aiplatform.NewPredictionClient(Ctx, option.WithEndpoint("us-central1-aiplatform.googleapis.com:443"))
	if err != nil {
		log.Fatalf("Error 1: %v", err)
	}
	defer C.Close()

	reqs := &aiplatformpb.PredictRequest{
		// Replace your-gcp-project to your GCP Project ID
		// Notice the model text-bison@001 at the end of the endpoint
		// If you want to use other model, change here
		Endpoint:  "projects/crafty-willow-399720/locations/us-central1/endpoints/3122105048511807488",
		Instances: []*structpb.Value{},
		/// need to reformat this object
	}
	//TODO::::: Get something to put into "Instances" above

	/////////////////////

	resp, err := C.Predict(Ctx, reqs)
	if err != nil {
		log.Fatalf("Error 4: %v", err)
	}

	// RespMap := resp.Predictions[0].GetStructValue().AsMap()

	Resp := resp.Predictions[0].GetStructValue()
	RespString := fmt.Sprintf("%+v", Resp)

	// Results = ModelResult{
	// 	Request: Requestb,
	// 	St:      fmt.Printf("resp: %v", RespMap),
	// }

	Results = ModelResult{
		Request: Requestb,
		St:      RespString,
	}

	tpl.ExecuteTemplate(w, "verify.tmpl", Details) // we pass the Details Submission to the template to render the fileds to be verified

}

///////////////////// This was the test code!

// resp, err := http.Post(posturl, "application/x-www-form-urlencoded", strings.NewReader(Requestb))
// if err != nil {
// 	fmt.Println("There was an error:", err)
// }
// defer resp.Body.Close()
// fmt.Println(resp)

// b, err := io.ReadAll(resp.Body)
// // fmt.Println("resp type is:", reflect.TypeOf(resp), "and is:", resp)

// b, err := io.ReadAll(reqs.Body) // This was the test code
// St = string(b)
// This was the test code

func responder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// TODO: otherwise put the response on the response page
	tpl.ExecuteTemplate(w, "response.tmpl", Results)
}
