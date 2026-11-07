package bunpou

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// Filter berisi kriteria penyaringan data bunpou.
type Filter struct {
	Lesson string
	Buku   string
}

// Service berisi logika bisnis untuk fitur bunpou.
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
		return "bp-000000000000"
	}
	return "bp-" + hex.EncodeToString(b)
}

func toKalimat(k *Kalimat) *Kalimat {
	if k == nil {
		return nil
	}
	return &Kalimat{
		Kalimat:    strings.TrimSpace(k.Kalimat),
		Terjemahan: strings.TrimSpace(k.Terjemahan),
	}
}

func toContoh(list []Contoh) []Contoh {
	out := make([]Contoh, 0, len(list))
	for _, c := range list {
		out = append(out, Contoh{
			Pertanyaan: toKalimat(c.Pertanyaan),
			Jawaban:    toKalimat(c.Jawaban),
		})
	}
	return out
}

func toPola(list []Pola) []Pola {
	out := make([]Pola, 0, len(list))
	for _, p := range list {
		out = append(out, Pola{
			NoBunpou:  strings.TrimSpace(p.NoBunpou),
			NamaPola:  strings.TrimSpace(p.NamaPola),
			Rumus:     strings.TrimSpace(p.Rumus),
			Deskripsi: strings.TrimSpace(p.Deskripsi),
			Contoh:    toContoh(p.Contoh),
		})
	}
	return out
}

func toBunpou(req Request) Bunpou {
	return Bunpou{
		Lesson: strings.TrimSpace(req.Lesson),
		Buku:   strings.TrimSpace(req.Buku),
		Pola:   toPola(req.Pola),
	}
}

// List mengembalikan data bunpou yang sudah disaring.
func (s *Service) List(f Filter) ([]Bunpou, error) {
	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	lesson := strings.TrimSpace(f.Lesson)
	buku := strings.ToLower(strings.TrimSpace(f.Buku))

	result := make([]Bunpou, 0, len(list))
	for _, b := range list {
		if lesson != "" && b.Lesson != lesson {
			continue
		}
		if buku != "" && !strings.Contains(strings.ToLower(b.Buku), buku) {
			continue
		}
		result = append(result, b)
	}
	return result, nil
}

func (s *Service) Get(id string) (Bunpou, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Create(req Request) (Bunpou, error) {
	b := toBunpou(req)
	b.ID = newID()
	return s.repo.Create(b)
}

func (s *Service) Update(id string, req Request) (Bunpou, error) {
	return s.repo.Update(id, toBunpou(req))
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

// BulkDelete menghapus beberapa bunpou sekaligus.
func (s *Service) BulkDelete(ids []string) (int, error) {
	return s.repo.DeleteMany(ids)
}
