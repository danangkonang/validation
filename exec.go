package validation

import (
	"fmt"
	"reflect"
	"strconv"
)

func isRequired(field reflect.Value) bool {
	// fmt.Println(field.Kind())
	// fmt.Println(field.IsValid())
	// fmt.Println(field.Interface())
	// fmt.Println(reflect.Zero(field.Type()).Interface())
	// return field.IsValid() && (field.Interface() != reflect.Zero(field.Type()).Interface())
	// return true
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		// if fl.(*validate).fldIsPointer && field.Interface() != nil {
		// 	return true
		// }
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

func isAlpha(fl reflect.Value) bool {
	return alphaRegex.MatchString(fl.String())
}

func isAlphanum(fl reflect.Value) bool {
	return alphaNumericRegex.MatchString(fl.String())
}

func isBoolean(fl reflect.Value) bool {
	bools := []string{"0", "1", "true", "false", "True", "False"}
	for _, b := range bools {
		if b == fl.String() {
			return true
		}
	}
	return false
}

func isIP(fl reflect.Value) bool {
	return ipRegex.MatchString(fl.String())
}

// Ref: https://en.wikipedia.org/wiki/IPv4
func isIPV4(fl reflect.Value) bool {
	return ipV4Regex.MatchString(fl.String())
}

// Ref: https://en.wikipedia.org/wiki/IPv6
func isIPV6(fl reflect.Value) bool {
	return ipV6Regex.MatchString(fl.String())
}

func isURL(fl reflect.Value) bool {
	return urlRegex.MatchString(fl.String())
}

// isDate check the date string is valid or not
func isDate(fl reflect.Value) bool {
	return dateRegex.MatchString(fl.String())
}

func isEmail(fl reflect.Value) bool {
	return emailRegex.MatchString(fl.String())
}

func isNumber(fl reflect.Value) bool {
	switch fl.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		return true
	default:
		return numberRegex.MatchString(fl.String())
	}
}

func isNumeric(fl reflect.Value) bool {
	switch fl.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		return true
	default:
		return numericRegex.MatchString(fl.String())
	}
}

func isLongitude(fl reflect.Value) bool {
	field := fl
	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 32)
	case reflect.Float64:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 64)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
	return longitudeRegex.MatchString(v)
}

func isLatitude(fl reflect.Value) bool {
	field := fl
	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 32)
	case reflect.Float64:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 64)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
	return latitudeRegex.MatchString(v)
}

func isMinimum(fl reflect.Value, rule int) bool {
	// Handle invalid or nil values
	if !fl.IsValid() || (fl.Kind() == reflect.Ptr && fl.IsNil()) {
		return rule <= 0 // A nil value satisfies min=0, fails otherwise
	}

	switch fl.Kind() {
	case reflect.String:
		return len(fl.String()) >= rule // or utf8.RuneCountInString(fl.String()) for characters
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fl.Int() >= int64(rule)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fl.Uint() >= uint64(rule)
	case reflect.Float32, reflect.Float64:
		return fl.Float() >= float64(rule)
	case reflect.Slice, reflect.Array:
		return fl.Len() >= rule
	case reflect.Ptr:
		// Recursively check the dereferenced value
		return isMinimum(fl.Elem(), rule)
	default:
		// Unsupported type, treat as invalid
		return false
	}
}

func isMaximum(fl reflect.Value, rule int) bool {
	// if len(fl.String()) <= rule {
	// 	return true
	// } else {
	// 	return false
	// }
	if !fl.IsValid() || (fl.Kind() == reflect.Ptr && fl.IsNil()) {
		return rule <= 0 // A nil value satisfies min=0, fails otherwise
	}

	switch fl.Kind() {
	case reflect.String:
		return len(fl.String()) <= rule // or utf8.RuneCountInString(fl.String()) for characters
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fl.Int() <= int64(rule)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fl.Uint() <= uint64(rule)
	case reflect.Float32, reflect.Float64:
		return fl.Float() <= float64(rule)
	case reflect.Slice, reflect.Array:
		return fl.Len() <= rule
	case reflect.Ptr:
		// Recursively check the dereferenced value
		return isMinimum(fl.Elem(), rule)
	default:
		// Unsupported type, treat as invalid
		return false
	}
}

func isEqualField(fl reflect.Value, rule string) bool {
	return fl.String() == rule
}
