package dokkai

import (
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidKind dikembalikan ketika tipe dokkai bukan kotoba/kanji.
var ErrInvalidKind = errors.New("tipe dokkai tidak valid (gunakan kotoba atau kanji)")

// ErrInvalidLesson dikembalikan ketika lesson kosong atau tidak valid.
var ErrInvalidLesson = errors.New("lesson wajib diisi")

// Service berisi logika bisnis untuk fitur dokkai.
type Service struct {
	repo *Repository
}

// NewService membuat service baru.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Lessons mengembalikan daftar lesson yang tersedia untuk tipe dokkai tertentu.
func (s *Service) Lessons(kind string) ([]string, error) {
	if !ValidKind(kind) {
		return nil, ErrInvalidKind
	}
	return s.repo.Lessons(kind)
}

// ByLesson mengembalikan dokkai berdasarkan tipe dan lesson.
func (s *Service) ByLesson(kind, lesson string) ([]Dokkai, error) {
	if !ValidKind(kind) {
		return nil, ErrInvalidKind
	}
	lesson = strings.TrimSpace(lesson)
	if lesson == "" {
		return nil, ErrInvalidLesson
	}
	return s.repo.FindByLesson(kind, lesson)
}

// Bulk memvalidasi lalu menyimpan dokkai dari input bulk.
func (s *Service) Bulk(kind string, items []Dokkai) (int, error) {
	if !ValidKind(kind) {
		return 0, ErrInvalidKind
	}
	if len(items) == 0 {
		return 0, errors.New("tidak ada data yang bisa diimpor")
	}

	cleaned := make([]Dokkai, 0, len(items))
	for i, item := range items {
		item.Lesson = strings.TrimSpace(item.Lesson)
		item.Judul = strings.TrimSpace(item.Judul)
		item.Paragraf = strings.TrimSpace(item.Paragraf)
		item.Terjemahan = strings.TrimSpace(item.Terjemahan)
		if item.Lesson == "" {
			return 0, fmt.Errorf("data %d: lesson wajib diisi", i+1)
		}
		if item.Paragraf == "" {
			return 0, fmt.Errorf("data %d: paragraf wajib diisi", i+1)
		}
		cleaned = append(cleaned, item)
	}
	return s.repo.AppendMany(kind, cleaned)
}
