package dokkai

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestLevelRoutes(t *testing.T) {
	repo := NewRepository(t.TempDir())
	svc := NewService(repo)
	if _, err := svc.BulkLevel("N5", []LevelDokkai{{Judul: "A", Paragraf: "a"}}); err != nil {
		t.Fatal(err)
	}

	app := fiber.New()
	RegisterRoutes(app, repo, func(c *fiber.Ctx) error { return c.Next() })

	req := httptest.NewRequest("GET", "/dokkai/level", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	var levels []string
	if err := json.Unmarshal(body, &levels); err != nil {
		t.Fatalf("unmarshal levels: %v (%s)", err, body)
	}
	if len(levels) != 1 || levels[0] != "N5" {
		t.Fatalf("unexpected levels: %v", levels)
	}

	req = httptest.NewRequest("GET", "/dokkai/level/N5", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(resp.Body)
	var stories []LevelDokkai
	if err := json.Unmarshal(body, &stories); err != nil {
		t.Fatalf("unmarshal stories: %v (%s)", err, body)
	}
	if len(stories) != 1 || stories[0].ID != "1" || stories[0].Level != "N5" {
		t.Fatalf("unexpected stories: %+v", stories)
	}
}
