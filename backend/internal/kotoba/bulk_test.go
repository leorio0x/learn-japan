package kotoba

import (
	"os"
	"path/filepath"
	"testing"

	"learn-japan/backend/internal/bulk"
)

func TestBulkFromJSONGeneratesIDAndKeepsPercakapan(t *testing.T) {
	path := filepath.Join(t.TempDir(), "kotoba.json")
	existing := `[{"id":"ktb-007","kanji":"私","furigana":"わたし","terjemahan":["saya"],"buku":"minna no nihonggo 1","lesson":"1"}]`
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	rows, err := bulk.Parse("data.json", []byte(`[
		{
			"id": "ktb-999",
			"kanji": "日本",
			"furigana": "にほん",
			"terjemahan": ["Jepang"],
			"contoh_percakapan": [
				{"pertanyaan": "日本ですか。", "arti_pertanyaan": "Jepang?", "jawaban": "はい。", "arti_jawaban": "Ya."}
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
	if created[0].ID != "ktb-008" {
		t.Errorf("want generated id ktb-008, got %q", created[0].ID)
	}
	if len(created[0].ContohPercakapan) != 1 {
		t.Fatalf("percakapan hilang: %+v", created[0])
	}
	if created[0].ContohPercakapan[0].Pertanyaan != "日本ですか。" {
		t.Errorf("unexpected percakapan: %+v", created[0].ContohPercakapan[0])
	}
}
