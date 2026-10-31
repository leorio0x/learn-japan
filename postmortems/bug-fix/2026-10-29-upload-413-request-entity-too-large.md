# Upload Bulk Input: 413 Request Entity Too Large (hanya di prod)

- **Tanggal**: 2026-10-29
- **Komponen**: Bulk input Kotoba/Kanji, reverse proxy (`frontend/nginx.conf`)

## Ringkasan

Saat mengunggah file bulk input (csv/xlsx/json) di lingkungan produksi, server
membalas:

```
413 Request Entity Too Large
```

Masalah hanya muncul di prod; di lingkungan dev upload berjalan normal.

## Penyebab

Di prod, request dari browser melewati container frontend (nginx) sebelum
diteruskan ke backend:

```
browser -> nginx (frontend:80) --/api/--> backend:3000
```

`frontend/nginx.conf` tidak pernah mengatur `client_max_body_size`, sehingga
nginx memakai nilai default **1 MB**. File yang lebih besar dari 1 MB ditolak
nginx dengan 413 dan tidak pernah sampai ke backend.

Di dev tidak ada nginx: Vite mem-proxy `/api` langsung ke `http://localhost:3000`
(`frontend/vite.config.js`), sehingga yang berlaku hanya `BodyLimit` Fiber
sebesar 10 MB (`backend/main.go:41`).

Bukti:

- `frontend/nginx.conf` — tidak ada `client_max_body_size` (default 1 MB).
- `frontend/src/lib/api.js:34` — upload multipart dikirim ke `/api/...`.
- `backend/main.go:41` — `BodyLimit: 10 * 1024 * 1024`.
- `backend/internal/bulk/bulk.go:20` — batas ukuran file 5 MB.

## Solusi

Tambahkan limit body di level `server` pada `frontend/nginx.conf`, disamakan
dengan backend (10 MB):

```nginx
server {
    listen 80;
    server_name _;

    client_max_body_size 10m;
    ...
}
```

## Verifikasi

- File bulk input berukuran antara 1 MB dan 10 MB berhasil diunggah di prod.
- File di atas 10 MB tetap ditolak oleh backend (bukan nginx) dengan pesan JSON
  `{"error": ...}` karena `BodyLimit` Fiber.
- Rebuild/deploy container frontend diperlukan agar konfigurasi baru aktif:
  `docker compose build frontend && docker compose up -d frontend`.

## Berkas

- `frontend/nginx.conf`
