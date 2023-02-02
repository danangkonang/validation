package validation

var (
	Lang = map[string]string{
		"required":  "This field is required",
		"alpha":     "This field must be a valid alpha",
		"alphanum":  "This field must be a valid alphanum",
		"number":    "This field must be a valid number",
		"numeric":   "This field must be a valid numeric",
		"email":     "This field must be a valid email address",
		"latitude":  "This field must be a valid latitude",
		"longitude": "This field must be a valid longitude",
		"ip":        "This field must be a valid ip",
		"boolean":   "This field must be a valid boolean",
		"ipv4":      "This field must be a valid ipv4",
		"ipv6":      "This field must be a valid ipv6",
		"url":       "This field must be a valid url",
		"date":      "This field must be a valid date",
		"min":       "This field must be a minimum {{.}}",
		"max":       "This field must be a maximum {{.}} charakter",
		"eqfield":   "This field must be a equal with {{.}}",
		"maxsize":   "This field must be a maximum size {{.}}",
		"minsize":   "This field must be a minimum size {{.}}",
		"maxwidth":  "This field must be a maximum width {{.}}",
		"minwidth":  "This field must be a minimum width {{.}}",
		"maxhight":  "This field must be a maximum hight {{.}}",
		"minhight":  "This field must be a minimum hight {{.}}",
	}
)
