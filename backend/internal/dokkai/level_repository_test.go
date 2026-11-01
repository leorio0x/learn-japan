package dokkai

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLevelRepositoryBulkGeneratesIDAndListsLevels(t *testing.T) {
	root := t.TempDir()
	n5 := filepath.Join(root, "level", "N5")
	if err := os.MkdirAll(n5, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(n5, "1.json"), []byte(`{"id":"1","level":"N5","judul":"A","paragraf":"a"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(n5, "3.json"), []byte(`{"id":"3","level":"N5","judul":"B","paragraf":"b"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	repo := NewRepository(root)
	svc := NewService(repo)

	imported, err := svc.BulkLevel("", []LevelDokkai{
		{Level: "N5", Judul: "C", Paragraf: "c"},
		{Level: "N4", Judul: "D", Paragraf: "d"},
	})
	if err != nil {
		t.Fatalf("bulk: %v", err)
	}
	if imported != 2 {
		t.Fatalf("want 2 imported, got %d", imported)
	}

	n5Stories, err := svc.ByLevel("n5")
	if err != nil {
		t.Fatalf("by level: %v", err)
	}
	if len(n5Stories) != 3 {
		t.Fatalf("want 3 stories, got %d", len(n5Stories))
	}
	if n5Stories[2].ID != "4" {
		t.Errorf("want generated id 4, got %q", n5Stories[2].ID)
	}
	if n5Stories[0].ID != "1" || n5Stories[1].ID != "3" {
		t.Errorf("stories not sorted by id: %+v", n5Stories)
	}

	levels, err := svc.Levels()
	if err != nil {
		t.Fatalf("levels: %v", err)
	}
	if len(levels) != 2 || levels[0] != "N5" || levels[1] != "N4" {
		t.Errorf("unexpected levels order: %v", levels)
	}
}

func TestLevelServiceUpdateAndDelete(t *testing.T) {
	root := t.TempDir()
	repo := NewRepository(root)
	svc := NewService(repo)

	if _, err := svc.BulkLevel("N5", []LevelDokkai{{Judul: "A", Paragraf: "a"}}); err != nil {
		t.Fatalf("bulk: %v", err)
	}
	updated, err := svc.UpdateLevel("N5", "1", LevelDokkai{Judul: "B", Paragraf: "b"})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.ID != "1" || updated.Level != "N5" || updated.Judul != "B" {
		t.Errorf("unexpected updated story: %+v", updated)
	}

	if err := svc.DeleteLevel("N5", "1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	stories, err := svc.ByLevel("N5")
	if err != nil {
		t.Fatalf("by level: %v", err)
	}
	if len(stories) != 0 {
		t.Errorf("want empty after delete, got %+v", stories)
	}
}
