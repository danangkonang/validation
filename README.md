## validation

```go
package main

import (
  "github.com/danangkonang/validation"
  "log"
)

func main() {
  err := validation.Validation("mystring", "required|minlength:5|maxlength:16")
  if err != nil {
		log.Fatal("mystring "+err.Error())
	}

  // now do something
}
```