package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func Validation(value interface{}, rules string) error {
	data := fmt.Sprintf("%v", value)
	rule := strings.Split(rules, "|")
	for _, r := range rule {

		// required
		if r == "required" {
			if data == "" {
				return errors.New(" is required")
			}
		}

		// min length
		min := regexp.MustCompile("minlength").MatchString(r)
		if min {
			min_length := strings.Split(r, ":")[1]
			length_minimum, _ := strconv.Atoi(min_length)
			length_value := len(data)

			if length_value < length_minimum {
				return errors.New(" min length is " + min_length)
			}
		}

		// max length
		max := regexp.MustCompile("maxlength").MatchString(r)
		if max {
			max_length := strings.Split(r, ":")[1]
			length_maximum, _ := strconv.Atoi(max_length)
			length_value := len(data)

			if length_value > length_maximum {
				return errors.New(" max length is " + max_length)
			}
		}

		// is email address
		if r == "email" {
			emailRegexString1 := "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
			re := regexp.MustCompile(emailRegexString1)
			isEmail := re.MatchString(data)
			if !isEmail {
				return errors.New(" is not valid email address")
			}
		}

	}
	return nil
}

func Trim(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
