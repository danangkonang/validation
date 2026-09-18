# Go Validation

Library validasi berbasis reflection untuk field exported pada `struct` Go.
Aturan dibaca dari tag `validate`; pilihan nilai string dapat menggunakan tag
`enum`.

## Persyaratan dan instalasi

- Go 1.16 atau lebih baru

```bash
go get github.com/danangkonang/validation
```

## Penggunaan dasar

```go
package main

import (
    "errors"
    "fmt"

    "github.com/danangkonang/validation"
)

type User struct {
    Email string `json:"email" validate:"required,email"`
    Role  string `json:"role" enum:"admin,user" validate:"required"`
    Age   int    `json:"age" validate:"gte=18,lte=120"`
}

func main() {
    validationErrors, err := validation.New().Validate(User{
        Email: "bad", Role: "guest", Age: 17,
    })
    if errors.Is(err, validation.ErrValidation) {
        for _, fieldError := range validationErrors {
            fmt.Printf("%s: %v\n", fieldError.Field, fieldError.Message)
        }
    }
}
```

`Validate` menerima `struct` atau pointer ke `struct`. Jika data valid, hasil
dan error sama-sama `nil`. Untuk penggunaan singkat tersedia `Check`, yang
membuat validator baru:

```go
validationErrors, err := validation.Check(User{Email: "bad"})
```

## Format tag dan built-in rules

Rule ditulis sebagai daftar yang dipisahkan koma. Argument ditulis dengan `=`:

```go
type Account struct {
    Username string `validate:"required,alphanum,min=3,max=30"`
    Age      int    `validate:"gte=18,lte=120"`
}
```

| Rule | Keterangan |
| --- | --- |
| `required` | Nilai harus terisi; string, collection, dan array tidak boleh kosong. Pointer tidak boleh `nil`. |
| `omitempty` | Jika nilai kosong, rule lain pada field dilewati. |
| `alpha` | String alfabet. |
| `alphanum` | String alfanumerik. |
| `number` | Angka integer. |
| `numeric` | Angka, termasuk desimal. |
| `email` | Alamat email valid. |
| `latitude` / `longitude` | Koordinat valid. |
| `ip` | Alamat IP valid; CIDR tidak diterima. |
| `ipv4` / `ipv6` | Alamat IP yang sesuai; CIDR juga diterima. |
| `boolean` | Boolean Go atau string `0`, `1`, `true`, `false`, `True`, `False`. |
| `url` | URL absolut dengan skema `http` atau `https`. |
| `date` | Tanggal `YYYY-MM-DD` atau `YYYY/MM/DD`. |
| `timezone` | Nama timezone yang valid, misalnya `Asia/Jakarta`. |
| `min=N` / `max=N` | Batas panjang string, collection, atau nilai numerik. |
| `len=N` | Panjang string, slice, array, atau map harus tepat `N`. |
| `gte=N` / `lte=N` | Batas nilai numerik; mendukung integer, unsigned integer, dan float. |
| `eqfield=Field` | Nilai harus sama dengan field lain pada struct yang sama. |

`min`, `max`, dan `len` menghitung string berdasarkan Unicode rune, bukan byte.
Spasi di sekitar rule dan argument diabaikan. Rule yang tidak dikenal, argument
yang salah, atau argument yang hilang menghasilkan `*InvalidRuleError`.

Contoh kombinasi rule:

```go
type Registration struct {
    Password        string `validate:"required,min=8"`
    ConfirmPassword string `validate:"required,eqfield=Password"`
    Website         string `validate:"omitempty,url"`
}
```

## Enum

`enum` hanya berlaku untuk field string. Nilai dipisahkan koma dan bersifat
case-sensitive. Validasi enum dilewati jika field kosong.

```go
type Request struct {
    Status string `validate:"required" enum:"pending,approved,rejected"`
}
```

Jika `enum` dipasang pada tipe selain string, validator mengembalikan
`*InvalidRuleError`.

## Nested struct dan collection

Validator menelusuri nested struct, pointer, slice, array, dan map. Nama field
menggunakan tag `json` jika tersedia. Contoh path error:

```go
type Child struct {
    Name string `json:"name" validate:"required"`
}

type Parent struct {
    Children []Child `json:"children"`
}
```

Error pada child pertama memiliki field `children[0].name`. Hanya field
exported yang divalidasi; field dengan `json:"-"` dilewati.

## Hasil dan error

Setiap error field memiliki struktur berikut:

```go
type ValidationErrorMessage struct {
    Index   string        `json:"index,omitempty"`
    Field   string        `json:"key,omitempty"`
    Message []interface{} `json:"message,omitempty"`
}
```

Jenis error utama:

- `ErrValidation`: data melanggar satu atau lebih rule.
- `ErrInvalidInput`: input `nil` atau bukan `struct`/pointer ke `struct`.
- `*InvalidRuleError`: rule atau konfigurasi rule tidak valid.

```go
validationErrors, err := validation.New().Validate(data)
switch {
case errors.Is(err, validation.ErrValidation):
    // Gunakan validationErrors untuk menampilkan error per field.
case errors.Is(err, validation.ErrInvalidInput):
    // Input tidak dapat divalidasi.
case err != nil:
    var ruleError *validation.InvalidRuleError
    if errors.As(err, &ruleError) {
        fmt.Println(ruleError)
    }
}
```

## Custom rule

Custom rule didaftarkan pada instance validator dan menerima `reflect.Value`:

```go
v := validation.New()
_ = v.RegisterValidation("startsWithA", func(value reflect.Value) bool {
    return value.Kind() == reflect.String && strings.HasPrefix(value.String(), "A")
})

type User struct {
    Name string `validate:"required,startsWithA"`
}

_, err := v.Validate(User{Name: "Budi"})
```

Nama rule tidak boleh kosong dan callback tidak boleh `nil`. Callback yang panic
dianggap menghasilkan validasi gagal.

## Bahasa dan pesan error

Setiap validator memiliki salinan language map sendiri. `SetLanguage` tidak
mengubah default package atau validator lain. Placeholder `{{.}}` berisi
argument rule.

```go
v := validation.New()
v.SetLanguage(map[string]string{
    "required": "Field ini wajib diisi",
    "min":      "Minimal {{.}} karakter",
})
```

Pesan custom rule dapat ditambahkan menggunakan nama rule custom tersebut.

## Validasi file dan gambar

`FileValidate` menerima `*os.File` dan mendukung:

- `minsize=N`, `maxsize=N`: ukuran file dalam byte.
- `minwidth=N`, `maxwidth=N`: lebar gambar dalam pixel.
- `minheight=N`, `maxheight=N`: tinggi gambar dalam pixel.
- `minhight=N`, `maxhight=N`: ejaan lama yang tetap didukung.

Rule width dan height memerlukan file gambar yang dapat dibaca decoder image
Go. File harus dibuka terlebih dahulu:

```go
file, err := os.Open("avatar.jpg")
if err != nil {
    panic(err)
}
defer file.Close()

validationError, err := validation.New().FileValidate(
    file, "maxsize=1000000,minwidth=200,minheight=200",
)
if errors.Is(err, validation.ErrValidation) {
    fmt.Println(validationError.Message)
}
```

Jika semua rule lolos, hasil dan error sama-sama `nil`. File `nil`, rule
malformed, atau file non-gambar untuk rule dimensi menghasilkan error.

## Menjalankan contoh

```bash
go run ./_example
go run ./_example/file
```

Contoh validasi file memerlukan `test.jpg` pada working directory saat program
dijalankan.

## Development

```bash
GOCACHE=/tmp/validation-go-cache go test ./...
GOCACHE=/tmp/validation-go-cache go test -race ./...
GOCACHE=/tmp/validation-go-cache go vet ./...
```

Sebelum membuat perubahan, jalankan `gofmt` pada file Go dan tambahkan regression
test untuk perubahan perilaku.
