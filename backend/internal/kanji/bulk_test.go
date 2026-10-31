package kanji

import (
	"os"
	"path/filepath"
	"testing"

	"learn-japan/backend/internal/bulk"
)

func TestBulkFromJSONGeneratesIDAndKeepsPercakapan(t *testing.T) {
	path := filepath.Join(t.TempDir(), "kanji.json")
	existing := `[{"id":"knj-007","kanji":"日","furigana":"ひ","terjemahan":["hari"],"buku":"minna no nihonggo 1","lesson":"1"}]`
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	rows, err := bulk.Parse("data.json", []byte(`[
		{
			"id": "knj-999",
			"kanji": "本",
			"furigana": "ほん",
			"terjemahan": ["buku"],
			"contoh_percakapan": [
				{"pertanyaan": "本ですか。", "arti_pertanyaan": "Buku?", "jawaban": "はい。", "arti_jawaban": "Ya."}
			],
			"buku": "minna no nihonggo 1",
			"lesson": "1"
		}
	]`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	created, err := NewService(NewRepository(path)).Bulk(rows)
	if err != nil {
		t.Fatalf("bulk: %v", err)
	}
	if len(created) != 1 {
		t.Fatalf("want 1 created, got %d", len(created))
	}
	if created[0].ID != "knj-008" {
		t.Errorf("want generated id knj-008, got %q", created[0].ID)
	}
	if len(created[0].ContohPercakapan) != 1 {
		t.Fatalf("percakapan hilang: %+v", created[0])
	}
	if created[0].ContohPercakapan[0].Pertanyaan != "本ですか。" {
		t.Errorf("unexpected percakapan: %+v", created[0].ContohPercakapan[0])
	}
}
