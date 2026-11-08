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
func (s *Service) GetByLesson(lesson string) ([]Soal, error) {
	return s.repo.FindByLesson(lesson)
}

// Create menambahkan soal latihan baru.
func (s *Service) Create(req Request) (Soal, error) {
	if err := validateRequest(req); err != nil {
		return Soal{}, err
	}
	return s.repo.Create(req.Lesson, req)
}

// Update memperbarui soal latihan.
func (s *Service) Update(lesson, id string, req Request) (Soal, error) {
	req.Lesson = lesson
	if err := validateRequest(req); err != nil {
		return Soal{}, err
	}
	return s.repo.Update(lesson, id, req)
}

// Delete menghapus soal latihan.
func (s *Service) Delete(lesson, id string) error {
	return s.repo.Delete(lesson, id)
}

func validateRequest(req Request) error {
	lesson := strings.TrimSpace(req.Lesson)
	if _, ok := normalizeLesson(lesson); !ok {
		return ErrInvalidLesson
	}

	kategori, ok := normalizeKategori(req.Kategori)
	if !ok {
		return ErrInvalidKategori
	}

	soal := strings.TrimSpace(req.Soal)
	if soal == "" {
		return ErrInvalidRequest
	}

	switch kategori {
	case "jtest":
		if strings.TrimSpace(req.Jawaban) == "" {
			return ErrInvalidRequest
		}
	case "pg", "bintang":
		if !validChoices(req.PilihanJawaban, req.JawabanBenar) {
			return ErrInvalidRequest
		}
	case "dokkai":
		if strings.TrimSpace(req.Dokkai) == "" {
			return ErrInvalidRequest
		}
		if !validChoices(req.PilihanJawaban, req.JawabanBenar) {
			return ErrInvalidRequest
		}
	}

	return nil
}

func validChoices(choices []string, answer string) bool {
	cleaned := trimChoices(choices)
	if len(cleaned) < 2 {
		return false
	}
	index := answerIndex(answer)
	return index >= 0 && index < len(cleaned)
}

func answerIndex(value string) int {
	text := strings.ToLower(strings.TrimSpace(value))
	if len(text) != 1 {
		return -1
	}
	char := text[0]
	if char < 'a' || char > 'z' {
		return -1
	}
	return int(char - 'a')
}
