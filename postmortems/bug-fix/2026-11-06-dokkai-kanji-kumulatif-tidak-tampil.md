# Dokkai Kanji Kumulatif 1-15 Tidak Tampil di Frontend

- **Tanggal**: 2026-11-06
- **Komponen**: Dokkai Kanji (frontend `DokkaiView`, backend `internal/dokkai`)

## Ringkasan

File `data/dokkai/kanji/1-15.json` sudah berhasil dibuat, tetapi opsi
"Lesson 1-15" tidak muncul pada halaman Dokkai Kanji sehingga isinya tidak bisa
dibuka dari frontend.

## Penyebab

Frontend membangun daftar pilihan lesson dari data kanji
(`api.listKanji()`), bukan dari file dokkai yang benar-benar tersedia. Data
kanji hanya memiliki lesson `1` sampai `15`, sehingga nilai kumulatif `1-15`
tidak pernah masuk ke dropdown.

Selain itu, backend belum menyediakan endpoint untuk mendaftar file dokkai yang
ada pada suatu tipe (`kotoba`/`kanji`). Akibatnya frontend tidak punya cara
mengetahui keberadaan `1-15.json`.

## Solusi

1. Menambahkan `Repository.Lessons(kind)` yang membaca direktori
   `data/dokkai/{kind}` dan mengembalikan nama file `.json` yang valid.
2. Menambahkan endpoint `GET /api/dokkai/:type/lessons`, didaftarkan sebelum
   route generik `/:type/:lesson`.
3. Mengubah `DokkaiView` agar mengambil daftar lesson dari endpoint baru
   tersebut, sehingga file kumulatif seperti `1-15.json` otomatis tampil.

Urutan sorting dibuat numerik untuk lesson biasa, lalu lesson rentang
(misalnya `1-15`) ditampilkan setelahnya.

## Verifikasi

- `go test ./internal/dokkai/...` berhasil.
- Test baru memastikan `1-15.json` terdaftar dan endpoint
  `/dokkai/kanji/lessons` mengembalikan `["1", "1-15"]`.
- `npm run build` pada frontend berhasil.
- Perlu restart backend dan rebuild frontend agar perubahan aktif.

## Berkas

- `backend/internal/dokkai/repository.go`
- `backend/internal/dokkai/service.go`
- `backend/internal/dokkai/handler.go`
- `backend/internal/dokkai/routes.go`
- `backend/internal/dokkai/lessons_test.go`
- `frontend/src/lib/api.js`
- `frontend/src/components/DokkaiView.svelte`
- `data/dokkai/kanji/1-15.json`
