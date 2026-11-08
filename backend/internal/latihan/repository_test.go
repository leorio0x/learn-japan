package latihan

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepositoryCategorizedLessons(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "lesson1", "jtest", "1.json"), `[
  {"lesson":"1","kategori":"jtest","instruksi":"susun kalimat","soal":"[1. わたし 2. がくせい ]です","jawaban":"わたしはがくせいです"}
]`)
	writeTestFile(t, filepath.Join(dir, "lesson1", "pg", "1.json"), `[
  {"lesson":"1","kategori":"pilihan ganda","soal":"<私:わたし:saya> ___ <学生:がくせい:mahasiswa> です","pilihan_jawaban":["が","は","を","に"],"jawaban_benar":"b"}
]`)
	writeTestFile(t, filepath.Join(dir, "lesson1", "dokkai", "1.json"), `[
  {
    "lesson":"1",
    "kategori":"dokkai",
    "dokkai":"わたし は <田中:たなか:Tanaka> です。",
    "soal":[
      {"no":"1","pertanyaan":"<田中:たなか:Tanaka> さん は <何:なん:apa> ですか。","pilihan_jawaban":["<学生:がくせい:mahasiswa>","<会社員:かいしゃいん:pegawai perusahaan>"],"jawaban_benar":"b"},
      {"no":"2","pertanyaan":"サントスさん は <学生:がくせい:mahasiswa> ですか。","pilihan_jawaban":["はい","いいえ"],"jawaban_benar":"a"}
    ]
  }
]`)

	repo := NewRepository(dir)
	lessons, err := repo.ListLessons()
	if err != nil {
		t.Fatal(err)
	}
	if len(lessons) != 1 || lessons[0] != "1" {
		t.Fatalf("unexpected lessons: %v", lessons)
	}

	list, err := repo.FindByLesson("1")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 4 {
		t.Fatalf("expected 4 soal, got %d", len(list))
	}
	for _, soal := range list {
		if soal.ID == "" {
			t.Fatal("id soal kosong")
		}
		if soal.Lesson != "1" {
			t.Fatalf("lesson tidak sesuai: %+v", soal)
		}
		if soal.Kategori == "jtest" && soal.JawabanBenar == "" {
			t.Fatalf("jawaban jtest kosong: %+v", soal)
		}
		if soal.Kategori == "dokkai" && soal.Dokkai == "" {
			t.Fatalf("dokkai kosong: %+v", soal)
		}
	}

	created, err := repo.Create("1", Request{
		Kategori:       "pg",
		Soal:           "ini soal baru",
		PilihanJawaban: []string{"a", "b"},
		JawabanBenar:   "a",
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.Kategori != "pg" {
		t.Fatalf("create gagal: %+v", created)
	}

	updated, err := repo.Update("1", created.ID, Request{
		Kategori:       "pg",
		Soal:           "ini soal diubah",
		PilihanJawaban: []string{"a", "b"},
		JawabanBenar:   "b",
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Soal != "ini soal diubah" || updated.JawabanBenar != "b" {
		t.Fatalf("update gagal: %+v", updated)
	}

	if err := repo.Delete("1", created.ID); err != nil {
		t.Fatal(err)
	}
	list, err = repo.FindByLesson("1")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 4 {
		t.Fatalf("expected 4 soal setelah delete, got %d", len(list))
	}

	if _, err := repo.FindByLesson("abc"); err != ErrInvalidLesson {
		t.Fatalf("expected ErrInvalidLesson, got %v", err)
	}
}

func TestRepositoryDokkaiCRUD(t *testing.T) {
	dir := t.TempDir()
	repo := NewRepository(dir)

	created, err := repo.Create("2", Request{
		Kategori:       "dokkai",
		Dokkai:         "ini cerita dokkai",
		Soal:           "pertanyaan pertama",
		PilihanJawaban: []string{"a", "b", "c", "d"},
		JawabanBenar:   "a",
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.Kategori != "dokkai" || created.Dokkai != "ini cerita dokkai" {
		t.Fatalf("create dokkai gagal: %+v", created)
	}

	updated, err := repo.Update("2", created.ID, Request{
		Kategori:       "dokkai",
		Dokkai:         "cerita diperbarui",
		Soal:           "pertanyaan kedua",
		PilihanJawaban: []string{"a", "b", "c", "d"},
		JawabanBenar:   "c",
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Soal != "pertanyaan kedua" || updated.Dokkai != "cerita diperbarui" {
		t.Fatalf("update dokkai gagal: %+v", updated)
	}

	if err := repo.Delete("2", created.ID); err != nil {
		t.Fatal(err)
	}
	list, err := repo.FindByLesson("2")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Fatalf("expected dokkai kosong, got %d", len(list))
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
