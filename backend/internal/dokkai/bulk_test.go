package dokkai

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBulkDokkaiGeneratesIDFromExisting(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "kanji")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	existing := `[{"id":"5","lesson":"1","judul":"A","paragraf":"a"}]`
	if err := os.WriteFile(filepath.Join(dir, "1.json"), []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewService(NewRepository(root))
	imported, err := svc.Bulk("kanji", []Dokkai{
		{ID: "99", Lesson: "1", Judul: "B", Paragraf: "b"},
		{ID: "100", Lesson: "1", Judul: "C", Paragraf: "c"},
	})
	if err != nil {
		t.Fatalf("bulk: %v", err)
	}
	if imported != 2 {
		t.Fatalf("want 2 imported, got %d", imported)
	}

	stories, err := svc.ByLesson("kanji", "1")
	if err != nil {
		t.Fatalf("by lesson: %v", err)
	}
	if len(stories) != 3 {
		t.Fatalf("want 3 stories, got %d", len(stories))
	}
	if stories[1].ID != "6" || stories[2].ID != "7" {
		t.Errorf("unexpected generated ids: %+v", stories)
	}
}
