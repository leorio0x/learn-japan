## 1. Halaman kotoba

- display kotoba dengan opsi [per lesson, semua]
- default display = lesson 1
- opsi display [card, table, flashcard(seperti quizlet)]
- untuk display flashcard front side berupa kanji dan furigana saja, backside berupa artinya saja
- display table&card berupa [Kanji, furigana, arti, contoh_percakapan]
- setiap data contoh_percakapan yang mengandung format <kanji:furigana>
  tolong ubah agak furigananya berada diatas kanji sebagaimana mestinya
- di setiap opsi tampilan card berikan tombol untuk "lihat detail" yang akan
  trigger modal box, tombol tersebut terletak di bagian paling bawah kiri dari card
  isi dari modal box adalah seluruh data dari kotobanya, buat serapih mungkin
- berikan juga tombol "lihat detail" di setiap baris data dengan opsi tampilan
  table seperti pada opsi tampilan card, akan tetapi letak tombolnya di kolom paling kanan
- opsi sorting data yang ditampilkan
  [default (ini berarti sesuai urutan data yang diambil dulu), acak]
- pagination hanya berlaku ketika display kotoba = semua
- opsi untuk hide kanji, furigana, arti, contoh_percakapan

### 1.1 Halaman dokkai

- pada halaman ini terdapat dokkai berupa cerita pendek yang memuat seluruh kotoba perlessonnya.
  data diambil dari data/dokkai/kotoba/[lesson].json
- opsi untuk hide kanji
- data paragraf berformat <kanji:kana:arti>, dimana kana dan
  arti muncul ketika klik 2 kali, gunakan floating-ui untuk fitur ini
- trigger dismiss toast dari floating-ui ketika klik dimanapun
- Format paragraf tanpa spasi antar kata, mengikuti gaya penulisan bahasa Jepang

### 1.2 Halaman tebak kotoba

- pada halaman ini berikan toggle untuk hide kanji sehingga
  hanya memunculkan furigana saja.
- setiap pertanyaan yang berupa tanggal/mengandung angka, otomatis sembunyikan kanjinya
- toggle sembunyikan kanji tampil sepanjang pertanyaan, state tidak berubah sampai diubah secara manual
- opsi memilih jumlah pertanyaan perbatch tebak-tebakan
- keterangan total jumlah kotoba
- pada pertanyaan dengan opsi pilihan ganda, susun 4 pilihan (A-D) yang
  terdiri dari 1 jawaban benar + 3 distraktor. distraktor diambil dari
  data terdekat sebelum/sesudah kotoba yang ditanyakan (bukan acak).
  jika data sekitar tidak cukup, ambil dari data lain yang tersisa.

### 1.3 Halaman CRUD kotoba

- input "buku" tipe select dengan pilihan [minna no nihonggo 1, minna no nihonggo 2]
- opsi bulk input (csv,excel,json) yang ditransformasi ke format data/format/kotoba.json
- opsi untuk bulk delete
- input contoh contoh_percakapan bisa lebih dari 1
