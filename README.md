# Simple GoLang Validation

## Install

```bash
go get github.com/danangkonang/validation
```

## Example

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/danangkonang/validation"
)

func main() {
	type Children struct {
		Level string `json:"level" enum:"beginner,intermediate,advanced"`
	}
	type T struct {
		Username        string     `json:"username" validate:"required,email"`
		Password        string     `json:"password" validate:"min=20"`
		ConfirmPassword string     `json:"confirm_password" validate:"eqfield=Password"`
		Children        []Children `json:"children" validate:"required"`
	}
	data := T{
		Username:        "",
		Password:        "foo",
		ConfirmPassword: "bar",
		Children: []Children{
			{
				Level: "",
			},
		},
	}
	v := validation.New()
	customMessage := map[string]string{
		"required": "your custom message",
		"min":      "minimum {{.}} char",
	}
	v.SetLanguage(customMessage)
	ValidationErrorMessage, err := v.Validate(data)
	if err != nil {
		userJson, _ := json.Marshal(ValidationErrorMessage)
		fmt.Println(string(userJson))
		// [
		// 	{
		// 		"key": "username",
		// 		"message": [
		// 			"This field is required",
		// 			"This field must be a valid email address"
		// 		]
		// 	},
		// 	{
		// 		"key": "password",
		// 		"message": [
		// 			"This field must be a minimum 20"
		// 		]
		// 	},
		// 	{
		// 		"key": "confirm_password",
		// 		"message": [
		// 			"This field must be a equal with password"
		// 		]
		// 	},
		// 	{
		// 		"key": "children",
		// 		"message": [
		// 			{
		// 				"key": "level",
		// 				"message": [
		// 					"This field must be one of [beginner,intermediate,advanced]"
		// 				]
		// 			}
		// 		]
		// 	}
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
- enum