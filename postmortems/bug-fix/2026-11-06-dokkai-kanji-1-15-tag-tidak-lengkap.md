# Dokkai Kanji 1-15: Tag Kata Tidak Lengkap

- **Tanggal**: 2026-11-06
- **Komponen**: Data dokkai kanji `data/dokkai/kanji/1-15.json`

## Ringkasan

Sebagian kata pada dokkai kumulatif kanji lesson `1-15` tidak memiliki tooltip
floating-ui. Kanji masih ada yang ditulis tanpa tag `<kanji:kana:arti>`, dan
sebagian kata beraksara kana masih ditulis apa adanya tanpa tag `<:kana:arti>`.

## Penyebab

`parseParagraph` hanya mengubah teks di dalam `<...>` menjadi komponen
`DokkaiWord` yang memiliki tooltip. Teks di luar tag dirender sebagai teks biasa,
sehingga kata yang tidak dianotasi tidak memunculkan toast furigana/arti.

Saat data dibuat, aturan format `<kanji:kana:arti>` dan `<:kana:arti>` belum
diterapkan secara merata. Akibatnya muncul kata seperti `とうきょう`, `えき`,
`やさしい`, `あります`, `友達`, `先生`, `本`, dan `電車` yang masih telanjang.

## Solusi

1. Menganotasi seluruh kanji yang masih telanjang dengan format
   `<kanji:kana:arti>`.
2. Menganotasi kata kana yang masih telanjang dengan format `<:kana:arti>`.
3. Membiarkan partikel/gramatika (`は`, `を`, `に`, `で`, `と`, `も`, `の`,
   `が`, `です`, `でした`, `さん`, dan sejenisnya) tetap tanpa tag agar kalimat
   tetap terbaca wajar.
4. Mempertahankan teks Jepang yang terlihat dan seluruh `terjemahan`; hanya
   menambahkan tag anotasi.

## Verifikasi

- Tidak ada lagi segmen di luar tag yang mengandung kanji.
- Segmen di luar tag hanya berisi partikel, kopula, dan sufiks.
- Hasil render tag (kanji atau kana jika kanji kosong) identik dengan teks
  sebelum perbaikan.
- Seluruh tag memiliki format `kanji:kana:arti` atau `:kana:arti`.
- File tetap JSON valid.

## Berkas

- `data/dokkai/kanji/1-15.json`
