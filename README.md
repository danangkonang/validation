## validation

```go
package main

import (
  "fmt"

  "github.com/danangkonang/validation"
)

func main() {
  err := validation.Validation("mystring", "required|minlength:5|maxlength:16")
  fmt.Println(err)
  // now do something
}
```

## rules
- required
- minlength
- maxlength
- email