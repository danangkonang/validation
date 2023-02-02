package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/danangkonang/validation"
)

func main() {
	// type T struct {
	// 	Avatar *os.File `json:"avatar" validate:"maxsize=10000,minsize=10000,maxwidth=200,minhight=200"`
	// }
	r, err := os.Open("test.jpg")
	if err != nil {
		panic(err)
	}
	a := validation.New()
	ValidationErrors, err := a.FileValidate(r, "maxsize=10000,minsize=10000,maxwidth=200,minhight=200")
	if err != nil {
		userJson, _ := json.Marshal(ValidationErrors)
		fmt.Println(string(userJson))
	}
}
