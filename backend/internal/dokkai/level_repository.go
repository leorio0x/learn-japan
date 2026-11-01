package dokkai

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var (
	// ErrInvalidLevel dikembalikan ketika level tidak valid.
	ErrInvalidLevel = errors.New("level tidak valid")
	// ErrInvalidStoryID dikembalikan ketika id dokkai tidak valid.
	ErrInvalidStoryID = errors.New("id dokkai tidak valid")
	// ErrStoryNotFound dikembalikan ketika dokkai tidak ditemukan.
	ErrStoryNotFound = errors.New("dokkai tidak ditemukan")
)

func sanitizeLevel(level string) string {
	level = strings.ToUpper(strings.TrimSpace(level))
	if level == "" {
		return ""
	}
	for _, r := range level {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			continue
		}
		return ""
	}
	return level
}

func sanitizeStoryID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return ""
	}
	for _, r := range id {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r == '-' || r == '_':
		default:
			return ""
		}
	}
	return id
}

func (r *Repository) levelDir(level string) string {
	return filepath.Join(r.root, "level", level)
}

func (r *Repository) levelStoryPath(level, id string) (string, error) {
	level = sanitizeLevel(level)
	if level == "" {
		return "", ErrInvalidLevel
	}
	id = sanitizeStoryID(id)
	if id == "" {
		return "", ErrInvalidStoryID
	}
	return filepath.Join(r.levelDir(level), id+".json"), nil
}

// Levels mengembalikan daftar level yang tersedia diurutkan dari N5, N4, dst.
func (r *Repository) Levels() ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entries, err := os.ReadDir(filepath.Join(r.root, "level"))
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	levels := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		level := sanitizeLevel(entry.Name())
		if level == "" {
			continue
		}
		levels = append(levels, level)
	}

	sort.Slice(levels, func(i, j int) bool {
		ni, oki := levelNumber(levels[i])
		nj, okj := levelNumber(levels[j])
		if oki && okj && ni != nj {
			return ni > nj
		}
		return levels[i] < levels[j]
	})
	return levels, nil
}

func levelNumber(level string) (int, bool) {
	if len(level) < 2 || level[0] != 'N' {
		return 0, false
	}
	n, err := strconv.Atoi(level[1:])
	if err != nil {
		return 0, false
	}
	return n, true
}

// FindByLevel mengembalikan seluruh dokkai pada level tertentu.
func (r *Repository) FindByLevel(level string) ([]LevelDokkai, error) {
	level = sanitizeLevel(level)
	if level == "" {
		return nil, ErrInvalidLevel
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.readLevelLocked(level)
}

func (r *Repository) readLevelLocked(level string) ([]LevelDokkai, error) {
	dir := r.levelDir(level)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []LevelDokkai{}, nil
		}
		return nil, err
	}

	list := make([]LevelDokkai, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || strings.ToLower(filepath.Ext(entry.Name())) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		if len(data) == 0 {
			continue
		}

		var story LevelDokkai
		if err := json.Unmarshal(data, &story); err != nil {
			return nil, fmt.Errorf("gagal membaca %s: %w", entry.Name(), err)
		}
		if story.ID == "" {
			story.ID = strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		}
		if story.Level == "" {
			story.Level = level
		}
		list = append(list, story)
	}

	sort.Slice(list, func(i, j int) bool {
		return compareStoryID(list[i].ID, list[j].ID) < 0
	})
	return list, nil
}

func compareStoryID(a, b string) int {
	na, errA := strconv.Atoi(a)
	nb, errB := strconv.Atoi(b)
	if errA == nil && errB == nil {
		switch {
		case na < nb:
			return -1
		case na > nb:
			return 1
		default:
			return 0
		}
	}
	return strings.Compare(a, b)
}

// AppendLevels menambahkan beberapa dokkai berlevel sekaligus.
// ID dibuat otomatis mengikuti data terakhir pada level masing-masing.
func (r *Repository) AppendLevels(defaultLevel string, items []LevelDokkai) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	grouped := make(map[string][]LevelDokkai)
	order := make([]string, 0)
	for _, item := range items {
		level := sanitizeLevel(item.Level)
		if level == "" {
			level = sanitizeLevel(defaultLevel)
		}
		if level == "" {
			return 0, ErrInvalidLevel
		}
		item.Level = level
		if _, ok := grouped[level]; !ok {
			order = append(order, level)
		}
		grouped[level] = append(grouped[level], item)
	}

	appended := 0
	for _, level := range order {
		existing, err := r.readLevelLocked(level)
		if err != nil {
			return appended, err
		}
		next := nextStoryID(existing)
		for _, item := range grouped[level] {
			item.ID = strconv.Itoa(next)
			next++
			if err := r.writeLevelLocked(level, item); err != nil {
				return appended, err
			}
			appended++
		}
	}
	return appended, nil
}

func nextStoryID(items []LevelDokkai) int {
	max := 0
	for _, item := range items {
		if n, err := strconv.Atoi(item.ID); err == nil && n > max {
			max = n
		}
	}
	return max + 1
}

func (r *Repository) writeLevelLocked(level string, story LevelDokkai) error {
	path, err := r.levelStoryPath(level, story.ID)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(story, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

// UpdateLevel menimpa dokkai berlevel berdasarkan level dan id.
func (r *Repository) UpdateLevel(level, id string, story LevelDokkai) (LevelDokkai, error) {
	level = sanitizeLevel(level)
	id = sanitizeStoryID(id)
	if level == "" {
		return LevelDokkai{}, ErrInvalidLevel
	}
	if id == "" {
		return LevelDokkai{}, ErrInvalidStoryID
	}
	path, err := r.levelStoryPath(level, id)
	if err != nil {
		return LevelDokkai{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return LevelDokkai{}, ErrStoryNotFound
		}
		return LevelDokkai{}, err
	}

	story.ID = id
	story.Level = level
	if err := r.writeLevelLocked(level, story); err != nil {
		return LevelDokkai{}, err
	}
	return story, nil
}

// DeleteLevel menghapus dokkai berlevel berdasarkan level dan id.
func (r *Repository) DeleteLevel(level, id string) error {
	level = sanitizeLevel(level)
	id = sanitizeStoryID(id)
	if level == "" {
		return ErrInvalidLevel
	}
	if id == "" {
		return ErrInvalidStoryID
	}
	path, err := r.levelStoryPath(level, id)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return ErrStoryNotFound
		}
		return err
	}
	return nil
}
