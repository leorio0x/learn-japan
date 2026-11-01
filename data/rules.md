# Format data

format untuk kotoba.json :
./format/kotoba.json

format untuk kanji.json :
./format/kanji.json

format untuk kaiwa.json :
./format/kaiwa.json

format untuk dokkai.json :
./format/dokkai.json

# Format Header

format header dari input bulk excel dan csv kotoba :
kanji,furigana,terjemahan,buku,lesson

format header dari input bulk excel dan csv kanji :
kanji,furigana,terjemahan,buku,lesson

# Format Bulk Input

format bulk input file json (kanji, kotoba) :

- setiap data dari input bulk tanpa attribute id, id digenerate dan
  disesuaikan dengan data terakhir dari existing data

format bulk input file .csv, .xlsx (kanji, kotoba) :

- kolom csv dipisahkan dengan koma (,)
- field csv yang mengandung koma wajib dibungkus tanda kutip ganda ("...")
  contoh : kanji,furigana,"Maaf, saya mau pamit dulu.",minna no nihonggo 1,8
- setiap data dari input bulk, id akan digenerate dan
  disesuaikan dengan data terakhir dari existing data
- setiap terjemahan multi arti delimiternya berupa "|" untuk csv,
  ditulis di dalam field yang dibungkus tanda kutip ganda
  contoh : kanji,furigana,"ini arti 1|ini arti2",minna no nihonggo 1, 5
- setiap terjemahan multi arti delimiternya berupa "," untuk excel
- tiap arti di-trim, entri kosong dibuang, minimal 1 arti
- nilai "buku" mengikuti pilihan [minna no nihonggo 1, minna no nihonggo 2]

format bulk input file json dokkai :

- setiap data dari input bulk, id akan digenerate dan
  disesuaikan dengan data terakhir dari existing data
