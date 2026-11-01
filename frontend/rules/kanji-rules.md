## 1 Halaman kanji

- display kanji dengan opsi [per lesson, semua]
- default display = lesson 1
- opsi display [card, table, flashcard(seperti quizlet)]
- display table&card berupa [Kanji, furigana, arti, contoh_percakapan]
- untuk display flashcard front side berupa kanji saja, backside berupa furigana dan artinya (furigana diatas, arti dibawahnya)
- pagination hanya berlaku ketika display kanji = semua
- setiap data contoh_percakapan yang mengandung format <kanji:furigana>
  tolong ubah agak furigananya berada diatas kanji sebagaimana mestinya
- di setiap opsi tampilan card berikan tombol untuk "lihat detail" yang akan
  trigger modal box, tombol tersebut terletak di bagian paling bawah kiri dari card
  isi dari modal box adalah seluruh data dari kotobanya, buat serapih mungkin
- opsi untuk hide kanji, furigana, arti, contoh_percakapan

### 1.1 Halaman Dokkai

- pada halaman ini terdapat dokkai berupa cerita pendek yang memuat seluruh kanji perlessonnya.
  data diambil dari data/dokkai/kanji/[lesson].json
- data paragraf berformat <kanji:kana:arti>, dimana kana dan
  arti muncul ketika klik 2 kali, gunakan "floating-ui" untuk fitur ini
- trigger dismiss toast dari floating-ui ketika klik dimanapun

### 1.2 Halaman tebak kanji

- buat agar ketika scroll kebawah, position dari bagian soal (blok kanji) = sticky

### 1.3 Halaman CRUD kanji

- input "buku" tipe select dengan pilihan [minna no nihonggo 1, minna no nihonggo 2]
- setiap kanji bisa terdapat 1 - >1 arti, pastikan untuk handle ini juga
- display kanji buatkan pagination per 10 item
- keterangan total jumlah kanji
- opsi bulk input (csv,excel,json) yang ditransformasi ke format data/format/kanji.json
- opsi untuk bulk delete
- input contoh contoh_percakapan bisa lebih dari 1
