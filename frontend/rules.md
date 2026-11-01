# Tech stack

framework = svelte 5.16.*
runtime = bun 1.4.*
style = tailwindcss 3.4.17.*

buat agar responsif untuk mobile, tablet dan desktop

# Warna

untuk warna navbar, tombol, text, gunakan warna dibawah, buat agar indah dan tidak aneh

- #0D1B2A
- #1B263B
- #415A77
- #778D7A
- #D4C4A8
- #F4F1DE

## Navbar

- opsi hide
- auto hide ketika scroll

# Features

- tombol navigasi ke halaman fitur crud&tebak2an disimpan di masing masing halaman, misal
  CRUD kotoba disimpan di halaman kotoba, begitupun fitur crud&tebak2an lain
- tombol navigasi ke halaman fitur crud disembunyikan jika belum memasukkan password
- logo untuk tombol navigasi fitur crud,tebak2an,dokkai
- untuk setiap fitur CRUD, jangan beri hak akses kecuali memasukkan password: armanganteng123#
- halaman untuk memasukkan password terdapat di url /ansatsunoadmin
- password tidak boleh terlihat melalui fitur browser

check detail rules features di :
frontend/rules/*.md
