package dokkai

import (
	"errors"
	"fmt"
	"strings"
)

// Levels mengembalikan daftar level dokkai yang tersedia.
func (s *Service) Levels() ([]string, error) {
	return s.repo.Levels()
}

// ByLevel mengembalikan dokkai berlevel berdasarkan level.
func (s *Service) ByLevel(level string) ([]LevelDokkai, error) {
	level = sanitizeLevel(level)
	if level == "" {
		return nil, ErrInvalidLevel
	}
	return s.repo.FindByLevel(level)
}

// BulkLevel memvalidasi lalu menyimpan dokkai berlevel dari input bulk.
// Level per data dipakai bila ada, jika tidak memakai defaultLevel.
func (s *Service) BulkLevel(defaultLevel string, items []LevelDokkai) (int, error) {
	if len(items) == 0 {
		return 0, errors.New("tidak ada data yang bisa diimpor")
	}

	defaultLevel = sanitizeLevel(defaultLevel)
	cleaned := make([]LevelDokkai, 0, len(items))
	for i, item := range items {
		level := sanitizeLevel(item.Level)
		if level == "" {
			level = defaultLevel
		}
		if level == "" {
			return 0, fmt.Errorf("data %d: level wajib diisi", i+1)
		}

		item.Level = level
		item.Judul = strings.TrimSpace(item.Judul)
		item.Paragraf = strings.TrimSpace(item.Paragraf)
		item.Terjemahan = strings.TrimSpace(item.Terjemahan)
		if item.Paragraf == "" {
			return 0, fmt.Errorf("data %d: paragraf wajib diisi", i+1)
		}
		cleaned = append(cleaned, item)
	}
	return s.repo.AppendLevels(defaultLevel, cleaned)
}

// UpdateLevel memperbarui satu dokkai berlevel.
func (s *Service) UpdateLevel(level, id string, item LevelDokkai) (LevelDokkai, error) {
	level = sanitizeLevel(level)
	if level == "" {
		return LevelDokkai{}, ErrInvalidLevel
	}
	if sanitizeStoryID(id) == "" {
		return LevelDokkai{}, ErrInvalidStoryID
	}

	item.Judul = strings.TrimSpace(item.Judul)
	item.Paragraf = strings.TrimSpace(item.Paragraf)
	item.Terjemahan = strings.TrimSpace(item.Terjemahan)
	if item.Paragraf == "" {
		return LevelDokkai{}, errors.New("paragraf wajib diisi")
	}
	return s.repo.UpdateLevel(level, id, item)
}

// DeleteLevel menghapus satu dokkai berlevel.
func (s *Service) DeleteLevel(level, id string) error {
	level = sanitizeLevel(level)
	if level == "" {
		return ErrInvalidLevel
	}
	if sanitizeStoryID(id) == "" {
		return ErrInvalidStoryID
	}
	return s.repo.DeleteLevel(level, id)
}
