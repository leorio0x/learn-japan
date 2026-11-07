## 1. Halaman bunpou

- display bunpou dengan opsi [lesson, semua lesson]
- data bunpou diambil dari data/bunpou.json
- pagination
- opsi hide kanji
- pada attribute rumus ditulis dalam format sesuai aturan yaitu [～:orang|kata benda|kata sifat|kata tunjuk],
  yang ditampilkan hanya "~" saja, dimana "~" adalah sesuatu yang bisa diisi/diganti dengan
  orang,kata benda,kata sifat,kata tunjuk dll.
- setiap simbol ~ pada rumus bersifat interaktif (bisa diketuk atau diarahkan kursornya),
  dimana sebuah komponen Floating UI (seperti tooltip, popover, atau modal kecil) akan
  muncul menampilkan daftar kategori/penjelas yang valid untuk mengisi slot tersebut
  (misalnya: orang, kata benda, kata sifat, kata tunjuk).
- data paragraf yang mengandung format <kanji:kana:arti>,
  kana dan arti muncul ketika klik 1 kali, gunakan "floating-ui" untuk fitur ini
- trigger dismiss toast dari floating-ui ketika klik dimanapun
- tombol navigasi untuk ke halaman latihan bunpou
- tombol navigasi untuk ke halaman CRUD bunpou (muncul ketika sudah input password)

## 1.1. Halaman latihan bunpou

data pertanyaan diambil dari `data/latihan/bunpou/*`
misal bunpou lesson 1, diambil dari `data/latihan/bunpou/lesson1/*`

- jawaban bisa pilihan ganda atau essay (tergantung data)
-

## 1.2. Halaman CRUD bunpou

sesuaikan input dengan format `data/format/soal_bunpou.json`
