package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	ok := form.Has("first_name")
	if ok {
		t.Error("shows valid when required field is empty/missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "asdf")

	//r.PostForm = postedData
	form = New(postedData)
	ok = form.Has("a")
	//form.Has("a", r)
	if ok == false {
		t.Error("form shows not valid when required field first name is in POST ")
	}

}

func TestForm_MinLenght(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("a", "e")
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLenght("a", 3)
	if form.Valid() {
		t.Error("shows valid when required field is only one character")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("field_b", "abcd")
	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)

	form.MinLenght("field_b", 3)
	if !form.Valid() {
		t.Error("shows not valid eventhough 4 characters (3needed)")
	}

	isError = form.Errors.Get("field_b")
	if isError != "" {
		t.Error("should not have an error but got one")
	}

}

func TestForm_ValidEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	postedData.Add("fieldA", "kake.com")
	form.ValidEmail("fieldA")
	if form.Valid() {
		t.Error("shows valid email when @ is missing")
	}

	postedData = url.Values{}
	postedData.Add("fieldb", "kake.kake@kake.com")

	form = New(postedData)

	form.ValidEmail("fieldb")
	if !form.Valid() {
		t.Error("shows not valid eventhough email ok (has @)")
	}

}
