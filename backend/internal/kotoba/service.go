package kotoba

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// Service berisi logika bisnis untuk fitur kotoba.
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
		return "ktb-000000000000"
	}
	return "ktb-" + hex.EncodeToString(b)
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

func toKotoba(req Request) Kotoba {
	return Kotoba{
		Kanji:            strings.TrimSpace(req.Kanji),
		Furigana:         strings.TrimSpace(req.Furigana),
		Terjemahan:       cleanTerjemahan(req.Terjemahan),
		ContohPercakapan: cleanPercakapan(req.ContohPercakapan),
		Buku:             strings.TrimSpace(req.Buku),
		Lesson:           strings.TrimSpace(req.Lesson),
	}
}

func (s *Service) List() ([]Kotoba, error) {
	return s.repo.FindAll()
}

func (s *Service) Get(id string) (Kotoba, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Create(req Request) (Kotoba, error) {
	k := toKotoba(req)
	k.ID = newID()
	return s.repo.Create(k)
}

func (s *Service) Update(id string, req Request) (Kotoba, error) {
	return s.repo.Update(id, toKotoba(req))
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

// BulkDelete menghapus beberapa kotoba sekaligus.
func (s *Service) BulkDelete(ids []string) (int, error) {
	return s.repo.DeleteMany(ids)
}
