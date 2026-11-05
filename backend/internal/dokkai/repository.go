package dokkai

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Repository menyimpan data dokkai pada direktori:
// <root>/<tipe>/<lesson>.json
type Repository struct {
	root string
	mu   sync.RWMutex
}

// NewRepository membuat repository dengan root direktori data dokkai.
func NewRepository(root string) *Repository {
	return &Repository{root: root}
}

// ValidKind mengecek apakah tipe dokkai didukung.
func ValidKind(kind string) bool {
	return kind == "kotoba" || kind == "kanji"
}

func sanitizeLesson(lesson string) string {
	lesson = strings.TrimSpace(lesson)
	var b strings.Builder
	for _, r := range lesson {
		switch {
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r == '-' || r == '_':
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (r *Repository) filePath(kind, lesson string) (string, error) {
	if !ValidKind(kind) {
		return "", ErrInvalidKind
	}
	name := sanitizeLesson(lesson)
	if name == "" {
		return "", ErrInvalidLesson
	}
	return filepath.Join(r.root, kind, name+".json"), nil
}

func (r *Repository) readLocked(kind, lesson string) ([]Dokkai, error) {
	path, err := r.filePath(kind, lesson)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Dokkai{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []Dokkai{}, nil
	}
	var list []Dokkai
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []Dokkai{}
	}
	return list, nil
}

func (r *Repository) writeLocked(kind, lesson string, list []Dokkai) error {
	path, err := r.filePath(kind, lesson)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

// Lessons mengembalikan daftar lesson yang memiliki file dokkai pada tipe tertentu.
func (r *Repository) Lessons(kind string) ([]string, error) {
	if !ValidKind(kind) {
		return nil, ErrInvalidKind
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	entries, err := os.ReadDir(filepath.Join(r.root, kind))
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	lessons := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || strings.ToLower(filepath.Ext(entry.Name())) != ".json" {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		if name == "" || sanitizeLesson(name) != name {
			continue
		}
		lessons = append(lessons, name)
	}

	sort.Slice(lessons, func(i, j int) bool {
		return compareLesson(lessons[i], lessons[j]) < 0
	})
	return lessons, nil
}

func compareLesson(a, b string) int {
	na, errA := strconv.Atoi(a)
	nb, errB := strconv.Atoi(b)
	switch {
	case errA == nil && errB == nil:
		switch {
		case na < nb:
			return -1
		case na > nb:
			return 1
		default:
			return 0
		}
	case errA == nil:
		return -1
	case errB == nil:
		return 1
	}
	return strings.Compare(a, b)
}

// FindByLesson mengembalikan seluruh dokkai untuk tipe dan lesson tertentu.
func (r *Repository) FindByLesson(kind, lesson string) ([]Dokkai, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.readLocked(kind, lesson)
}

// AppendMany menambahkan beberapa dokkai sekaligus, dikelompokkan per lesson.
func (r *Repository) AppendMany(kind string, items []Dokkai) (int, error) {
	if !ValidKind(kind) {
		return 0, ErrInvalidKind
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	grouped := make(map[string][]Dokkai)
	order := make([]string, 0)
	for _, item := range items {
		lesson := sanitizeLesson(item.Lesson)
		if lesson == "" {
			return 0, ErrInvalidLesson
		}
		if _, ok := grouped[lesson]; !ok {
			order = append(order, lesson)
		}
		item.Lesson = lesson
		grouped[lesson] = append(grouped[lesson], item)
	}

	appended := 0
	for _, lesson := range order {
		existing, err := r.readLocked(kind, lesson)
		if err != nil {
			return appended, err
		}

		next := nextDokkaiID(existing)
		for i := range grouped[lesson] {
			grouped[lesson][i].ID = strconv.Itoa(next)
			next++
		}

		existing = append(existing, grouped[lesson]...)
		if err := r.writeLocked(kind, lesson, existing); err != nil {
			return appended, err
		}
		appended += len(grouped[lesson])
	}
	return appended, nil
}

func nextDokkaiID(items []Dokkai) int {
	max := 0
	for _, item := range items {
		if n, err := strconv.Atoi(item.ID); err == nil && n > max {
			max = n
		}
	}
	return max + 1
}
