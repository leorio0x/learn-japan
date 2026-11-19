## 1 Halaman kanji

- display kanji dengan opsi [per lesson, semua]
- default display = lesson 1
- opsi display [card, table, flashcard(seperti quizlet)]
- display card berupa [Kanji, furigana, arti, contoh_percakapan]
- sedangkan display table hanya berupa [Kanji, furigana, arti]
- untuk display flashcard front side berupa kanji saja, backside berupa furigana dan artinya (furigana diatas, arti dibawahnya)
- pagination hanya berlaku ketika display kanji = semua
- setiap data contoh_percakapan yang mengandung format <kanji:furigana>
  tolong ubah agak furigananya berada diatas kanji sebagaimana mestinya
- di setiap opsi tampilan card berikan tombol untuk "lihat detail" yang akan
  trigger modal box, tombol tersebut terletak di bagian paling bawah kiri dari card
  isi dari modal box adalah seluruh data dari kotobanya, buat serapih mungkin
- opsi untuk hide kanji, furigana, arti, contoh_percakapan
- fitur search berdasarkan [kanji, furigana, arti]
- setiap kanji punya tombol toggle "tandai lupa" di tampilan card, table, dan flashcard
- kanji yang ditandai ditampilkan pada panel ringkasan berisi jumlah serta daftar
  kanji, furigana, dan arti, dengan aksi hapus per item dan hapus semua
- tombol "latihan kanji ditandai" pada panel mengarah ke halaman tebak kanji dengan
  sumber soal "kanji ditandai"
- tombol latihan nonaktif jika belum ada kanji yang ditandai

### 1.1 Halaman Dokkai

- pada halaman ini terdapat dokkai berupa cerita pendek yang memuat seluruh kanji perlessonnya.
  data diambil dari data/dokkai/kanji/[lesson].json
- data paragraf berformat <kanji:kana:arti>, dimana kana dan
  arti muncul ketika klik 2 kali, gunakan "floating-ui" untuk fitur ini
- trigger dismiss toast dari floating-ui ketika klik dimanapun
- Format paragraf tanpa spasi antar kata, mengikuti gaya penulisan bahasa Jepang

### 1.2 Penanda kanji lupa (sementara)

- tombol tandai berupa toggle (tandai lupa / batal tandai)
- kanji yang ditandai hanya disimpan sementara di frontend menggunakan sessionStorage,
  tanpa backend dan tanpa mengubah data kanji
- yang disimpan cukup id kanji, data lengkap diambil dari data yang sudah dimuat
- penanda hilang otomatis ketika tab/browser ditutup
- id kanji yang tidak ditemukan pada data saat dimuat diabaikan
- kanji yang ditandai bisa dihapus satu per satu atau seluruhnya dari panel ringkasan
- kanji yang ditandai dapat dijadikan soal melalui halaman tebak kanji
  (sumber soal: kanji ditandai)

### 1.3 Halaman tebak kanji

- buat agar ketika scroll kebawah, position dari bagian soal (blok kanji) = sticky
- ketika essay, input auto focus di input cara baca (furigana), ini hanya berlaku ketika device=komputer
- opsi filter atau pemilihan multi-lesson untuk sesi tebak kanji (gunakan CSS GRID, 4 col perbaris)
- opsi sumber soal [semua kanji, kanji ditandai]
- ketika sumber = kanji ditandai: filter lesson tidak dipakai, jumlah soal default = semua
  kanji yang ditandai
- opsi pilihan ganda tetap mengambil pilihan pengecoh dari seluruh data kanji
- tampilkan empty state jika belum ada kanji yang ditandai

### 1.4 Halaman CRUD kanji

- input "buku" tipe select dengan pilihan [minna no nihonggo 1, minna no nihonggo 2]
- setiap kanji bisa terdapat 1 - >1 arti, pastikan untuk handle ini juga
- display kanji buatkan pagination per 10 item
- keterangan total jumlah kanji
- opsi bulk input (csv,excel,json) yang ditransformasi ke format data/format/kanji.json
- opsi untuk bulk delete
- input contoh contoh_percakapan bisa lebih dari 1
