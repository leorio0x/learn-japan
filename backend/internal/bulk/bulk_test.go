package bulk

import (
	"bytes"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestParseCSV(t *testing.T) {
	data := []byte("kanji,hiragana/katakana,terjemahan,buku,lesson\n私,わたし,saya|aku,minna no nihonggo 1,1\n")
	rows, err := Parse("data.csv", data)
	if err != nil {
		t.Fatalf("parse csv: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("want 1 row, got %d", len(rows))
	}
	r := rows[0]
	if r.Kanji != "私" || r.Furigana != "わたし" {
		t.Errorf("unexpected kanji/furigana: %+v", r)
	}
	if len(r.Terjemahan) != 2 || r.Terjemahan[0] != "saya" || r.Terjemahan[1] != "aku" {
		t.Errorf("unexpected terjemahan: %+v", r.Terjemahan)
	}
	if r.Buku != "minna no nihonggo 1" || r.Lesson != "1" {
		t.Errorf("unexpected buku/lesson: %+v", r)
	}
	if r.Line != 2 {
		t.Errorf("want line 2, got %d", r.Line)
	}
}

func TestParseXLSX(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()
	sheet := f.GetSheetName(0)
	content := [][]interface{}{
		{"kanji", "hiragana/katakana", "terjemahan", "buku", "lesson"},
		{"日本", "にほん", "Jepang, Nippon", "minna no nihonggo 1", "1"},
		{},
		{"学校", "がっこう", "sekolah", "minna no nihonggo 1", "2"},
	}
	for i, row := range content {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		if err := f.SetSheetRow(sheet, cell, &row); err != nil {
			t.Fatalf("set row: %v", err)
		}
	}
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		t.Fatalf("write xlsx: %v", err)
	}

	rows, err := Parse("data.xlsx", buf.Bytes())
	if err != nil {
		t.Fatalf("parse xlsx: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("want 2 rows (baris kosong dilewati), got %d", len(rows))
	}
	if len(rows[0].Terjemahan) != 2 || rows[0].Terjemahan[1] != "Nippon" {
		t.Errorf("unexpected terjemahan: %+v", rows[0].Terjemahan)
	}
	if rows[1].Kanji != "学校" || rows[1].Lesson != "2" {
		t.Errorf("unexpected row: %+v", rows[1])
	}
}

func TestParseJSON(t *testing.T) {
	data := []byte(`[
  {
    "id": "ktb-999",
    "kanji": " 私 ",
    "furigana": " わたし ",
    "terjemahan": [" saya ", "", "aku"],
    "contoh_percakapan": [
      {
        "pertanyaan": " これは何ですか。 ",
        "arti_pertanyaan": " Ini apa? ",
        "jawaban": " それは本です。 ",
        "arti_jawaban": " Itu buku. "
      },
      {}
    ],
    "buku": "minna no nihonggo 1",
    "lesson": " 1 "
  },
  {
    "kanji": "日本",
    "furigana": "にほん",
    "terjemahan": ["Jepang", "Nippon"],
    "buku": "minna no nihonggo 2",
    "lesson": "2"
  }
]`)
	rows, err := Parse("data.json", data)
	if err != nil {
		t.Fatalf("parse json: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("want 2 rows, got %d", len(rows))
	}

	r := rows[0]
	if r.Line != 1 {
		t.Errorf("want line 1, got %d", r.Line)
	}
	if r.Kanji != "私" || r.Furigana != "わたし" {
		t.Errorf("unexpected kanji/furigana: %+v", r)
	}
	if len(r.Terjemahan) != 2 || r.Terjemahan[0] != "saya" || r.Terjemahan[1] != "aku" {
		t.Errorf("unexpected terjemahan: %+v", r.Terjemahan)
	}
	if r.Buku != "minna no nihonggo 1" || r.Lesson != "1" {
		t.Errorf("unexpected buku/lesson: %+v", r)
	}
	if len(r.ContohPercakapan) != 1 {
		t.Fatalf("want 1 percakapan, got %d", len(r.ContohPercakapan))
	}
	p := r.ContohPercakapan[0]
	if p.Pertanyaan != "これは何ですか。" || p.ArtiPertanyaan != "Ini apa?" {
		t.Errorf("unexpected pertanyaan: %+v", p)
	}
	if p.Jawaban != "それは本です。" || p.ArtiJawaban != "Itu buku." {
		t.Errorf("unexpected jawaban: %+v", p)
	}

	if rows[1].Line != 2 {
		t.Errorf("want line 2, got %d", rows[1].Line)
	}
	if len(rows[1].ContohPercakapan) != 0 {
		t.Errorf("want no percakapan, got %+v", rows[1].ContohPercakapan)
	}
}

func TestParseJSONEmpty(t *testing.T) {
	if _, err := Parse("data.json", []byte("[]")); err == nil {
		t.Fatal("want error for empty json array")
	}
}

func TestParseJSONInvalid(t *testing.T) {
	if _, err := Parse("data.json", []byte("{")); err == nil {
		t.Fatal("want error for invalid json")
	}
}

func TestParseJSONObject(t *testing.T) {
	if _, err := Parse("data.json", []byte(`{"kanji":"私"}`)); err == nil {
		t.Fatal("want error for non-array json")
	}
}

func TestParseCSVExample(t *testing.T) {
	// Contoh persis dari data/rules.md (baris data), field multi arti dibungkus kutip.
	data := []byte("kanji,hiragana,terjemahan,buku,lesson\n私,わたし,\"ini arti 1|ini arti2\",minna no nihonggo 1, 5\n")
	rows, err := Parse("data.csv", data)
	if err != nil {
		t.Fatalf("parse csv: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("want 1 row, got %d", len(rows))
	}
	r := rows[0]
	if len(r.Terjemahan) != 2 || r.Terjemahan[0] != "ini arti 1" || r.Terjemahan[1] != "ini arti2" {
		t.Errorf("unexpected terjemahan: %+v", r.Terjemahan)
	}
	if r.Lesson != "5" {
		t.Errorf("lesson harus ter-trim menjadi \"5\", got %q", r.Lesson)
	}
}

func TestParseCSVQuotedComma(t *testing.T) {
	// Field terjemahan berisi koma wajib dibungkus kutip ganda (data/rules.md).
	data := []byte("kanji,hiragana/katakana,terjemahan,buku,lesson\n失礼,そろそろ 失礼します。,\"Maaf, saya mau pamit dulu.\",Minna no Nihongo 1,8\n")
	rows, err := Parse("data.csv", data)
	if err != nil {
		t.Fatalf("parse csv: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("want 1 row, got %d", len(rows))
	}
	r := rows[0]
	if len(r.Terjemahan) != 1 || r.Terjemahan[0] != "Maaf, saya mau pamit dulu." {
		t.Errorf("unexpected terjemahan: %+v", r.Terjemahan)
	}
	if r.Buku != "Minna no Nihongo 1" || r.Lesson != "8" {
		t.Errorf("kolom bergeser, buku/lesson: %+v", r)
	}
}

func TestParseCSVHeaderMissing(t *testing.T) {
	data := []byte("kanji,buku,lesson\n私,minna 1,1\n")
	if _, err := Parse("data.csv", data); err == nil {
		t.Fatal("want error for missing terjemahan header")
	}
}

func TestParseCSVFuriganaHeader(t *testing.T) {
	data := []byte("kanji,furigana,terjemahan,buku,lesson\n日本,にほん,Jepang,minna no nihonggo 1,1\n")
	rows, err := Parse("data.csv", data)
	if err != nil {
		t.Fatalf("parse csv: %v", err)
	}
	if len(rows) != 1 || rows[0].Furigana != "にほん" {
		t.Fatalf("unexpected rows: %+v", rows)
	}
}

func TestParseUnsupportedExt(t *testing.T) {
	if _, err := Parse("data.txt", []byte("x")); err == nil {
		t.Fatal("want error for unsupported extension")
	}
}
