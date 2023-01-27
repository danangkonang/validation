package validation

import (
	"errors"
	"reflect"
	"strings"
)

type Validation struct {
	Language map[string]string
}

func New() *Validation {
	data := map[string]string{
		"required":  "this field is required",
		"alpha":     "this field is required",
		"alphanum":  "this field is required",
		"number":    "this field is required",
		"numeric":   "this field is required",
		"email":     "this field is required",
		"latitude":  "this field is required",
		"longitude": "this field is required",
	}
	srv := &Validation{
		Language: data,
	}
	return srv
}

func (s *Validation) SetLanguage(lang map[string]string) {
	for x, y := range s.Language {
		if m, n := lang[x]; n {
			y = m
		}
		lang[x] = y
	}
	s.Language = lang
}

type validationErrors struct {
	Key     string `json:"key,omitempty"`
	Message string `json:"message,omitempty"`
}

func (s *Validation) MustValid(data interface{}) ([]*validationErrors, error) {
	typeT := reflect.TypeOf(data)
	for i := 0; i < typeT.NumField(); i++ {
		field := typeT.Field(i)
		validate := field.Tag.Get("validate")
		key := strings.Split(field.Tag.Get("json"), ",")[0]
		if validate != "" {
			rules := strings.Split(validate, ",")
			out := make([]*validationErrors, 0)
			formErr := new(validationErrors)
			msg := []string{}
			for _, rule := range rules {
				value := reflect.ValueOf(data).FieldByName(field.Name)
				switch rule {
				case "required":
					if !isRequired(value) {
						msg = append(msg, s.Language["required"])
					}
				case "alpha":
					if !isAlpha(value) {
						msg = append(msg, s.Language["alpha"])
					}
				case "alphanum":
					if !isAlphanum(value) {
						msg = append(msg, s.Language["alphanum"])
					}
				case "number":
					if !isNumber(value) {
						msg = append(msg, s.Language["number"])
					}
				case "numeric":
					if !isNumeric(value) {
						msg = append(msg, s.Language["numeric"])
					}
				case "email":
					if !isEmail(value) {
						msg = append(msg, s.Language["email"])
					}
				case "latitude":
					if !isLatitude(value) {
						msg = append(msg, s.Language["latitude"])
					}
				case "longitude":
					if !isLongitude(value) {
						msg = append(msg, s.Language["longitude"])
					}
				}
			}
			formErr.Message = strings.Join(msg, ",")
			formErr.Key = key
			out = append(out, formErr)
			if len(out) > 0 {
				return out, errors.New("form error")
			}
		}
	}

	return nil, nil
}
