package validation

import (
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"
	"unicode/utf8"
)

func isRequired(value reflect.Value) bool {
	if !value.IsValid() {
		return false
	}
	if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		return !value.IsNil()
	}
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Map || value.Kind() == reflect.Array || value.Kind() == reflect.String {
		return value.Len() > 0
	}
	return !value.IsZero()
}

func stringValue(value reflect.Value) (string, bool) {
	for value.IsValid() && (value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return "", false
		}
		value = value.Elem()
	}
	if !value.IsValid() || value.Kind() != reflect.String {
		return "", false
	}
	return value.String(), true
}

func runBuiltin(name string, value reflect.Value) bool {
	switch name {
	case "alpha":
		return isAlpha(value)
	case "alphanum":
		return isAlphanum(value)
	case "number":
		return isNumber(value)
	case "numeric":
		return isNumeric(value)
	case "email":
		return isEmail(value)
	case "latitude":
		return isLatitude(value)
	case "longitude":
		return isLongitude(value)
	case "ip":
		return isIP(value)
	case "boolean":
		return isBoolean(value)
	case "ipv4":
		return isIPV4(value)
	case "ipv6":
		return isIPV6(value)
	case "url":
		return isURL(value)
	case "date":
		return isDate(value)
	case "timezone":
		return isTimezone(value)
	default:
		return false
	}
}

func isAlpha(value reflect.Value) bool {
	s, ok := stringValue(value)
	return ok && alphaRegex.MatchString(s)
}
func isAlphanum(value reflect.Value) bool {
	s, ok := stringValue(value)
	return ok && alphaNumericRegex.MatchString(s)
}
func isEmail(value reflect.Value) bool {
	s, ok := stringValue(value)
	return ok && emailRegex.MatchString(s)
}

func isBoolean(value reflect.Value) bool {
	if value.IsValid() && value.Kind() == reflect.Bool {
		return true
	}
	s, ok := stringValue(value)
	return ok && (s == "0" || s == "1" || s == "true" || s == "false" || s == "True" || s == "False")
}

func isIP(value reflect.Value) bool { s, ok := stringValue(value); return ok && net.ParseIP(s) != nil }
func isIPV4(value reflect.Value) bool {
	s, ok := stringValue(value)
	if !ok {
		return false
	}
	if ip := net.ParseIP(s); ip != nil {
		return ip.To4() != nil
	}
	ip, _, err := net.ParseCIDR(s)
	return err == nil && ip.To4() != nil
}
func isIPV6(value reflect.Value) bool {
	s, ok := stringValue(value)
	if !ok {
		return false
	}
	ip, _, err := net.ParseCIDR(s)
	if err == nil {
		return ip.To4() == nil
	}
	parsed := net.ParseIP(s)
	return parsed != nil && parsed.To4() == nil
}

func isURL(value reflect.Value) bool {
	s, ok := stringValue(value)
	if !ok || s == "" {
		return false
	}
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return u.Host != "" && (u.Scheme == "http" || u.Scheme == "https")
}

func isDate(value reflect.Value) bool {
	s, ok := stringValue(value)
	if !ok {
		return false
	}
	_, err := time.Parse("2006-01-02", strings.ReplaceAll(s, "/", "-"))
	return err == nil
}

func numericString(value reflect.Value) (string, bool) {
	for value.IsValid() && (value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return "", false
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return "", false
	}
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(value.Uint(), 10), true
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, value.Type().Bits()), true
	default:
		return stringValue(value)
	}
}

func isNumber(value reflect.Value) bool {
	s, ok := numericString(value)
	return ok && numberRegex.MatchString(s)
}
func isNumeric(value reflect.Value) bool {
	s, ok := numericString(value)
	return ok && numericRegex.MatchString(s)
}
func isLatitude(value reflect.Value) bool {
	s, ok := numericString(value)
	return ok && latitudeRegex.MatchString(s)
}
func isLongitude(value reflect.Value) bool {
	s, ok := numericString(value)
	return ok && longitudeRegex.MatchString(s)
}

func isMinimum(value reflect.Value, rule int) bool {
	value, ok := dereference(value)
	if !ok {
		return false
	}
	switch value.Kind() {
	case reflect.String:
		return utf8.RuneCountInString(value.String()) >= rule
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() >= int64(rule)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() >= uint64(rule)
	case reflect.Float32, reflect.Float64:
		return value.Float() >= float64(rule)
	case reflect.Slice, reflect.Array, reflect.Map:
		return value.Len() >= rule
	default:
		return false
	}
}

func isMaximum(value reflect.Value, rule int) bool {
	value, ok := dereference(value)
	if !ok {
		return false
	}
	switch value.Kind() {
	case reflect.String:
		return utf8.RuneCountInString(value.String()) <= rule
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() <= int64(rule)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() <= uint64(rule)
	case reflect.Float32, reflect.Float64:
		return value.Float() <= float64(rule)
	case reflect.Slice, reflect.Array, reflect.Map:
		return value.Len() <= rule
	default:
		return false
	}
}

func isGTE(value reflect.Value, min float64) bool {
	n, ok := numericFloat(value)
	return ok && n >= min
}
func isLTE(value reflect.Value, max float64) bool {
	n, ok := numericFloat(value)
	return ok && n <= max
}
func numericFloat(value reflect.Value) (float64, bool) {
	value, ok := dereference(value)
	if !ok {
		return 0, false
	}
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(value.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(value.Uint()), true
	case reflect.Float32, reflect.Float64:
		return value.Float(), true
	default:
		return 0, false
	}
}
func isTimezone(value reflect.Value) bool {
	s, ok := stringValue(value)
	if !ok {
		return false
	}
	_, err := time.LoadLocation(s)
	return err == nil
}
func isLength(value reflect.Value, length int) bool {
	value, ok := dereference(value)
	if !ok {
		return false
	}
	if value.Kind() == reflect.String {
		return utf8.RuneCountInString(value.String()) == length
	}
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array || value.Kind() == reflect.Map {
		return value.Len() == length
	}
	return false
}

func dereference(value reflect.Value) (reflect.Value, bool) {
	for value.IsValid() && (value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return reflect.Value{}, false
		}
		value = value.Elem()
	}
	return value, value.IsValid()
}
func isEmptyValue(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}
	if value.Kind() == reflect.String {
		return strings.TrimSpace(value.String()) == ""
	}
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		return value.IsNil()
	}
	if value.Kind() == reflect.Struct && value.Type() == reflect.TypeOf(time.Time{}) {
		return value.Interface().(time.Time).IsZero()
	}
	return value.IsZero()
}
