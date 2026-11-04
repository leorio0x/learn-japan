package kotoba

import (
	"crypto/rand"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"
)

// Filter berisi kriteria penyaringan data kotoba.
type Filter struct {
	Lesson string
}

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

// List mengembalikan data kotoba yang sudah disaring.
func (s *Service) List(f Filter) ([]Kotoba, error) {
	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	lesson := strings.TrimSpace(f.Lesson)
	if lesson == "" {
		return list, nil
	}

	result := make([]Kotoba, 0, len(list))
	for _, k := range list {
		if k.Lesson == lesson {
			result = append(result, k)
		}
	}
	return result, nil
}

// Lessons mengembalikan daftar lesson unik yang tersedia, terurut numerik.
func (s *Service) Lessons() ([]string, error) {
	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{}, len(list))
	result := make([]string, 0)
	for _, k := range list {
		lesson := strings.TrimSpace(k.Lesson)
		if lesson == "" {
			continue
		}
		if _, ok := seen[lesson]; ok {
			continue
		}
		seen[lesson] = struct{}{}
		result = append(result, lesson)
	}

	sort.Slice(result, func(i, j int) bool {
		ni, errI := strconv.Atoi(result[i])
		nj, errJ := strconv.Atoi(result[j])
		if errI == nil && errJ == nil && ni != nj {
			return ni < nj
		}
		return result[i] < result[j]
	})
	return result, nil
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
