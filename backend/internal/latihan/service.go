package latihan

import (
	"errors"
	"strings"
)

// ErrInvalidRequest dikembalikan ketika payload soal tidak valid.
var ErrInvalidRequest = errors.New("data soal tidak valid")

// Service berisi logika bisnis untuk latihan bunpou.
type Service struct {
	repo *Repository
}

// NewService membuat service baru.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ListLessons mengembalikan nomor lesson yang memiliki soal.
func (s *Service) ListLessons() ([]string, error) {
	return s.repo.ListLessons()
}

// GetByLesson mengembalikan soal latihan untuk lesson tertentu.
func (s *Service) GetByLesson(lesson string) ([]SoalBunpou, error) {
	return s.repo.FindByLesson(lesson)
}

// Create menambahkan soal latihan baru.
func (s *Service) Create(req Request) (SoalBunpou, error) {
	soal, err := buildSoal(req)
	if err != nil {
		return SoalBunpou{}, err
	}
	return s.repo.Create(req.Lesson, soal)
}

// Update memperbarui soal latihan.
func (s *Service) Update(lesson, id string, req Request) (SoalBunpou, error) {
	soal, err := buildSoal(req)
	if err != nil {
		return SoalBunpou{}, err
	}
	return s.repo.Update(lesson, id, soal)
}

// Delete menghapus soal latihan.
func (s *Service) Delete(lesson, id string) error {
	return s.repo.Delete(lesson, id)
}

func buildSoal(req Request) (SoalBunpou, error) {
	soal := SoalBunpou{
		Soal:           strings.TrimSpace(req.Soal),
		TipeSoal:       strings.TrimSpace(req.TipeSoal),
		PilihanJawaban: req.PilihanJawaban,
		JawabanBenar:   strings.TrimSpace(req.JawabanBenar),
	}
	normalized, ok := normalizeSoal(soal)
	if !ok {
		return SoalBunpou{}, ErrInvalidRequest
	}
	return normalized, nil
}
