# Simple GoLang Validation

## Example

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/danangkonang/validation"
)

func main() {
	type T struct {
		Username        string `json:"username" validate:"required,email"`
		Password        string `json:"password" validate:"min=20"`
		ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
	}
	data := T{
		Username:        "",
		Password:        "foo",
		ConfirmPassword: "bar",
	}
	a := validation.New()
	customMessage := map[string]string{
		"required": "your custom message",
		"min":      "minimum {{.}} char",
	}
  a.SetLanguage(customMessage)
	ValidationErrors, err := a.Validate(data)
	if err != nil {
		userJson, _ := json.Marshal(ValidationErrors)
		fmt.Println(string(userJson))
		// [
		//   {
		//     "key": "username",
		//     "message": [
		//       "your message",
		//       "This field must be a valid email address"
		//     ]
		//   },
		//   {
		//     "key": "password",
		//     "message": [
		//       "minimum 20 char"
		//     ]
		//   },
		//   {
		//     "key": "confirm_password",
		//     "message": [
		//       "This field must be a equal with password"
		//     ]
		//   }
		// ]
	}
}

```

## Rules
- required
- alpha
- alphanum
- number
- numeric
- email
- latitude
- longitude
- ip
- boolean
- ipv4
- ipv6
- url
- date
- min
- max
- eqfield