package latihan

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// ErrInvalidLesson dikembalikan ketika nomor lesson tidak valid.
var ErrInvalidLesson = errors.New("lesson tidak valid")

// ErrNotFound dikembalikan ketika soal latihan tidak ditemukan.
var ErrNotFound = errors.New("soal latihan tidak ditemukan")

// Repository membaca dan menulis kumpulan soal latihan bunpou.
type Repository struct {
	path string
	mu   sync.Mutex
}

// NewRepository membuat repository dengan direktori data latihan bunpou.
func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

type lessonFile struct {
	name  string
	items []SoalBunpou
}

// ListLessons mengembalikan nomor lesson yang memiliki soal latihan.
func (r *Repository) ListLessons() ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	entries, err := os.ReadDir(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	lessons := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		number, ok := normalizeLesson(entry.Name())
		if !ok {
			continue
		}
		files, err := r.readLessonFilesLocked("lesson" + number)
		if err != nil {
			return nil, err
		}
		if totalItems(files) > 0 {
			lessons = append(lessons, number)
		}
	}

	sort.Slice(lessons, func(i, j int) bool {
		a, _ := strconv.Atoi(lessons[i])
		b, _ := strconv.Atoi(lessons[j])
		return a < b
	})
	return lessons, nil
}

// FindByLesson mengembalikan seluruh soal latihan pada satu lesson.
func (r *Repository) FindByLesson(lesson string) ([]SoalBunpou, error) {
	number, ok := normalizeLesson(lesson)
	if !ok {
		return nil, ErrInvalidLesson
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	files, err := r.readLessonFilesLocked("lesson" + number)
	if err != nil {
		return nil, err
	}

	result := make([]SoalBunpou, 0)
	for _, file := range files {
		result = append(result, file.items...)
	}
	return result, nil
}

// Create menambahkan soal latihan baru ke lesson.
func (r *Repository) Create(lesson string, soal SoalBunpou) (SoalBunpou, error) {
	number, ok := normalizeLesson(lesson)
	if !ok {
		return SoalBunpou{}, ErrInvalidLesson
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	folder := "lesson" + number
	files, err := r.readLessonFilesLocked(folder)
	if err != nil {
		return SoalBunpou{}, err
	}
	soal.ID = newID()

	target := -1
	for i := range files {
		if files[i].name == "1.json" {
			target = i
			break
		}
	}
	if target == -1 {
		files = append(files, lessonFile{name: "1.json", items: []SoalBunpou{}})
		target = len(files) - 1
	}
	files[target].items = append(files[target].items, soal)
	if err := r.writeLessonFileLocked(folder, files[target]); err != nil {
		return SoalBunpou{}, err
	}
	return soal, nil
}

// Update memperbarui soal latihan berdasarkan lesson dan id.
func (r *Repository) Update(lesson, id string, soal SoalBunpou) (SoalBunpou, error) {
	number, ok := normalizeLesson(lesson)
	if !ok {
		return SoalBunpou{}, ErrInvalidLesson
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	folder := "lesson" + number
	files, err := r.readLessonFilesLocked(folder)
	if err != nil {
		return SoalBunpou{}, err
	}
	for fi := range files {
		for i := range files[fi].items {
			if files[fi].items[i].ID != id {
				continue
			}
			soal.ID = id
			files[fi].items[i] = soal
			if err := r.writeLessonFileLocked(folder, files[fi]); err != nil {
				return SoalBunpou{}, err
			}
			return soal, nil
		}
	}
	return SoalBunpou{}, ErrNotFound
}

// Delete menghapus soal latihan berdasarkan lesson dan id.
func (r *Repository) Delete(lesson, id string) error {
	number, ok := normalizeLesson(lesson)
	if !ok {
		return ErrInvalidLesson
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	folder := "lesson" + number
	files, err := r.readLessonFilesLocked(folder)
	if err != nil {
		return err
	}
	for fi := range files {
		for i := range files[fi].items {
			if files[fi].items[i].ID != id {
				continue
			}
			files[fi].items = append(files[fi].items[:i], files[fi].items[i+1:]...)
			return r.writeLessonFileLocked(folder, files[fi])
		}
	}
	return ErrNotFound
}

func (r *Repository) readLessonFilesLocked(folder string) ([]lessonFile, error) {
	dir := filepath.Join(r.path, folder)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []lessonFile{}, nil
		}
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		a, _ := strconv.Atoi(strings.TrimSuffix(entries[i].Name(), filepath.Ext(entries[i].Name())))
		b, _ := strconv.Atoi(strings.TrimSuffix(entries[j].Name(), filepath.Ext(entries[j].Name())))
		if a != 0 && b != 0 && a != b {
			return a < b
		}
		return entries[i].Name() < entries[j].Name()
	})

	files := make([]lessonFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		if len(strings.TrimSpace(string(data))) == 0 {
			continue
		}

		var items []SoalBunpou
		if err := json.Unmarshal(data, &items); err != nil {
			return nil, fmt.Errorf("file %s: %w", entry.Name(), err)
		}

		normalized := make([]SoalBunpou, 0, len(items))
		changed := false
		for _, item := range items {
			soal, ok := normalizeSoal(item)
			if !ok {
				continue
			}
			if soal.ID == "" {
				soal.ID = newID()
				changed = true
			}
			normalized = append(normalized, soal)
		}

		file := lessonFile{name: entry.Name(), items: normalized}
		if changed {
			if err := r.writeLessonFileLocked(folder, file); err != nil {
				return nil, err
			}
		}
		files = append(files, file)
	}
	return files, nil
}

func (r *Repository) writeLessonFileLocked(folder string, file lessonFile) error {
	dir := filepath.Join(r.path, folder)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(file.items, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, file.name), data, 0o644)
}

func totalItems(files []lessonFile) int {
	total := 0
	for _, file := range files {
		total += len(file.items)
	}
	return total
}

func newID() string {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "lt-000000000000"
	}
	return "lt-" + hex.EncodeToString(b)
}

func normalizeLesson(value string) (string, bool) {
	text := strings.TrimSpace(strings.ToLower(value))
	text = strings.TrimPrefix(text, "lesson")
	text = strings.TrimSpace(text)
	if text == "" {
		return "", false
	}
	for _, char := range text {
		if char < '0' || char > '9' {
			return "", false
		}
	}
	number, err := strconv.Atoi(text)
	if err != nil || number < 1 {
		return "", false
	}
	return strconv.Itoa(number), true
}

func normalizeSoal(soal SoalBunpou) (SoalBunpou, bool) {
	soal.Soal = strings.TrimSpace(soal.Soal)
	soal.TipeSoal = strings.TrimSpace(soal.TipeSoal)
	soal.JawabanBenar = strings.TrimSpace(soal.JawabanBenar)

	choices := make([]string, 0, len(soal.PilihanJawaban))
	for _, choice := range soal.PilihanJawaban {
		choice = strings.TrimSpace(choice)
		if choice != "" {
			choices = append(choices, choice)
		}
	}
	if soal.TipeSoal == "" {
		if len(choices) > 0 {
			soal.TipeSoal = "pilihan ganda"
		} else {
			soal.TipeSoal = "essay"
		}
	}
	if soal.TipeSoal == "essay" {
		choices = nil
	}
	if choices == nil {
		choices = []string{}
	}

	if soal.Soal == "" || soal.JawabanBenar == "" {
		return SoalBunpou{}, false
	}
	if soal.TipeSoal != "pilihan ganda" && soal.TipeSoal != "essay" {
		return SoalBunpou{}, false
	}
	if soal.TipeSoal == "pilihan ganda" {
		if len(choices) < 2 {
			return SoalBunpou{}, false
		}
		matched := false
		for _, choice := range choices {
			if choice == soal.JawabanBenar {
				matched = true
				break
			}
		}
		if !matched {
			return SoalBunpou{}, false
		}
	}
	soal.PilihanJawaban = choices
	return soal, true
}
