package validation

func Check(data interface{}) ([]*ValidationErrorMessage, error) {
	a := New()
	return a.Validate(data)
}
