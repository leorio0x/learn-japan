package kanji

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// Filter berisi kriteria penyaringan data kanji.
type Filter struct {
	Lesson string
	Buku   string
}

// Service berisi logika bisnis untuk fitur kanji.
type Service struct {
	repo *Repository
}

// NewService membuat service baru.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func newID() string {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "knj-000000000000"
	}
	return "knj-" + hex.EncodeToString(b)
}

// cleanTerjemahan membersihkan daftar arti: trim dan buang entri kosong.
func cleanTerjemahan(items []string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		if v := strings.TrimSpace(item); v != "" {
			result = append(result, v)
		}
	}
	return result
}

// cleanPercakapan membersihkan contoh percakapan: trim dan buang entri kosong.
func cleanPercakapan(items []Percakapan) []Percakapan {
	result := make([]Percakapan, 0, len(items))
	for _, item := range items {
		p := Percakapan{
			Pertanyaan:     strings.TrimSpace(item.Pertanyaan),
			ArtiPertanyaan: strings.TrimSpace(item.ArtiPertanyaan),
			Jawaban:        strings.TrimSpace(item.Jawaban),
			ArtiJawaban:    strings.TrimSpace(item.ArtiJawaban),
		}
		if p == (Percakapan{}) {
			continue
		}
		result = append(result, p)
	}
	return result
}

func toKanji(req Request) Kanji {
	return Kanji{
		Kanji:            strings.TrimSpace(req.Kanji),
		Furigana:         strings.TrimSpace(req.Furigana),
		Terjemahan:       cleanTerjemahan(req.Terjemahan),
		ContohPercakapan: cleanPercakapan(req.ContohPercakapan),
		Buku:             strings.TrimSpace(req.Buku),
		Lesson:           strings.TrimSpace(req.Lesson),
	}
}

// List mengembalikan data kanji yang sudah disaring.
func (s *Service) List(f Filter) ([]Kanji, error) {
	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	lesson := strings.TrimSpace(f.Lesson)
	buku := strings.ToLower(strings.TrimSpace(f.Buku))

	result := make([]Kanji, 0, len(list))
	for _, k := range list {
		if lesson != "" && k.Lesson != lesson {
			continue
		}
		if buku != "" && !strings.Contains(strings.ToLower(k.Buku), buku) {
			continue
		}
		result = append(result, k)
	}
	return result, nil
}

func (s *Service) Get(id string) (Kanji, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Create(req Request) (Kanji, error) {
	k := toKanji(req)
	k.ID = newID()
	return s.repo.Create(k)
}

func (s *Service) Update(id string, req Request) (Kanji, error) {
	return s.repo.Update(id, toKanji(req))
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

// BulkDelete menghapus beberapa kanji sekaligus.
func (s *Service) BulkDelete(ids []string) (int, error) {
	return s.repo.DeleteMany(ids)
}
