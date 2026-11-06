# Terjemahan Dokkai Tidak Memisahkan Paragraf

- **Tanggal**: 2026-11-06
- **Komponen**: Dokkai (frontend `Dokkai`, `DokkaiView`, `CrudDokkai`)

## Ringkasan

Terjemahan dokkai tampil, tetapi karakter baris baru ganda (`\n\n`) pada
field `terjemahan` tidak menghasilkan jarak antar paragraf. Seluruh teks
terjemahan terlihat menyatu dalam satu blok.

## Penyebab

Nilai `terjemahan` dirender di dalam elemen `<p>` tanpa kelas
`whitespace-pre-line`. Secara default HTML menciutkan whitespace, termasuk
baris baru, sehingga `\n\n` tidak terlihat. Paragraf bahasa Jepang sudah
memakai `whitespace-pre-line`, sedangkan blok terjemahannya belum.

## Solusi

1. Menambahkan kelas `whitespace-pre-line` pada elemen `<p>` terjemahan di
   `frontend/src/pages/Dokkai.svelte`.
2. Menambahkan kelas yang sama pada
   `frontend/src/components/DokkaiView.svelte`.
3. Menambahkan kelas yang sama pada preview terjemahan di
   `frontend/src/pages/CrudDokkai.svelte` agar konsisten.

## Verifikasi

- `bun run build` pada frontend berhasil.
- Kelas `whitespace-pre-line` terkonfirmasi ada pada ketiga komponen dan
  ikut terkompilasi di CSS hasil build.

## Berkas

- `frontend/src/pages/Dokkai.svelte`
- `frontend/src/components/DokkaiView.svelte`
- `frontend/src/pages/CrudDokkai.svelte`
