package forms

import "net/url"

//Form creates custom form struct, embeds url.Values object
type Form struct {
	url.Values
	Errors error
}

//New initializes a form structure
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}
