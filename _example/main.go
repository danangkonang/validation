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
		"required": "your message",
		"min":      "minimum {{.}} char",
	}
	a.SetLanguage(customMessage)
	ValidationErrors, err := a.Validate(data)
	if err != nil {
		userJson, _ := json.Marshal(ValidationErrors)
		fmt.Println(string(userJson))
	}
}
