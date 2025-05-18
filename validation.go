package validation

import (
	"errors"
	"fmt"
	"html/template"
	"image"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Validation struct {
	Language map[string]string
}

func New() *Validation {
	srv := &Validation{
		Language: Lang,
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

type ValidationErrorMessage struct {
	Index   string        `json:"index,omitempty"`
	Field   string        `json:"key,omitempty"`
	Message []interface{} `json:"message,omitempty"`
}

func format(s string, v interface{}) string {
	t1 := template.New("t1")
	t1 = template.Must(t1.Parse(s))
	sb := new(strings.Builder)
	t1.Execute(sb, v)
	return sb.String()
}

func (s *Validation) Validate(data interface{}) ([]*ValidationErrorMessage, error) {
	typeT := reflect.TypeOf(data)
	typeV := reflect.ValueOf(data)
	out := make([]*ValidationErrorMessage, 0)
	for i := 0; i < typeT.NumField(); i++ {
		fieldType := typeT.Field(i)
		fieldValue := typeV.Field(i)
		var key string
		if fieldType.Tag.Get("json") == "" {
			key = fieldType.Name
		} else {
			key = strings.Split(fieldType.Tag.Get("json"), ",")[0]
		}

		formErr := new(ValidationErrorMessage)
		msg := []interface{}{}
		if validate := fieldType.Tag.Get("validate"); validate != "" {
			rules := strings.Split(strings.ReplaceAll(validate, " ", ""), ",")
			for _, rule := range rules {
				value := reflect.ValueOf(data).FieldByName(fieldType.Name)
				rl := strings.Split(rule, "=")
				switch rl[0] {
				case "eqfield":
					ok, _ := reflect.TypeOf(data).FieldByName(rl[1])
					pp := strings.Split(ok.Tag.Get("json"), ",")[0]
					if !isEqualField(value, reflect.ValueOf(data).FieldByName(rl[1]).String()) {
						msg = append(msg, format(s.Language["eqfield"], pp))
					}
				case "min":
					mn, err := strconv.Atoi(rl[1])
					if err == nil {
						if !isMinimum(value, mn) {
							msg = append(msg, format(s.Language["min"], mn))
						}
					}
				case "max":
					mx, err := strconv.Atoi(rl[1])
					if err == nil {
						if !isMaximum(value, mx) {
							msg = append(msg, format(s.Language["max"], mx))
						}
					}
				case "required":
					switch value.Kind() {
					case reflect.Slice:
						if value.Len() == 0 {
							msg = append(msg, s.Language["required"])
						}
						for j := 0; j < value.Len(); j++ {
							validationErrorMessage, err := s.Validate(value.Index(j).Interface())
							if err != nil {
								for _, a := range validationErrorMessage {
									msg = append(msg, a)
								}
							}
						}
					case reflect.String:
						if !isRequired(value) {
							msg = append(msg, s.Language["required"])
						}
					case reflect.Struct:
						if value.Type() == reflect.TypeOf(time.Time{}) {
							if value.Interface().(time.Time).IsZero() {
								msg = append(msg, s.Language["required"])
							}
						}
						// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						// 	if value.Int() == 0 {
						// 		msg = append(msg, s.Language["required"])
						// 		// return fmt.Errorf("field %s is required and must not be 0", fieldType.Name)
						// 	}
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
				case "ip":
					if !isIP(value) {
						msg = append(msg, s.Language["ip"])
					}
				case "boolean":
					if !isBoolean(value) {
						msg = append(msg, s.Language["boolean"])
					}
				case "ipv4":
					if !isIPV4(value) {
						msg = append(msg, s.Language["ipv4"])
					}
				case "ipv6":
					if !isIPV6(value) {
						msg = append(msg, s.Language["ipv6"])
					}
				case "url":
					if !isURL(value) {
						msg = append(msg, s.Language["url"])
					}
				case "date":
					if !isDate(value) {
						msg = append(msg, s.Language["date"])
					}
				}
			}
		}
		if enumTag := fieldType.Tag.Get("enum"); enumTag != "" {
			allowedValues := map[string]bool{}
			for _, val := range split(enumTag, ",") {
				allowedValues[val] = true
			}
			if !allowedValues[fieldValue.String()] {
				msg = append(msg, format(s.Language["enum"], fmt.Sprintf("[%s]", enumTag)))
			}
		}
		if len(msg) > 0 {
			formErr.Message = msg
			formErr.Field = key
			out = append(out, formErr)
		}
	}
	if len(out) > 0 {
		return out, errors.New("form error")
	}
	return nil, nil
}

func (s *Validation) FileValidate(f *os.File, rules string) (*ValidationErrorMessage, error) {
	formErr := new(ValidationErrorMessage)
	msg := []interface{}{}

	for _, v := range strings.Split(rules, ",") {
		rv := strings.Split(v, "=")
		switch rv[0] {
		case "maxsize":
			fi, _ := f.Stat()
			it, _ := strconv.ParseInt(rv[1], 0, 64)
			if fi.Size() > it {
				msg = append(msg, format(s.Language["maxsize"], it))
			}
		case "minsize":
			fi, _ := f.Stat()
			it, _ := strconv.ParseInt(rv[1], 0, 64)
			if fi.Size() < it {
				msg = append(msg, format(s.Language["minsize"], it))
			}
		case "maxwidth":
			c, _, _ := image.DecodeConfig(f)
			it, _ := strconv.ParseInt(rv[1], 0, 64)
			if int64(c.Width) > it {
				msg = append(msg, format(s.Language["maxwidth"], it))
			}
		case "minwidth":
			c, _, _ := image.DecodeConfig(f)
			it, _ := strconv.ParseInt(rv[1], 0, 64)
			if int64(c.Width) < it {
				msg = append(msg, format(s.Language["minwidth"], it))
			}
		case "minhight":
			c, _, _ := image.DecodeConfig(f)
			it, _ := strconv.ParseInt(rv[1], 0, 64)
			if int64(c.Height) < it {
				msg = append(msg, format(s.Language["minhight"], it))
			}
		case "maxhight":
			c, _, _ := image.DecodeConfig(f)
			it, _ := strconv.ParseInt(rv[1], 0, 64)
			if int64(c.Height) > it {
				msg = append(msg, format(s.Language["maxhight"], it))
			}
		}
	}
	if len(msg) > 0 {
		formErr.Message = msg
		formErr.Field = f.Name()
		return formErr, errors.New("form error")
	}
	return nil, nil
}

func split(s string, delim string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if string(s[i]) == delim {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}
