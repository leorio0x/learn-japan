package latihan

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepositoryLessonCRUD(t *testing.T) {
	dir := t.TempDir()
	lessonDir := filepath.Join(dir, "lesson1")
	if err := os.MkdirAll(lessonDir, 0o755); err != nil {
		t.Fatal(err)
	}
	seed := `[
  {"soal":"a（　）b","tipe_soal":"pilihan ganda","pilihan_jawaban":["は","が"],"jawaban_benar":"は"},
  {"soal":"c（　）d","tipe_soal":"essay","jawaban_benar":"です"}
]`
	if err := os.WriteFile(filepath.Join(lessonDir, "1.json"), []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}

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
	if len(list) != 2 {
		t.Fatalf("expected 2 soal, got %d", len(list))
	}
	for _, soal := range list {
		if soal.ID == "" {
			t.Fatal("id soal kosong")
		}
	}

	created, err := repo.Create("1", SoalBunpou{
		Soal:           "e（　）f",
		TipeSoal:       "essay",
		JawabanBenar:   "ます",
		PilihanJawaban: []string{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == "" {
		t.Fatal("id soal baru kosong")
	}

	updated, err := repo.Update("1", created.ID, SoalBunpou{
		Soal:           "e（　）g",
		TipeSoal:       "essay",
		JawabanBenar:   "ました",
		PilihanJawaban: []string{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Soal != "e（　）g" {
		t.Fatalf("update gagal: %+v", updated)
	}

	if err := repo.Delete("1", created.ID); err != nil {
		t.Fatal(err)
	}
	list, err = repo.FindByLesson("1")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 soal setelah delete, got %d", len(list))
	}

	if _, err := repo.FindByLesson("abc"); err != ErrInvalidLesson {
		t.Fatalf("expected ErrInvalidLesson, got %v", err)
	}
}
