# Tombol "Lihat detail" Kotoba: ReferenceError k is not defined

- **Tanggal**: 2026-10-21
- **Komponen**: Halaman Kotoba, mode card (`frontend/src/pages/Kotoba.svelte`)

## Ringkasan

Pada mode tampilan card, tombol "Lihat detail" melempar error di browser:

```
Uncaught ReferenceError: k is not defined
    at HTMLButtonElement.on_click_5 (Kotoba.svelte:231:58)
```

Klik tombol tidak membuka modal detail.

## Penyebab

Handler tombol memakai variabel `k` yang tidak ada:

```svelte
{#each paged as w (w.id)}
  ...
  onclick={() => (detailItem = k)}
```

Variabel item pada blok `{#each}` bernama `w`, bukan `k`. Akibatnya terjadi
`ReferenceError` saat klik. Mode tabel pada file yang sama sudah benar memakai
`w`.

## Solusi

Ganti referensi `k` menjadi `w` pada handler tombol mode card:

```svelte
onclick={() => (detailItem = w)}
```

## Verifikasi

- `Kotoba.svelte` tidak lagi memiliki referensi ke `k`.
- Tombol "Lihat detail" mode card mengisi `detailItem` dengan item yang dipilih
  sehingga modal detail tampil.

## Berkas

- `frontend/src/pages/Kotoba.svelte`
