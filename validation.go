package validation

import (
	"errors"
	"fmt"
	"html/template"
	"image"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrValidation   = errors.New("form error")
	ErrInvalidInput = errors.New("invalid validation input")
)

type InvalidRuleError struct{ Field, Rule string }

func (e *InvalidRuleError) Error() string {
	return fmt.Sprintf("invalid rule %q on field %q", e.Rule, e.Field)
}

type Validation struct {
	Language         map[string]string
	customValidators map[string]func(reflect.Value) bool
}

func New() *Validation {
	language := make(map[string]string, len(Lang))
	for key, message := range Lang {
		language[key] = message
	}
	return &Validation{Language: language, customValidators: make(map[string]func(reflect.Value) bool)}
}

func (s *Validation) SetLanguage(lang map[string]string) {
	if lang == nil {
		return
	}
	if s.Language == nil {
		s.Language = make(map[string]string)
	}
	for key, message := range lang {
		s.Language[key] = message
	}
}

func (s *Validation) RegisterValidation(tag string, fn func(reflect.Value) bool) error {
	tag = strings.TrimSpace(tag)
	if tag == "" || fn == nil {
		return errors.New("custom validation tag and callback are required")
	}
	if s.customValidators == nil {
		s.customValidators = make(map[string]func(reflect.Value) bool)
	}
	s.customValidators[tag] = fn
	return nil
}

type ValidationErrorMessage struct {
	Index   string        `json:"index,omitempty"`
	Field   string        `json:"key,omitempty"`
	Message []interface{} `json:"message,omitempty"`
}

func format(message string, value interface{}) (string, error) {
	tmpl, err := template.New("validation-message").Parse(message)
	if err != nil {
		return "", err
	}
	var builder strings.Builder
	if err := tmpl.Execute(&builder, value); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func (s *Validation) Validate(data interface{}) ([]*ValidationErrorMessage, error) {
	value, err := indirect(reflect.ValueOf(data))
	if err != nil {
		return nil, err
	}
	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: input must be a struct", ErrInvalidInput)
	}
	result, err := s.validateStruct(value, "", 0)
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result, ErrValidation
	}
	return nil, nil
}

func indirect(value reflect.Value) (reflect.Value, error) {
	if !value.IsValid() {
		return reflect.Value{}, fmt.Errorf("%w: input data cannot be nil", ErrInvalidInput)
	}
	for value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return reflect.Value{}, fmt.Errorf("%w: input data cannot be nil", ErrInvalidInput)
		}
		value = value.Elem()
	}
	return value, nil
}

func (s *Validation) validateStruct(value reflect.Value, prefix string, depth int) ([]*ValidationErrorMessage, error) {
	if depth > 100 {
		return nil, errors.New("validation nesting exceeds 100 levels")
	}
	result := make([]*ValidationErrorMessage, 0)
	typeValue := value.Type()
	for i := 0; i < value.NumField(); i++ {
		fieldType := typeValue.Field(i)
		if fieldType.PkgPath != "" {
			continue
		}
		fieldValue := value.Field(i)
		key := jsonName(fieldType)
		if key == "-" {
			continue
		}
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}
		messages, skip, err := s.validateField(fieldType, fieldValue, value)
		if err != nil {
			return nil, err
		}
		if len(messages) > 0 {
			result = append(result, &ValidationErrorMessage{Field: path, Message: messages})
		}
		if skip {
			continue
		}
		nested, err := s.validateNested(fieldValue, path, depth+1)
		if err != nil {
			return nil, err
		}
		result = append(result, nested...)
	}
	return result, nil
}

func (s *Validation) validateField(fieldType reflect.StructField, value, parent reflect.Value) ([]interface{}, bool, error) {
	raw := strings.TrimSpace(fieldType.Tag.Get("validate"))
	if raw == "" {
		return s.appendEnum(nil, fieldType, value, false)
	}
	rules, err := parseRules(raw)
	if err != nil {
		return nil, false, &InvalidRuleError{fieldType.Name, err.Error()}
	}
	if hasRule(rules, "omitempty") && isEmptyValue(value) {
		return nil, true, nil
	}
	messages := make([]interface{}, 0)
	for _, rule := range rules {
		messageKey, argument, valid, err := s.applyRule(rule, value, fieldType, parent)
		if err != nil {
			return nil, false, err
		}
		if !valid {
			templateMessage := s.Language[messageKey]
			if templateMessage == "" {
				templateMessage = fmt.Sprintf("Invalid %s", messageKey)
			}
			message, err := format(templateMessage, argument)
			if err != nil {
				return nil, false, fmt.Errorf("format message for %s: %w", messageKey, err)
			}
			messages = append(messages, message)
		}
	}
	return s.appendEnum(messages, fieldType, value, false)
}

func (s *Validation) applyRule(rule string, value reflect.Value, fieldType reflect.StructField, parent reflect.Value) (string, interface{}, bool, error) {
	name, argument, err := splitRule(rule)
	if err != nil {
		return "", nil, false, &InvalidRuleError{fieldType.Name, rule}
	}
	switch name {
	case "omitempty":
		return name, nil, true, nil
	case "required":
		return name, nil, isRequired(value), nil
	case "min", "max", "len":
		n, err := strconv.Atoi(argument)
		if err != nil {
			return "", nil, false, &InvalidRuleError{fieldType.Name, rule}
		}
		if name == "min" {
			return name, n, isMinimum(value, n), nil
		}
		if name == "max" {
			return name, n, isMaximum(value, n), nil
		}
		return name, n, isLength(value, n), nil
	case "gte", "lte":
		n, err := strconv.ParseFloat(argument, 64)
		if err != nil {
			return "", nil, false, &InvalidRuleError{fieldType.Name, rule}
		}
		if name == "gte" {
			return name, argument, isGTE(value, n), nil
		}
		return name, argument, isLTE(value, n), nil
	case "eqfield":
		other := parent.FieldByName(argument)
		if !other.IsValid() {
			return "", nil, false, &InvalidRuleError{fieldType.Name, rule}
		}
		otherName := argument
		if field, ok := parent.Type().FieldByName(argument); ok {
			otherName = jsonName(field)
		}
		return name, otherName, reflect.DeepEqual(indirectValue(value), indirectValue(other)), nil
	case "alpha", "alphanum", "number", "numeric", "email", "latitude", "longitude", "ip", "boolean", "ipv4", "ipv6", "url", "date", "timezone":
		return name, nil, runBuiltin(name, value), nil
	default:
		fn, ok := s.customValidators[name]
		if !ok {
			return "", nil, false, &InvalidRuleError{fieldType.Name, rule}
		}
		valid := false
		func() { defer func() { _ = recover() }(); valid = fn(value) }()
		return name, nil, valid, nil
	}
}

func indirectValue(value reflect.Value) interface{} {
	for value.IsValid() && (value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if !value.IsValid() || !value.CanInterface() {
		return nil
	}
	return value.Interface()
}

func (s *Validation) validateNested(value reflect.Value, prefix string, depth int) ([]*ValidationErrorMessage, error) {
	value, err := indirectNested(value)
	if err != nil {
		return nil, nil
	}
	switch value.Kind() {
	case reflect.Struct:
		if value.Type() == reflect.TypeOf(time.Time{}) {
			return nil, nil
		}
		return s.validateStruct(value, prefix, depth)
	case reflect.Slice, reflect.Array:
		result := make([]*ValidationErrorMessage, 0)
		for i := 0; i < value.Len(); i++ {
			nested, err := s.validateNested(value.Index(i), fmt.Sprintf("%s[%d]", prefix, i), depth+1)
			if err != nil {
				return nil, err
			}
			result = append(result, nested...)
		}
		return result, nil
	case reflect.Map:
		result := make([]*ValidationErrorMessage, 0)
		for _, key := range value.MapKeys() {
			nested, err := s.validateNested(value.MapIndex(key), fmt.Sprintf("%s[%v]", prefix, key.Interface()), depth+1)
			if err != nil {
				return nil, err
			}
			result = append(result, nested...)
		}
		return result, nil
	}
	return nil, nil
}

func indirectNested(value reflect.Value) (reflect.Value, error) {
	for value.IsValid() && (value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return reflect.Value{}, errors.New("nil")
		}
		value = value.Elem()
	}
	return value, nil
}

func (s *Validation) appendEnum(messages []interface{}, fieldType reflect.StructField, value reflect.Value, skip bool) ([]interface{}, bool, error) {
	if skip {
		return messages, true, nil
	}
	tag := strings.TrimSpace(fieldType.Tag.Get("enum"))
	if tag == "" || isEmptyValue(value) {
		return messages, false, nil
	}
	if value.Kind() != reflect.String {
		return nil, false, &InvalidRuleError{fieldType.Name, "enum"}
	}
	allowed := false
	for _, item := range strings.Split(tag, ",") {
		if strings.TrimSpace(item) == value.String() {
			allowed = true
			break
		}
	}
	if !allowed {
		message, err := format(s.Language["enum"], "["+tag+"]")
		if err != nil {
			return nil, false, err
		}
		messages = append(messages, message)
	}
	return messages, false, nil
}

func parseRules(raw string) ([]string, error) {
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, errors.New("empty rule")
		}
		result = append(result, part)
	}
	return result, nil
}

func splitRule(rule string) (string, string, error) {
	parts := strings.SplitN(rule, "=", 2)
	name := strings.TrimSpace(parts[0])
	if name == "" {
		return "", "", errors.New("empty rule")
	}
	argument := ""
	if len(parts) == 2 {
		argument = strings.TrimSpace(parts[1])
	}
	needsArgument := name == "min" || name == "max" || name == "len" || name == "gte" || name == "lte" || name == "eqfield" || name == "maxsize" || name == "minsize" || name == "maxwidth" || name == "minwidth" || name == "maxheight" || name == "minheight" || name == "maxhight" || name == "minhight"
	if needsArgument && argument == "" {
		return "", "", errors.New("missing argument")
	}
	if !needsArgument && len(parts) == 2 {
		return "", "", errors.New("unexpected argument")
	}
	return name, argument, nil
}

func hasRule(rules []string, target string) bool {
	for _, rule := range rules {
		name, _, _ := splitRule(rule)
		if name == target {
			return true
		}
	}
	return false
}
func jsonName(field reflect.StructField) string {
	name := strings.Split(field.Tag.Get("json"), ",")[0]
	if name == "" {
		return field.Name
	}
	return name
}

func (s *Validation) FileValidate(file *os.File, rules string) (*ValidationErrorMessage, error) {
	if file == nil {
		return nil, fmt.Errorf("%w: file cannot be nil", ErrInvalidInput)
	}
	parts, err := parseRules(strings.TrimSpace(rules))
	if err != nil {
		return nil, &InvalidRuleError{file.Name(), err.Error()}
	}
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}
	messages := make([]interface{}, 0)
	for _, rule := range parts {
		name, argument, err := splitRule(rule)
		if err != nil {
			return nil, &InvalidRuleError{file.Name(), rule}
		}
		limit, err := strconv.ParseInt(argument, 10, 64)
		if err != nil || limit < 0 {
			return nil, &InvalidRuleError{file.Name(), rule}
		}
		failed := false
		switch name {
		case "maxsize":
			failed = info.Size() > limit
		case "minsize":
			failed = info.Size() < limit
		case "maxwidth", "minwidth", "maxheight", "minheight", "maxhight", "minhight":
			config, err := decodeConfig(file)
			if err != nil {
				return nil, fmt.Errorf("decode image: %w", err)
			}
			switch name {
			case "maxwidth":
				failed = int64(config.Width) > limit
			case "minwidth":
				failed = int64(config.Width) < limit
			case "maxheight", "maxhight":
				failed = int64(config.Height) > limit
			default:
				failed = int64(config.Height) < limit
			}
		default:
			return nil, &InvalidRuleError{file.Name(), rule}
		}
		if failed {
			message, err := format(s.Language[name], limit)
			if err != nil {
				return nil, err
			}
			messages = append(messages, message)
		}
	}
	if len(messages) == 0 {
		return nil, nil
	}
	return &ValidationErrorMessage{Field: file.Name(), Message: messages}, ErrValidation
}

func decodeConfig(file *os.File) (image.Config, error) {
	position, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return image.Config{}, err
	}
	defer file.Seek(position, io.SeekStart)
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return image.Config{}, err
	}
	config, _, err := image.DecodeConfig(file)
	return config, err
}
