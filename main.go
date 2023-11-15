package main

//This is a demo project for MSDS434
import (
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strconv"
)

// submission defines the information of our web app and the info that will be sent to the endpoint
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

type ModelRequest struct {
	instances []int
}

func main() {
	// var ref string = "{\"instances\":[[ 10,40,1,0,1,0,0,1,0,0,1,0,0,0]]}"
	var pre string = "{\"instances\":[["
	var post string = "]]}"
	var body string
	var request string

	//test.html is our webform
	tmpl := template.Must(template.ParseFiles("./templates/test.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}
		// the webform exports strings, so internal to go we'll do the type conversion of the
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

		//details encodes the values from webform into our submission struct
		details := Submission{
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

		// do something with details print it for now - it'll have to go to the endpoint
		_ = details
		// body := ModelRequest{
		// 	instances:  details,
		// }
		tmpl.Execute(w, struct{ Success bool }{true})
		// fmt.Printf("%+v\n", details)
		v := reflect.ValueOf(details)
		// typeOfdetails := v.Type()

		for i := 0; i < v.NumField(); i++ {
			// fmt.Printf("Field: %s\tValue: %v\n", typeOfdetails.Field(i).Name, v.Field(i).Interface())
			body = body + fmt.Sprintf("%v", v.Field(i).Interface()) + ","
		}
		if last := len(body) - 1; last >= 0 && body[last] == ',' {
			body = body[:last]
		}

		request = pre + body + post
		fmt.Println(request)
		// fmt.Println(ref)

	})

	//for now we'll put it up on the localhost at port 8080
	http.ListenAndServe(":8080", nil)
	fmt.Println("Listening on port 8080")
}
