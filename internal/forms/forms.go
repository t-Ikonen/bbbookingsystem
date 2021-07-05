package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

//Form creates custom form struct, embeds url.Values object
type Form struct {
	url.Values
	Errors errors
}

//Valid return TRUE if there are no errors (form is valid) otherwise returns FALSE
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

//New initializes a form structure
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

//Required check if required fiels has data
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be empty")
		}
	}
}

//MinLenght checks the lenght of a field to be enough
func (f *Form) MinLenght(field string, lenght int, r *http.Request) bool {
	x := r.Form.Get(field)

	if len(x) < lenght {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d long", lenght))
		return false
	}
	return true

}

//Has check if form field  is in  POST and field is empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	} else {
		return true
	}
}

//ValidEmail checks the format of email to be valid
func (f *Form) ValidEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "This must a valid email format")
	}
}
