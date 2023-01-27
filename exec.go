package validation

import (
	"fmt"
	"reflect"
	"strconv"
)

func isRequired(field reflect.Value) bool {
	return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
}

func isAlpha(fl reflect.Value) bool {
	return alphaRegex.MatchString(fl.String())
}

func isAlphanum(fl reflect.Value) bool {
	return alphaNumericRegex.MatchString(fl.String())
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
