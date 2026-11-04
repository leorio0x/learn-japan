# Halaman Kotoba Mengambil Seluruh Data Walau Filter Lesson Aktif

- **Tanggal**: 2026-11-04
- **Komponen**: Kotoba (backend + frontend), reverse proxy (`frontend/nginx.conf`)

## Ringkasan

Saat membuka halaman Kotoba dengan filter Lesson 1 (38 entri), browser tetap
mengunduh seluruh dataset Kotoba (2.910 entri, ±1,68 MB). Di koneksi lambat ini
membuat halaman terasa berat dan lama tampil.

## Penyebab

Filter lesson hanya dilakukan di sisi browser setelah seluruh data terunduh:

- `frontend/src/pages/Kotoba.svelte` memanggil `api.list()` tanpa parameter.
- `frontend/src/lib/api.js` memanggil `GET /api/kotoba` tanpa query filter.
- `backend/internal/kotoba/handler.go` memanggil `svc.List()` tanpa filter.
- `backend/internal/kotoba/service.go` mengembalikan `repo.FindAll()` (semua data).

Selain itu response API tidak dikompresi: `frontend/nginx.conf` belum
mengaktifkan `gzip`, dan backend belum memakai ETag sehingga browser selalu
mengunduh ulang response penuh.

Bukti ukuran (data saat ini):

| Cakupan | Raw | Gzip |
| --- | ---: | ---: |
| Semua lesson | ±1.677.779 B | ±225.043 B |
| Lesson 1 | ±15.267 B | ±2.527 B |

## Solusi

1. **Filter lesson di backend**
   - `internal/kotoba/service.go`: tambah `Filter{Lesson}` dan saring data pada
     `List`.
   - `internal/kotoba/handler.go`: baca `?lesson=` dan teruskan ke service.
   - `GET /api/kotoba?lesson=1` kini hanya mengembalikan data Lesson 1.

2. **Endpoint daftar lesson**
   - Tambah `GET /api/kotoba/lessons` (`Service.Lessons`) yang mengembalikan
     lesson unik terurut numerik.
   - Route didaftarkan sebelum `/:id` agar tidak dianggap sebagai id.

3. **Frontend memuat per lesson**
   - `api.list(params)` dan `api.listKotobaLessons()`.
   - `Kotoba.svelte` mengambil daftar lesson sekali, lalu memuat data sesuai
     lesson terpilih. Hasil disimpan pada cache `Map` per lesson sehingga
     berpindah lesson bolak-balik tidak mengunduh ulang.
   - Request token mencegah response lama menimpa pilihan lesson terbaru.

4. **Kompresi dan cache**
   - `frontend/nginx.conf`: aktifkan `gzip` untuk JSON/JS/CSS/teks dengan
     `gzip_proxied any` agar response dari backend ikut dikompresi.
   - `backend/main.go`: pasang middleware `etag`.
   - Handler Kotoba mengirim `Cache-Control: no-cache` sehingga browser
     merevalidasi lewat ETag dan mendapat `304 Not Modified` bila data tidak
     berubah.

## Verifikasi

- `go test ./...` di `backend/` lulus.
- `bun run build` di `frontend/` lulus.
- Ukuran response Lesson 1 turun dari ±1,68 MB menjadi ±15 KB (±2,5 KB dengan
  gzip).
- Request berikutnya dengan `If-None-Match` yang sama menghasilkan `304`.

## Berkas

- `backend/internal/kotoba/service.go`
- `backend/internal/kotoba/handler.go`
- `backend/internal/kotoba/routes.go`
- `backend/main.go`
- `frontend/src/lib/api.js`
- `frontend/src/pages/Kotoba.svelte`
- `frontend/nginx.conf`
