package validation

func Check(data interface{}) ([]*ValidationErrors, error) {
	a := New()
	return a.Validate(data)
}
