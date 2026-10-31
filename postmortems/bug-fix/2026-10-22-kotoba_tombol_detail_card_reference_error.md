# Tombol "Lihat detail" pada card kotoba memakai variabel yang salah

- **Tanggal**: 2026-10-22
- **Komponen**: Frontend — Halaman Kotoba (tampilan card)

## Ringkasan

Pada tampilan card halaman Kotoba, tombol **"Lihat detail"** tidak membuka modal
detail. Handler mengacu ke variabel `k` yang tidak ada di halaman Kotoba
(variabel loop di sana bernama `w`), sehingga `detailItem` tidak terisi data
kotoba yang dipilih.

## Penyebab

Sisa salin-tempel dari `Kanji.svelte`:

```js
onclick={() => (detailItem = k)}
```

Sedangkan loop pada Kotoba memakai `{#each paged as w (w.id)}`.

## Solusi

`frontend/src/pages/Kotoba.svelte` — ganti referensi ke item loop yang benar:

```js
onclick={() => (detailItem = w)}
```

## Verifikasi

- `bun run build` berhasil.
- Tombol detail pada card Kotoba kini membuka modal dengan data yang dipilih.
