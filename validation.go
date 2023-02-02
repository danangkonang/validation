package validation

import (
	"errors"
	"html/template"
	"image"
	"os"
	"reflect"
	"strconv"
	"strings"
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

type ValidationErrors struct {
	Index   string   `json:"index,omitempty"`
	Key     string   `json:"key,omitempty"`
	Message []string `json:"message,omitempty"`
}

func format(s string, v interface{}) string {
	t1 := template.New("t1")
	t1 = template.Must(t1.Parse(s))
	sb := new(strings.Builder)
	t1.Execute(sb, v)
	return sb.String()
}

func (s *Validation) Validate(data interface{}) ([]*ValidationErrors, error) {
	typeT := reflect.TypeOf(data)
	out := make([]*ValidationErrors, 0)
	for i := 0; i < typeT.NumField(); i++ {
		field := typeT.Field(i)
		var key string
		if field.Tag.Get("json") == "" {
			key = field.Name
		} else {
			key = strings.Split(field.Tag.Get("json"), ",")[0]
		}
		validate := field.Tag.Get("validate")
		if validate != "" {
			rules := strings.Split(validate, ",")
			formErr := new(ValidationErrors)
			msg := []string{}
			for _, rule := range rules {
				value := reflect.ValueOf(data).FieldByName(field.Name)
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
			if len(msg) > 0 {
				formErr.Message = msg
				formErr.Key = key
				out = append(out, formErr)
			}
		}
	}
	if len(out) > 0 {
		return out, errors.New("form error")
	}
	return nil, nil
}

func (s *Validation) FileValidate(f *os.File, rules string) (*ValidationErrors, error) {
	formErr := new(ValidationErrors)
	msg := []string{}

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
		formErr.Key = f.Name()
		return formErr, errors.New("form error")
	}
	return nil, nil
}
