# TODO Proyek Validation

Dokumen ini disusun dari audit kode pada 2026-09-06. Proyek saat ini adalah
library validasi struct Go berbasis reflection, dengan dukungan tag `validate`,
tag `enum`, custom validator, pesan bahasa, dan validasi file/gambar.

## Selesai pada pass ini

- [x] Hardening reflection, nested path, pointer handling, `eqfield`, `enum`,
  `omitempty`, malformed rule, dan error type dasar.
- [x] Isolasi language map per instance dan validasi callback nil/tag kosong.
- [x] Validasi IP berbasis `net`, date berbasis `time.Parse`, dan URL berbasis
  `net/url`.
- [x] Hardening `FileValidate`, termasuk nil/error handling, validasi image
  config, posisi file, dan alias height legacy.
- [x] Menambahkan regression tests, fuzz test, CI untuk test/race/vet/coverage,
  serta memperbarui README dan CONTRIBUTION.

## Status saat ini

- `go test ./...` lulus.
- `go test -race ./...` lulus.
- `go vet ./...` lulus.
- Cakupan test masih rendah: hanya beberapa skenario `required`, nested slice,
  dan email.
- `go.mod` menyatakan Go 1.16, tetapi kode memakai `reflect.Pointer`; kontrak
  versi Go perlu diputuskan dan diuji secara eksplisit.

## P0 - Stabilitas dan correctness

- [x] Tambahkan test regresi untuk input pointer, pointer nil, pointer field,
  struct nested, slice struct, slice scalar, map, array, dan interface.
  Validasi nested jangan bergantung pada adanya rule `required`; nested value
  yang memiliki tag sendiri harus tetap divalidasi.
- [x] Perbaiki `required` agar konsisten untuk semua tipe yang didukung,
  termasuk pointer, interface, map, bool, angka, array, dan time.Time. Saat ini
  beberapa tipe tidak diproses dan dapat dianggap valid tanpa pemeriksaan.
- [x] Ganti akses reflection yang dapat memakai nilai invalid dengan error
  konfigurasi yang jelas. Kasus utama: `eqfield` menunjuk field yang tidak ada,
  field pembanding pointer/non-string, atau field yang tidak bisa dibandingkan.
- [x] Pastikan validator tidak panic karena tipe field yang tidak sesuai. Rule
  seperti `latitude`/`longitude` saat ini melakukan panic pada tipe unsupported;
  kembalikan validation error atau error konfigurasi yang terdokumentasi.
- [x] Perbaiki parser rule/tag agar argumen hilang, argumen berlebih, angka
  invalid, dan rule tidak dikenal tidak diam-diam diabaikan. Pilih satu kontrak:
  menolak konfigurasi dengan error, atau menandainya sebagai validation error.
- [x] Hardening `FileValidate`: validasi file nil, cek semua error `Stat`, parse
  angka, dan `image.DecodeConfig`; jangan mengakses `rv[1]` sebelum memastikan
  argumen tersedia. Reset/duplikasi reader sebelum decode berulang agar hasil
  rule dimensi tidak dipengaruhi posisi file.

## P1 - Kontrak API dan kualitas hasil validasi

- [x] Tetapkan spesifikasi rule: tipe yang diterima, arti `min`/`max`/`len`,
  apakah panjang dihitung dalam byte atau rune, perilaku nilai kosong, dan
  urutan evaluasi rule. Tulis tabel ini di README dan jadikan test sebagai
  executable specification.
- [x] Buat error type yang dapat diinspeksi, misalnya sentinel/type untuk
  invalid input, invalid rule, dan validation failure. Hindari hanya
  mengandalkan string `"form error"`.
- [x] Ubah `ValidationErrorMessage` agar struktur nested memiliki index/path
  yang stabil. Field seperti `children[0].level` lebih mudah dipakai client
  daripada object message bersarang tanpa index.
- [x] Perbaiki `eqfield` agar membandingkan nilai berdasarkan tipe sebenarnya,
  termasuk angka, bool, pointer, dan string; jangan menggunakan `.String()`
  sebagai universal equality.
- [x] Evaluasi ulang perilaku `omitempty`: nilai kosong harus didefinisikan
  untuk setiap tipe dan skip harus berlaku konsisten pada rule `enum` juga bila
  memang itu yang diinginkan.
- [x] Validasi `enum` berdasarkan tipe field atau batasi secara eksplisit ke
  string. Saat ini implementasinya memakai `Value.String()` sehingga tipe
  non-string tidak menghasilkan kontrak yang intuitif.
- [ ] Validasi callback custom: tolak tag kosong dan callback nil, dokumentasikan
  apakah callback boleh mengubah state, dan pertimbangkan API callback yang
  mengembalikan error agar alasan kegagalan dapat disampaikan.

## P1 - Keamanan, concurrency, dan input parsing

- [x] Jangan membagikan map global `Lang` ke setiap instance `Validation`.
  `New()` sebaiknya membuat salinan map agar `SetLanguage` pada satu instance
  tidak mengubah instance lain dan aman dari race saat dipakai server.
- [x] Ganti `template.Must` di `format` dengan parsing/eksekusi yang
  mengembalikan error. Pesan custom yang salah format tidak boleh menjatuhkan
  proses aplikasi.
- [x] Ganti regex IP/IPv4/IPv6 dengan `net.ParseIP`/`net.ParseCIDR` dan tinjau
  regex URL, email, serta date. Regex IPv4 saat ini dapat menerima oktet di atas
  255; regex IPv6 juga memiliki pola yang tampak typo.
- [x] Tambahkan batas dan kebijakan untuk validasi file: ukuran maksimum,
  tipe/format yang diizinkan, penanganan file seekable/non-seekable, dan error
  untuk format gambar tidak valid.

## P2 - Refactoring maintainability

- [ ] Pisahkan parser tag, traversal reflection, rule registry, formatter pesan,
  dan validator file ke modul/fungsi yang lebih kecil. `Validate` saat ini
  memuat traversal, dispatch semua rule, nested validation, serta formatting.
- [ ] Buat registry rule built-in sehingga dispatch tidak terus membesar dalam
  satu `switch`; registry juga memudahkan dokumentasi dan test per rule.
- [x] Hapus komentar debug/kode lama yang tidak aktif dan rapikan penamaan,
  termasuk typo public/internal seperti `hight` dan `charakter`, dengan alias
  kompatibilitas bila tag lama sudah dipakai pengguna.
- [x] Ganti helper `split` dengan `strings.Split` kecuali ada alasan khusus;
  tambahkan trimming whitespace yang konsisten pada tag dan enum.
- [ ] Pertimbangkan cache metadata struct/tag untuk mengurangi biaya reflection
  pada endpoint yang memvalidasi tipe sama berulang kali. Ukur benchmark dulu.
- [x] Putuskan dukungan Go: naikkan `go.mod` ke versi yang mendukung
  `reflect.Pointer`, atau gunakan API yang kompatibel dengan Go 1.16. Tambahkan
  CI pada semua versi Go yang benar-benar dijanjikan.

## P2 - Test, dokumentasi, dan release

- [ ] Ubah test menjadi table-driven untuk seluruh built-in rule, dengan kasus
  valid, invalid, empty, boundary, tipe unsupported, dan malformed parameter.
- [ ] Tambahkan test untuk error/panic boundary pada `FileValidate`, custom
  validator, bahasa custom, field tidak ditemukan, dan nested pointer.
- [x] Tambahkan fuzz test untuk parser rule, `Validate`, regex/input string, dan
  pesan template. Target utamanya adalah memastikan input arbitrary tidak panic.
- [x] Jalankan `go test -race`, `go vet`, dan coverage di CI; tambahkan linting
  yang kompatibel dengan versi Go proyek. Tetapkan target coverage berdasarkan
  behavior penting, bukan angka semata.
- [x] Perbarui README dengan semua rule yang benar-benar didukung (`len`, `gte`,
  `lte`, `timezone`, custom validator, file validation), tipe input, contoh
  pointer/nested error, serta daftar breaking change.
- [ ] Tambahkan LICENSE, CONTRIBUTING yang lebih lengkap, changelog, dan
  workflow release/tag yang konsisten. `CONTRIBUTION.md` saat ini masih berisi
  instruksi tag versi spesifik `v0.0.17`.

## Urutan pengerjaan yang disarankan

1. Kunci perilaku melalui spesifikasi dan test regresi P0.
2. Perbaiki panic dan silent failure pada reflection, parser, serta file
   validation.
3. Pisahkan error type dan rapikan format/path error agar API stabil.
4. Benahi global state bahasa, parser network-like input, dan kontrak versi Go.
5. Refactor internal setelah behavior terlindungi oleh test.
6. Lengkapi CI, fuzz/benchmark, dokumentasi, dan proses release.

## Definition of Done

- Tidak ada input tag atau nilai pengguna yang menyebabkan panic.
- Semua built-in rule memiliki test boundary dan dokumentasi tipe yang didukung.
- Nested error memiliki path/index stabil dan dapat diproses API client.
- Instance validator dapat dipakai bersamaan tanpa berbagi mutable global state.
- CI menjalankan test normal, race test, vet, lint, dan coverage pada versi Go
  yang didukung.
- README, contoh, changelog, dan versi module sesuai dengan behavior aktual.
