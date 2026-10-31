# Bulk input CSV: terjemahan berkoma tidak dibungkus kutip

- **Tanggal**: 2026-10-14
- **Komponen**: Bulk input Kotoba & Kanji (`backend/internal/bulk`)

## Ringkasan

Baris CSV seperti
`kanji,そろそろ 失礼します。,Maaf, saya mau pamit dulu.,Minna no Nihongo 1,8`
menghasilkan data terjemahan yang salah. Nilai terjemahan
`Maaf, saya mau pamit dulu.` mengandung koma sehingga dipecah menjadi dua kolom
saat parsing, lalu menggeser seluruh kolom setelahnya:

- terjemahan = `Maaf`
- buku = `saya mau pamit dulu.`
- lesson = `Minna no Nihongo 1`
- `8` menjadi kolom ekstra yang dibuang

## Penyebab

Kolom CSV dipisahkan koma. Field yang mengandung koma **harus** dibungkus tanda
kutip ganda sesuai RFC 4180. Parser (`encoding/csv`, `bulk.go:87`) sudah
mendukung perilaku ini, tetapi panduan di `data/rules.md` maupun teks modal bulk
input belum menegaskan aturan quoting untuk field berkoma, sehingga pengguna
memasukkan nilai tanpa kutip. Karena `reader.FieldsPerRecord = -1`, baris
berkolom berlebih tetap diterima dan data bergeser tanpa error (silent
corruption).

## Solusi

Tidak ada perubahan pada logika parser. Perbaikan difokuskan pada kejelasan
kontrak input:

- `data/rules.md` — tambah aturan: field CSV yang mengandung koma wajib
  dibungkus tanda kutip ganda, termasuk multi arti `|` yang ditulis di dalam
  field terkutip; tambah contoh kasus koma.
- `frontend/src/pages/Crud.svelte` & `CrudKanji.svelte` — tampilkan contoh
  quoting pada keterangan modal bulk input.
- `backend/internal/bulk/bulk_test.go` — selaraskan `TestParseCSVExample`
  dengan contoh `rules.md` (terkutip) dan tambah `TestParseCSVQuotedComma`.

Format input yang benar:
`kanji,そろそろ 失礼します。,"Maaf, saya mau pamit dulu.",Minna no Nihongo 1,8`
dan multi arti: `...,"saya|aku",...`.

## Verifikasi

- `go test ./internal/bulk/...` lulus (termasuk `TestParseCSVQuotedComma`).
- Field terjemahan berkoma dan multi arti terkutip ter-parsing tanpa pergeseran
  kolom.

## Berkas

- `data/rules.md`
- `frontend/src/pages/Crud.svelte`
- `frontend/src/pages/CrudKanji.svelte`
- `backend/internal/bulk/bulk_test.go`
