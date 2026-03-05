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
	Language         map[string]string
	customValidators map[string]func(reflect.Value) bool
}

func New() *Validation {
	srv := &Validation{
		Language:         Lang,
		customValidators: make(map[string]func(reflect.Value) bool),
	}
	return srv
}

func (s *Validation) SetLanguage(lang map[string]string) {
	if lang == nil {
		return
	}
	for key, msg := range lang {
		s.Language[key] = msg
	}
}

func (s *Validation) RegisterValidation(tag string, fn func(reflect.Value) bool) {
	if s.customValidators == nil {
		s.customValidators = make(map[string]func(reflect.Value) bool)
	}
	s.customValidators[tag] = fn
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
	if data == nil {
		return nil, errors.New("input data cannot be nil")
	}
	typeV := reflect.ValueOf(data)
	if typeV.Kind() == reflect.Pointer {
		if typeV.IsNil() {
			return nil, errors.New("input data cannot be nil")
		}
		typeV = typeV.Elem()
	}
	if typeV.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct")
	}
	typeT := reflect.TypeOf(data)
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
				if rule == "" {
					continue
				}
				value := reflect.ValueOf(data).FieldByName(fieldType.Name)
				hasOmitEmpty := false
				for _, r := range rules {
					if r == "omitempty" {
						hasOmitEmpty = true
						break
					}
				}
				if hasOmitEmpty && isEmptyValue(value) {
					// skip semua validation untuk field ini
					break
				}

				rl := strings.Split(rule, "=")
				// if len(rl) != 2 {
				// 	continue
				// }
				ruleName := rl[0]
				switch ruleName {
				case "eqfield":
					if len(rl) != 2 {
						continue
					}
					ok, _ := reflect.TypeOf(data).FieldByName(rl[1])
					pp := strings.Split(ok.Tag.Get("json"), ",")[0]
					if !isEqualField(value, reflect.ValueOf(data).FieldByName(rl[1]).String()) {
						msg = append(msg, format(s.Language["eqfield"], pp))
					}
				case "min":
					if len(rl) != 2 {
						continue
					}
					mn, err := strconv.Atoi(rl[1])
					if err == nil {
						if !isMinimum(value, mn) {
							msg = append(msg, format(s.Language["min"], mn))
						}
					}
				case "max":
					if len(rl) != 2 {
						continue
					}
					mx, err := strconv.Atoi(rl[1])
					if err == nil {
						if !isMaximum(value, mx) {
							msg = append(msg, format(s.Language["max"], mx))
						}
					}
				case "len":
					if len(rl) != 2 {
						continue
					}
					l, err := strconv.Atoi(rl[1])
					if err != nil {
						continue
					}
					if !isLength(value, l) {
						msg = append(msg, format(s.Language["len"], l))
					}

				case "gte":
					if len(rl) != 2 {
						continue
					}
					val, err := strconv.ParseFloat(rl[1], 64)
					if err != nil {
						continue
					}
					if !isGTE(value, val) {
						msg = append(msg, format(s.Language["gte"], rl[1]))
					}

				case "lte":
					if len(rl) != 2 {
						continue
					}
					val, err := strconv.ParseFloat(rl[1], 64)
					if err != nil {
						continue
					}
					if !isLTE(value, val) {
						msg = append(msg, format(s.Language["lte"], rl[1]))
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
				case "timezone":
					if !isTimezone(value) {
						msg = append(msg, s.Language["timezone"])
					}
				case "omitempty":
					continue
				default:
					if fn, ok := s.customValidators[ruleName]; ok {
						if !fn(value) {
							if customMsg, exists := s.Language[ruleName]; exists {
								msg = append(msg, customMsg)
							} else {
								msg = append(msg, fmt.Sprintf("Invalid %s", ruleName))
							}
						}
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
