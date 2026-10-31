package kaiwa

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// Filter berisi kriteria penyaringan data kaiwa.
type Filter struct {
	Lesson string
	Buku   string
}

// Service berisi logika bisnis untuk fitur kaiwa.
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
		return "kw-000000000000"
	}
	return "kw-" + hex.EncodeToString(b)
}

func toDialog(d Dialog) Dialog {
	return Dialog{
		Kanji:      strings.TrimSpace(d.Kanji),
		Kana:       strings.TrimSpace(d.Kana),
		Terjemahan: strings.TrimSpace(d.Terjemahan),
	}
}

func toDialogs(list []Dialog) []Dialog {
	out := make([]Dialog, 0, len(list))
	for _, d := range list {
		out = append(out, toDialog(d))
	}
	return out
}

func toContoh(list []Contoh) []Contoh {
	out := make([]Contoh, 0, len(list))
	for _, c := range list {
		out = append(out, Contoh{
			Pertanyaan: toDialogs(c.Pertanyaan),
			Jawaban:    toDialogs(c.Jawaban),
		})
	}
	return out
}

func toPola(list []Pola) []Pola {
	out := make([]Pola, 0, len(list))
	for _, p := range list {
		out = append(out, Pola{
			PolaPertanyaan: strings.TrimSpace(p.PolaPertanyaan),
			PolaJawaban:    strings.TrimSpace(p.PolaJawaban),
			Deskripsi:      strings.TrimSpace(p.Deskripsi),
			Contoh:         toContoh(p.Contoh),
		})
	}
	return out
}

func toKaiwa(req Request) Kaiwa {
	return Kaiwa{
		Pola:   toPola(req.Pola),
		Buku:   strings.TrimSpace(req.Buku),
		Lesson: strings.TrimSpace(req.Lesson),
	}
}

// List mengembalikan data kaiwa yang sudah disaring.
func (s *Service) List(f Filter) ([]Kaiwa, error) {
	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	lesson := strings.TrimSpace(f.Lesson)
	buku := strings.ToLower(strings.TrimSpace(f.Buku))

	result := make([]Kaiwa, 0, len(list))
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

func (s *Service) Get(id string) (Kaiwa, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Create(req Request) (Kaiwa, error) {
	k := toKaiwa(req)
	k.ID = newID()
	return s.repo.Create(k)
}

func (s *Service) Update(id string, req Request) (Kaiwa, error) {
	return s.repo.Update(id, toKaiwa(req))
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

// BulkDelete menghapus beberapa kaiwa sekaligus.
func (s *Service) BulkDelete(ids []string) (int, error) {
	return s.repo.DeleteMany(ids)
}
