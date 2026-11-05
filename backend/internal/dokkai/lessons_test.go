package dokkai

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRepositoryLessons(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "kanji")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"1.json", "2.json", "10.json", "1-15.json", "catatan.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(`[]`), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	repo := NewRepository(root)
	lessons, err := repo.Lessons("kanji")
	if err != nil {
		t.Fatalf("lessons: %v", err)
	}
	want := []string{"1", "2", "10", "1-15"}
	if len(lessons) != len(want) {
		t.Fatalf("want %v, got %v", want, lessons)
	}
	for i := range want {
		if lessons[i] != want[i] {
			t.Fatalf("want %v, got %v", want, lessons)
		}
	}

	if _, err := repo.Lessons("invalid"); !errors.Is(err, ErrInvalidKind) {
		t.Fatalf("want ErrInvalidKind, got %v", err)
	}
}

func TestLessonsRoute(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "kanji")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "1.json"), []byte(`[]`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "1-15.json"), []byte(`[]`), 0o644); err != nil {
		t.Fatal(err)
	}

	repo := NewRepository(root)
	app := fiber.New()
	RegisterRoutes(app, repo, func(c *fiber.Ctx) error { return c.Next() })

	req := httptest.NewRequest("GET", "/dokkai/kanji/lessons", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	var lessons []string
	if err := json.Unmarshal(body, &lessons); err != nil {
		t.Fatalf("unmarshal lessons: %v (%s)", err, body)
	}
	if len(lessons) != 2 || lessons[0] != "1" || lessons[1] != "1-15" {
		t.Fatalf("unexpected lessons: %v", lessons)
	}

	req = httptest.NewRequest("GET", "/dokkai/invalid/lessons", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("want 400, got %d", resp.StatusCode)
	}
}
