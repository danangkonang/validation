package main

import (
	"encoding/json"
	"fmt"

	"github.com/danangkonang/validation"
)

func main() {
	type T struct {
		// Username string `json:"username" validate:"required,email"`
		// Password string `json:"password" validate:"min=20"`
		// ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
		// Hobbies []string `json:"hobbies" validate:"required"`
		Age int `json:"age" validate:"min=1"`
	}

	jsonData := []byte(`{"username": "Alice", "age": null, "password": "foo", "confirm_password": "bar"}`)

	// data := T{
	// 	// Username: "",
	// 	// Password:        "foo",
	// 	// ConfirmPassword: "bar",
	// 	Age: 0,
	// }
	var data T
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	a := validation.New()
	customMessage := map[string]string{
		// "required": "your message",
		// "min": "minimum {{.}} char",
	}
	a.SetLanguage(customMessage)
	ValidationErrors, err := a.Validate(data)
	if err != nil {
		userJson, _ := json.Marshal(ValidationErrors)
		fmt.Println(string(userJson))
	}
}
