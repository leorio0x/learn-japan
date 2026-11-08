package latihan

import (
	"bytes"
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

// ErrInvalidKategori dikembalikan ketika kategori soal tidak didukung.
var ErrInvalidKategori = errors.New("kategori soal tidak valid")

// ErrNotFound dikembalikan ketika soal latihan tidak ditemukan.
var ErrNotFound = errors.New("soal latihan tidak ditemukan")

// KategoriList adalah urutan kategori soal latihan yang didukung.
var KategoriList = []string{"jtest", "pg", "bintang", "dokkai"}

// Repository membaca dan menulis kumpulan soal latihan bunpou.
// Struktur data: <root>/lesson<N>/<kategori>/*.json
type Repository struct {
	path string
	mu   sync.Mutex
}

// NewRepository membuat repository dengan direktori data latihan bunpou.
func NewRepository(path string) *Repository {
	return &Repository{path: path}
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
		items, err := r.readLessonLocked(number)
		if err != nil {
			return nil, err
		}
		if len(items) > 0 {
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
func (r *Repository) FindByLesson(lesson string) ([]Soal, error) {
	number, ok := normalizeLesson(lesson)
	if !ok {
		return nil, ErrInvalidLesson
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	return r.readLessonLocked(number)
}

// Create menambahkan soal latihan baru ke lesson.
func (r *Repository) Create(lesson string, req Request) (Soal, error) {
	number, ok := normalizeLesson(lesson)
	if !ok {
		return Soal{}, ErrInvalidLesson
	}
	kategori, ok := normalizeKategori(req.Kategori)
	if !ok {
		return Soal{}, ErrInvalidKategori
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	path := r.filePath(number, kategori, "1")
	itemIndex := 0

	switch kategori {
	case "jtest":
		items, err := readJSONFile[rawJtest](path)
		if err != nil {
			return Soal{}, err
		}
		item := rawJtest{
			Lesson:    number,
			Kategori:  kategori,
			Instruksi: strings.TrimSpace(req.Instruksi),
			Soal:      strings.TrimSpace(req.Soal),
			Jawaban:   strings.TrimSpace(req.Jawaban),
		}
		items = append(items, item)
		if err := writeJSONFile(path, items); err != nil {
			return Soal{}, err
		}
		itemIndex = len(items) - 1
		return jtestToSoal(number, kategori, "1", itemIndex, item), nil

	case "pg", "bintang":
		items, err := readJSONFile[rawPilihan](path)
		if err != nil {
			return Soal{}, err
		}
		item := rawPilihan{
			Lesson:         number,
			Kategori:       kategori,
			Soal:           strings.TrimSpace(req.Soal),
			PilihanJawaban: trimChoices(req.PilihanJawaban),
			JawabanBenar:   strings.TrimSpace(req.JawabanBenar),
		}
		items = append(items, item)
		if err := writeJSONFile(path, items); err != nil {
			return Soal{}, err
		}
		itemIndex = len(items) - 1
		return pilihanToSoal(number, kategori, "1", itemIndex, item), nil

	case "dokkai":
		items, err := readJSONFile[rawDokkai](path)
		if err != nil {
			return Soal{}, err
		}
		no := strings.TrimSpace(req.No)
		if no == "" {
			no = "1"
		}
		item := rawDokkai{
			Lesson:   number,
			Kategori: kategori,
			Dokkai:   strings.TrimSpace(req.Dokkai),
			Soal: []rawDokkaiSoal{{
				No:             no,
				Pertanyaan:     strings.TrimSpace(req.Soal),
				PilihanJawaban: trimChoices(req.PilihanJawaban),
				JawabanBenar:   strings.TrimSpace(req.JawabanBenar),
			}},
		}
		items = append(items, item)
		if err := writeJSONFile(path, items); err != nil {
			return Soal{}, err
		}
		itemIndex = len(items) - 1
		return dokkaiToSoal(number, kategori, "1", itemIndex, 0, item, item.Soal[0]), nil
	}

	return Soal{}, ErrInvalidKategori
}

// Update memperbarui soal latihan berdasarkan lesson dan id.
func (r *Repository) Update(lesson, id string, req Request) (Soal, error) {
	number, ok := normalizeLesson(lesson)
	if !ok {
		return Soal{}, ErrInvalidLesson
	}
	loc, ok := parseID(number, id)
	if !ok {
		return Soal{}, ErrNotFound
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	path := r.filePath(number, loc.kategori, loc.fileBase)

	switch loc.kategori {
	case "jtest":
		items, err := readJSONFile[rawJtest](path)
		if err != nil {
			return Soal{}, err
		}
		if loc.itemIndex < 0 || loc.itemIndex >= len(items) {
			return Soal{}, ErrNotFound
		}
		item := items[loc.itemIndex]
		item.Lesson = number
		item.Kategori = loc.kategori
		item.Instruksi = strings.TrimSpace(req.Instruksi)
		item.Soal = strings.TrimSpace(req.Soal)
		item.Jawaban = strings.TrimSpace(req.Jawaban)
		items[loc.itemIndex] = item
		if err := writeJSONFile(path, items); err != nil {
			return Soal{}, err
		}
		return jtestToSoal(number, loc.kategori, loc.fileBase, loc.itemIndex, item), nil

	case "pg", "bintang":
		items, err := readJSONFile[rawPilihan](path)
		if err != nil {
			return Soal{}, err
		}
		if loc.itemIndex < 0 || loc.itemIndex >= len(items) {
			return Soal{}, ErrNotFound
		}
		item := rawPilihan{
			Lesson:         number,
			Kategori:       loc.kategori,
			Soal:           strings.TrimSpace(req.Soal),
			PilihanJawaban: trimChoices(req.PilihanJawaban),
			JawabanBenar:   strings.TrimSpace(req.JawabanBenar),
		}
		items[loc.itemIndex] = item
		if err := writeJSONFile(path, items); err != nil {
			return Soal{}, err
		}
		return pilihanToSoal(number, loc.kategori, loc.fileBase, loc.itemIndex, item), nil

	case "dokkai":
		items, err := readJSONFile[rawDokkai](path)
		if err != nil {
			return Soal{}, err
		}
		if loc.itemIndex < 0 || loc.itemIndex >= len(items) {
			return Soal{}, ErrNotFound
		}
		story := items[loc.itemIndex]
		story.Lesson = number
		story.Kategori = loc.kategori
		if text := strings.TrimSpace(req.Dokkai); text != "" {
			story.Dokkai = text
		}
		if loc.subIndex < 0 || loc.subIndex >= len(story.Soal) {
			return Soal{}, ErrNotFound
		}
		sub := story.Soal[loc.subIndex]
		if no := strings.TrimSpace(req.No); no != "" {
			sub.No = no
		}
		sub.Pertanyaan = strings.TrimSpace(req.Soal)
		sub.PilihanJawaban = trimChoices(req.PilihanJawaban)
		sub.JawabanBenar = strings.TrimSpace(req.JawabanBenar)
		story.Soal[loc.subIndex] = sub
		items[loc.itemIndex] = story
		if err := writeJSONFile(path, items); err != nil {
			return Soal{}, err
		}
		return dokkaiToSoal(number, loc.kategori, loc.fileBase, loc.itemIndex, loc.subIndex, story, sub), nil
	}

	return Soal{}, ErrNotFound
}

// Delete menghapus soal latihan berdasarkan lesson dan id.
func (r *Repository) Delete(lesson, id string) error {
	number, ok := normalizeLesson(lesson)
	if !ok {
		return ErrInvalidLesson
	}
	loc, ok := parseID(number, id)
	if !ok {
		return ErrNotFound
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	path := r.filePath(number, loc.kategori, loc.fileBase)

	switch loc.kategori {
	case "jtest":
		items, err := readJSONFile[rawJtest](path)
		if err != nil {
			return err
		}
		if loc.itemIndex < 0 || loc.itemIndex >= len(items) {
			return ErrNotFound
		}
		items = append(items[:loc.itemIndex], items[loc.itemIndex+1:]...)
		return writeJSONFile(path, items)

	case "pg", "bintang":
		items, err := readJSONFile[rawPilihan](path)
		if err != nil {
			return err
		}
		if loc.itemIndex < 0 || loc.itemIndex >= len(items) {
			return ErrNotFound
		}
		items = append(items[:loc.itemIndex], items[loc.itemIndex+1:]...)
		return writeJSONFile(path, items)

	case "dokkai":
		items, err := readJSONFile[rawDokkai](path)
		if err != nil {
			return err
		}
		if loc.itemIndex < 0 || loc.itemIndex >= len(items) {
			return ErrNotFound
		}
		story := items[loc.itemIndex]
		if loc.subIndex < 0 || loc.subIndex >= len(story.Soal) {
			return ErrNotFound
		}
		story.Soal = append(story.Soal[:loc.subIndex], story.Soal[loc.subIndex+1:]...)
		if len(story.Soal) == 0 {
			items = append(items[:loc.itemIndex], items[loc.itemIndex+1:]...)
		} else {
			items[loc.itemIndex] = story
		}
		return writeJSONFile(path, items)
	}

	return ErrNotFound
}

func (r *Repository) readLessonLocked(number string) ([]Soal, error) {
	result := make([]Soal, 0)

	for _, kategori := range KategoriList {
		dir := filepath.Join(r.path, "lesson"+number, kategori)
		names, err := jsonFiles(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		for _, name := range names {
			path := filepath.Join(dir, name)
			fileBase := strings.TrimSuffix(name, filepath.Ext(name))

			switch kategori {
			case "jtest":
				items, err := readJSONFile[rawJtest](path)
				if err != nil {
					return nil, err
				}
				for i, item := range items {
					soal := jtestToSoal(number, kategori, fileBase, i, item)
					if soal.Soal == "" || soal.Jawaban == "" {
						continue
					}
					result = append(result, soal)
				}

			case "pg", "bintang":
				items, err := readJSONFile[rawPilihan](path)
				if err != nil {
					return nil, err
				}
				for i, item := range items {
					soal := pilihanToSoal(number, kategori, fileBase, i, item)
					if soal.Soal == "" || soal.JawabanBenar == "" || len(soal.PilihanJawaban) == 0 {
						continue
					}
					result = append(result, soal)
				}

			case "dokkai":
				items, err := readJSONFile[rawDokkai](path)
				if err != nil {
					return nil, err
				}
				for i, story := range items {
					for j, sub := range story.Soal {
						soal := dokkaiToSoal(number, kategori, fileBase, i, j, story, sub)
						if soal.Soal == "" || soal.JawabanBenar == "" || len(soal.PilihanJawaban) == 0 {
							continue
						}
						result = append(result, soal)
					}
				}
			}
		}
	}

	return result, nil
}

func (r *Repository) filePath(lesson, kategori, fileBase string) string {
	return filepath.Join(r.path, "lesson"+lesson, kategori, fileBase+".json")
}

func jtestToSoal(lesson, kategori, fileBase string, itemIndex int, raw rawJtest) Soal {
	jawaban := strings.TrimSpace(raw.Jawaban)
	return Soal{
		ID:           flatID(lesson, kategori, fileBase, itemIndex),
		Lesson:       lesson,
		Kategori:     kategori,
		Instruksi:    strings.TrimSpace(raw.Instruksi),
		Soal:         strings.TrimSpace(raw.Soal),
		Jawaban:      jawaban,
		JawabanBenar: jawaban,
	}
}

func pilihanToSoal(lesson, kategori, fileBase string, itemIndex int, raw rawPilihan) Soal {
	return Soal{
		ID:             flatID(lesson, kategori, fileBase, itemIndex),
		Lesson:         lesson,
		Kategori:       kategori,
		Soal:           strings.TrimSpace(raw.Soal),
		PilihanJawaban: trimChoices(raw.PilihanJawaban),
		JawabanBenar:   strings.TrimSpace(raw.JawabanBenar),
	}
}

func dokkaiToSoal(lesson, kategori, fileBase string, itemIndex, subIndex int, story rawDokkai, sub rawDokkaiSoal) Soal {
	return Soal{
		ID:             dokkaiID(lesson, kategori, fileBase, itemIndex, subIndex),
		Lesson:         lesson,
		Kategori:       kategori,
		Soal:           strings.TrimSpace(sub.Pertanyaan),
		PilihanJawaban: trimChoices(sub.PilihanJawaban),
		JawabanBenar:   strings.TrimSpace(sub.JawabanBenar),
		Dokkai:         strings.TrimSpace(story.Dokkai),
		No:             strings.TrimSpace(sub.No),
	}
}

func flatID(lesson, kategori, fileBase string, itemIndex int) string {
	return fmt.Sprintf("l%s_%s_%s_%d", lesson, kategori, fileBase, itemIndex)
}

func dokkaiID(lesson, kategori, fileBase string, itemIndex, subIndex int) string {
	return fmt.Sprintf("l%s_%s_%s_%d_%d", lesson, kategori, fileBase, itemIndex, subIndex)
}

type idLocation struct {
	kategori  string
	fileBase  string
	itemIndex int
	subIndex  int
	hasSub    bool
}

func parseID(lesson, id string) (idLocation, bool) {
	prefix := "l" + lesson + "_"
	if !strings.HasPrefix(id, prefix) {
		return idLocation{}, false
	}
	parts := strings.Split(strings.TrimPrefix(id, prefix), "_")
	if len(parts) < 3 {
		return idLocation{}, false
	}
	kategori, ok := normalizeKategori(parts[0])
	if !ok {
		return idLocation{}, false
	}
	fileBase := strings.TrimSpace(parts[1])
	if fileBase == "" {
		return idLocation{}, false
	}
	itemIndex, err := strconv.Atoi(parts[2])
	if err != nil {
		return idLocation{}, false
	}
	loc := idLocation{kategori: kategori, fileBase: fileBase, itemIndex: itemIndex}
	if len(parts) >= 4 {
		subIndex, err := strconv.Atoi(parts[3])
		if err != nil {
			return idLocation{}, false
		}
		loc.subIndex = subIndex
		loc.hasSub = true
	}
	return loc, true
}

func jsonFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".json") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Slice(names, func(i, j int) bool {
		a, _ := strconv.Atoi(strings.TrimSuffix(names[i], filepath.Ext(names[i])))
		b, _ := strconv.Atoi(strings.TrimSuffix(names[j], filepath.Ext(names[j])))
		if a != 0 && b != 0 && a != b {
			return a < b
		}
		return names[i] < names[j]
	})
	return names, nil
}

func readJSONFile[T any](path string) ([]T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []T{}, nil
		}
		return nil, err
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	if len(strings.TrimSpace(string(data))) == 0 {
		return []T{}, nil
	}
	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("file %s: %w", path, err)
	}
	if items == nil {
		items = []T{}
	}
	return items, nil
}

func writeJSONFile[T any](path string, items []T) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func trimChoices(choices []string) []string {
	result := make([]string, 0, len(choices))
	for _, choice := range choices {
		choice = strings.TrimSpace(choice)
		if choice != "" {
			result = append(result, choice)
		}
	}
	return result
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

func normalizeKategori(value string) (string, bool) {
	text := strings.ToLower(strings.TrimSpace(value))
	if text == "jstest" {
		text = "jtest"
	}
	for _, kategori := range KategoriList {
		if text == kategori {
			return kategori, true
		}
	}
	return "", false
}
