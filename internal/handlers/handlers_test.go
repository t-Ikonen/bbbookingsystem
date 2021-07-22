package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"snowsuite", "/snowsuite", "GET", []postData{}, http.StatusOK},
	{"frostsuite", "/frostsuite", "GET", []postData{}, http.StatusOK},
	{"northernlights", "/northernlights", "GET", []postData{}, http.StatusOK},
	{"reservation", "/reservation", "GET", []postData{}, http.StatusOK},
	{"booking", "/booking", "GET", []postData{}, http.StatusOK},

	{"post-booking", "/booking", "POST", []postData{
		{key: "endDate", value: "01/02/2021"},
		{key: "startDate", value: "02/03/2021"},
	}, http.StatusOK},
	{"post-booking-json", "/bookingjson", "POST", []postData{
		{key: "endDate", value: "01/02/2021"},
		{key: "endDate", value: "02/03/2021"},
	}, http.StatusOK},
	{"make-reservation", "/reservation", "POST", []postData{
		{key: "firstName", value: "Erkki"},
		{key: "lastName", value: "Kuikka-Kiimanen"},
		{key: "phone", value: "123-3456-1234"},
		{key: "email", value: "kul.li@kiima.fi"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := testServer.Client().Get(testServer.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			//fmt.Println("POST")
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := testServer.Client().PostForm(testServer.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		}

	}
}
