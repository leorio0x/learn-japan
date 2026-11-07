# Ringkasan App

app untuk quiz kotoba, kanji dan bunpou bahasa jepang, referensi materi aplikasi ini dari buku minna no nihongo 1&2

# Backend

untuk spesifikasi tech stack dan detail lainnya lihat di backend/rules.md

# Frontend

untuk spesifikasi tech stack dan detail lainnya lihat di frontend/rules.md

# Data

untuk format datanya lihat di data/rules.md
untuk data kotoba lihat di data/kotoba.json
untuk data kanji lihat di data/kanji.json
untuk data bunpou lihat di data/bunpou.json

# Containerization

gunakan docker untuk membungkus app backend&frontend

# postmortem

setiap fixing bug tuliskan di folder postmortems/bug-fix dengan format YYYY-MM-DD-nama_insiden.md
setiap fixing performance tuliskan di folder postmortems/performance-fix dengan format YYYY-MM-DD-nama_performance_issue.md

# Features

## Tebak Kotoba

- pertanyaan bisa berupa aksara jepang dan jawaban berupa bahasa indonesia atau sebaliknya.
- jawaban bisa berupa pilihan ganda atau essay.
- bisa sortir quiz berdasarkan lesson.
- setiap batch pertanyaan ada akhir dan skornya.
- jika pertanyaan berupa aksara jepang, diatas kanjinya tuliskan furigana (untuk hint).

## Tebak Kanji

- pertanyaan berupa kanji dari data/kanji.json
- jawaban berupa penulisan hiragana (furigana) dan artinya dalam bentuk essay atau bisa dalam bentuk pilihan ganda
- bisa sortir quiz berdasarkan lesson
- setiap batch pertanyaan ada akhir dan skornya

## Tebak bunpou

- tebak pertanyaan dari sebuah jawaban
- tebak jawaban dari suatu pertanyaan

## Dokkai

- halaman dokkai berdasarkan level [N5,N4,dst]

## Bulk input

- fitur bulk input (csv,excel,json) untuk kotoba & kanji, dikonversi
  menjadi json sesuai format di data/rules.md
- fitur bulk input dokkai dengan format json, dikonversi agar sesuai dengan
  format json di data/rules.md
