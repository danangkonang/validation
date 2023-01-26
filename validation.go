package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// func Validation(value interface{}, rules string, message ...string) error {
// 	data := fmt.Sprintf("%v", value)
// 	rule := strings.Split(rules, "|")
// 	for _, r := range rule {

// 		// required
// 		if r == "required" {
// 			if data == "" {
// 				return errors.New(" is required")
// 			}
// 		}

// 		// min length
// 		min := regexp.MustCompile("minlength").MatchString(r)
// 		if min {
// 			min_length := strings.Split(r, ":")[1]
// 			length_minimum, _ := strconv.Atoi(min_length)
// 			length_value := len(data)

// 			if length_value < length_minimum {
// 				return errors.New(" min length is " + min_length)
// 			}
// 		}

// 		// max length
// 		max := regexp.MustCompile("maxlength").MatchString(r)
// 		if max {
// 			max_length := strings.Split(r, ":")[1]
// 			length_maximum, _ := strconv.Atoi(max_length)
// 			length_value := len(data)

// 			if length_value > length_maximum {
// 				return errors.New(" max length is " + max_length)
// 			}
// 		}

// 		// is email address
// 		if r == "email" {
// 			emailRegexString1 := "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
// 			re := regexp.MustCompile(emailRegexString1)
// 			isEmail := re.MatchString(data)
// 			if !isEmail {
// 				return errors.New(" is not valid email address")
// 			}
// 		}

// 	}
// 	return nil
// }

// func Trim(value string) string {
// 	return strings.Join(strings.Fields(value), " ")
// }

type validationErrors struct {
	// Index   int    `json:"index,omitempty"`
	Key     string `json:"key,omitempty"`
	Message string `json:"message,omitempty"`
}

func MustValid(data interface{}) ([]*validationErrors, error) {
	typeT := reflect.TypeOf(data)
	for i := 0; i < typeT.NumField(); i++ {
		field := typeT.Field(i)
		validate := field.Tag.Get("validate")
		key := strings.Split(field.Tag.Get("json"), ",")[0]
		if validate != "" {
			rules := strings.Split(validate, ",")
			out := make([]*validationErrors, 0)
			for _, rule := range rules {
				value := reflect.ValueOf(data).FieldByName(field.Name)
				switch rule {
				case "required":
					if !isRequired(value) {
						formErr := new(validationErrors)
						formErr.Message = fmt.Sprintf("%s required", key)
						formErr.Key = key
						out = append(out, formErr)
					}
				case "alpha":
					if !isAlpha(value) {
						formErr := new(validationErrors)
						formErr.Message = fmt.Sprintf("%s only [a-zA-Z]", key)
						formErr.Key = key
						out = append(out, formErr)
					}
				case "alphanum":
					if !isAlphanum(value) {
						formErr := new(validationErrors)
						formErr.Message = fmt.Sprintf("%s only [a-zA-Z0-9]", key)
						formErr.Key = key
						out = append(out, formErr)
					}
				case "number":
					if !isNumber(value) {
						formErr := new(validationErrors)
						formErr.Message = "mush be number"
						formErr.Key = key
						out = append(out, formErr)
					}
				case "numeric":
					if !isNumeric(value) {
						formErr := new(validationErrors)
						formErr.Message = "mush be number"
						formErr.Key = key
						out = append(out, formErr)
					}
				case "email":
					if !isEmail(value) {
						formErr := new(validationErrors)
						formErr.Message = "invalid email"
						formErr.Key = key
						out = append(out, formErr)
					}
				case "latitude":
					if !isLatitude(value) {
						formErr := new(validationErrors)
						formErr.Message = "invalid latitude"
						formErr.Key = key
						out = append(out, formErr)
					}
				case "longitude":
					if !isLongitude(value) {
						formErr := new(validationErrors)
						formErr.Message = "invalid longitude"
						formErr.Key = key
						out = append(out, formErr)
					}
				}
			}
			if len(out) > 0 {
				return out, errors.New("form error")
			}
		}
	}

	return nil, nil
}
