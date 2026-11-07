# Cloudflare Tunnel: connector restarting, domain 530 (error 1033)

- **Tanggal**: 2026-11-07
- **Komponen**: Deployment, `docker-compose.yml`, Cloudflare Tunnel, `frontend/nginx.conf`, `backend`

## Ringkasan

Setelah memasang Cloudflare Tunnel untuk `https://ansatsu.hirsuta.id`,
request dari internet membalas:

```
HTTP/1.1 530
error code: 1033
```

Container `learn-japan-tunnel` berstatus `Restarting (255)` dengan log:

```
Provided Tunnel token is not valid.
See 'cloudflared tunnel run --help'.
```

Setelah token diperbaiki, error berubah menjadi:

```
"cloudflared tunnel run" requires the ID or name of the tunnel to run
as the last command line argument or in the configuration file.
```

Sementara itu jalur lama `http://116.212.74.154:8080` tetap `200`, jadi
aplikasi sendiri tidak bermasalah.

## Penyebab

Ada beberapa masalah berurutan:

1. **Token yang disimpan bukan token connector.**
   Nilai di `.env` berformat `cfut_...` dengan panjang 53 karakter. Itu token
   API, bukan tunnel token. `cloudflared` hanya menerima token connector:
   base64 JSON (`{"a":...,"t":...,"s":...}`) yang diawali `eyJ` dan panjangnya
   ±184 karakter. Token `cfut_` gagal di-decode sehingga connector menolak
   start. Token diambil dari halaman **Install connector → Docker**, bukan dari
   **My Profile → API Tokens**.

2. **`.env` di VPS belum ada sejak awal.**
   Compose adalah otomatis membaca `.env` di folder yang sama dengan
   `docker-compose.yml`, tetapi selama ini nilainya memakai default. Saat
   `.env` dibuat, `DATA_DIR` harus disamakan dengan lokasi data yang sudah
   dipakai (`/opt/learn-japan/data`), kalau tidak aplikasi tampak kehilangan
   data.

3. **Perubahan `.env` tidak langsung terbaca container.**
   Variabel `environment:` di-inject saat container dibuat. `docker compose
   restart` tetap memakai env lama; perlu
   `docker compose up -d --force-recreate cloudflared`.

4. **Salah konteks perintah.**
   `cloudflared.exe service install <token>` adalah perintah untuk memasang
   service di Windows, tidak berlaku untuk deployment Docker di VPS.
   Container hanya butuh `TUNNEL_TOKEN` berisi token connector.

Catatan diagnosis yang sempat membingungkan: query
`Resolve-DnsName ansatsu.hirsuta.id -Type CNAME` mengembalikan SOA (seolah
record tidak ada). Itu normal karena Cloudflare melakukan **CNAME flattening**
untuk record proxied; buktinya query `-Type A` mengembalikan IP edge
Cloudflare (`104.21.78.177`, `172.67.136.37`).

## Solusi

1. Tambahkan service `cloudflared` ke `docker-compose.yml`:

   ```yaml
   cloudflared:
     image: cloudflare/cloudflared:latest
     container_name: learn-japan-tunnel
     command: tunnel --no-autoupdate run
     environment:
       TUNNEL_TOKEN: ${CLOUDFLARE_TUNNEL_TOKEN}
     depends_on:
       - frontend
     restart: unless-stopped
   ```

2. Buat `/opt/learn-japan/.env` di VPS (tidak di-commit, `chmod 600`) berisi
   `DATA_DIR`, `CRUD_PASSWORD`, `SESSION_SECRET`, dan
   `CLOUDFLARE_TUNNEL_TOKEN` yang diambil dari **Networking → Tunnels →
   [tunnel] → Overview → Add a replica → Docker** (string setelah `--token`).

3. Terapkan ulang container dengan env baru:

   ```bash
   docker compose up -d --force-recreate cloudflared
   ```

4. Karena token sempat terekspos, token di-rotate di dashboard dan nilai di
   `.env` diganti, lalu container di-recreate lagi.

5. Lengkapi konfigurasi HTTPS:

   - `frontend/nginx.conf`: `server_name ansatsu.hirsuta.id`, real IP dari
     `CF-Connecting-IP`, `X-Forwarded-Proto` diambil dari `CF-Visitor`,
     `Cache-Control: no-store` untuk `/api/`, header keamanan.
   - `backend/main.go`: `ProxyHeader: X-Forwarded-For` dan trusted proxy
     (range privat Docker), CORS ke `https://ansatsu.hirsuta.id`.
   - `backend/internal/auth/auth.go`: cookie `Secure` dinamis berdasarkan
     `X-Forwarded-Proto`.

6. Setelah tunnel stabil, port lama dibatasi ke localhost
   (`127.0.0.1:8080:80`) agar proteksi Cloudflare tidak bisa dilewati.

## Verifikasi

- `docker logs learn-japan-tunnel` menampilkan
  `Registered tunnel connection` pada beberapa `connIndex`.
- `https://ansatsu.hirsuta.id/api/health` → `200 {"status":"ok"}`.
- `https://ansatsu.hirsuta.id/` → `200`, dengan header
  `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`.
- `/api/health` mengembalikan `Cache-Control: no-store`.
- Jalur lama `http://116.212.74.154:8080/api/health` masih `200` selama
  masa cutover.
- Setelah Fase D: `http://116.212.74.154:8080` tidak lagi bisa diakses dari
  luar, sedangkan `https://ansatsu.hirsuta.id` tetap `200`.

## Berkas

- `docker-compose.yml`
- `frontend/nginx.conf`
- `backend/main.go`
- `backend/internal/auth/auth.go`
- `.env.example`
