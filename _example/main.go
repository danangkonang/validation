package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/danangkonang/validation"
)

func main() {
	type T struct {
		// Username string `json:"username" validate:"required,email"`
		// Password string `json:"password" validate:"min=20"`
		// ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
		// Hobbies []string `json:"hobbies" validate:"required"`
		Age    int     `json:"age" validate:"min=1"`
		Number []int   `json:"number" validate:"required"`
		Cuid2  *string `json:"cuid" validate:"omitempty,len=24"`
	}

	jsonData := []byte(`{"username": "Alice", "age": null, "password": "foo", "confirm_password": "bar", "number": [], "cuid": "1"}`)

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
	a.RegisterValidation("custom", func(v reflect.Value) bool { return true })
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
