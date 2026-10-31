// Package bulk menyediakan parser untuk fitur bulk input kotoba & kanji
// sesuai aturan di data/rules.md. Mendukung file .csv, .xlsx, dan .json.
package bulk

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

// maxUploadSize membatasi ukuran file yang diunggah (5 MB).
const maxUploadSize = 5 << 20

// Pemisah multi terjemahan: "|" untuk csv dan "," untuk excel (data/rules.md).
const (
	csvDelimiter   = "|"
	excelDelimiter = ","
)

// ErrNoFile dikembalikan ketika request tidak menyertakan file.
var ErrNoFile = errors.New(`file wajib diunggah pada field "file"`)

// Row merepresentasikan satu baris data hasil parsing bulk input.
type Row struct {
	Line             int
	Kanji            string
	Furigana         string
	Terjemahan       []string
	ContohPercakapan []Percakapan
	Buku             string
	Lesson           string
}

// Percakapan merepresentasikan satu contoh percakapan pada bulk input JSON.
type Percakapan struct {
	Pertanyaan     string `json:"pertanyaan"`
	ArtiPertanyaan string `json:"arti_pertanyaan"`
	Jawaban        string `json:"jawaban"`
	ArtiJawaban    string `json:"arti_jawaban"`
}

// jsonRow adalah bentuk data pada file JSON bulk input.
// Field id sengaja tidak didefinisikan agar selalu digenerate server.
type jsonRow struct {
	Kanji            string       `json:"kanji"`
	Furigana         string       `json:"furigana"`
	Terjemahan       []string     `json:"terjemahan"`
	ContohPercakapan []Percakapan `json:"contoh_percakapan"`
	Buku             string       `json:"buku"`
	Lesson           string       `json:"lesson"`
}

// ValidationError menandakan baris yang tidak lolos validasi.
type ValidationError struct {
	Line int
	Msg  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("baris %d: %s", e.Line, e.Msg)
}

// ParseUpload membaca file dari multipart form fiber lalu mem-parsing-nya.
func ParseUpload(c *fiber.Ctx, field string) ([]Row, error) {
	file, err := c.FormFile(field)
	if err != nil {
		return nil, ErrNoFile
	}
	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(io.LimitReader(f, maxUploadSize+1))
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file: %w", err)
	}
	if len(data) > maxUploadSize {
		return nil, fmt.Errorf("ukuran file melebihi %d MB", maxUploadSize>>20)
	}
	return Parse(file.Filename, data)
}

// Parse mem-parsing data csv/xlsx/json menjadi daftar Row.
func Parse(filename string, data []byte) ([]Row, error) {
	switch ext := strings.ToLower(filepath.Ext(filename)); ext {
	case ".csv":
		return parseCSV(data)
	case ".xlsx":
		return parseXLSX(data)
	case ".json":
		return parseJSON(data)
	case "":
		return nil, errors.New("nama file tidak memiliki ekstensi; gunakan .csv, .xlsx, atau .json")
	default:
		return nil, fmt.Errorf("format file %s tidak didukung (gunakan .csv, .xlsx, atau .json)", ext)
	}
}

func parseCSV(data []byte) ([]Row, error) {
	reader := csv.NewReader(bytes.NewReader(stripBOM(data)))
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("gagal membaca csv: %w", err)
	}
	if len(records) == 0 {
		return nil, errors.New("file kosong")
	}
	idx, err := headerIndex(records[0])
	if err != nil {
		return nil, err
	}
	return rowsFromRecords(records[1:], idx, csvDelimiter), nil
}

func parseXLSX(data []byte) ([]Row, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("gagal membuka xlsx: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, errors.New("file excel tidak memiliki sheet")
	}
	records, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet: %w", err)
	}
	if len(records) == 0 {
		return nil, errors.New("file kosong")
	}
	idx, err := headerIndex(records[0])
	if err != nil {
		return nil, err
	}
	return rowsFromRecords(records[1:], idx, excelDelimiter), nil
}

func parseJSON(data []byte) ([]Row, error) {
	var items []jsonRow
	if err := json.Unmarshal(stripBOM(data), &items); err != nil {
		return nil, fmt.Errorf("gagal membaca json: %w", err)
	}
	if len(items) == 0 {
		return nil, errors.New("file json kosong")
	}

	rows := make([]Row, 0, len(items))
	for i, item := range items {
		rows = append(rows, Row{
			Line:             i + 1,
			Kanji:            strings.TrimSpace(item.Kanji),
			Furigana:         strings.TrimSpace(item.Furigana),
			Terjemahan:       cleanValues(item.Terjemahan),
			ContohPercakapan: cleanPercakapan(item.ContohPercakapan),
			Buku:             strings.TrimSpace(item.Buku),
			Lesson:           strings.TrimSpace(item.Lesson),
		})
	}
	return rows, nil
}

func cleanValues(items []string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		if v := strings.TrimSpace(item); v != "" {
			out = append(out, v)
		}
	}
	return out
}

func cleanPercakapan(items []Percakapan) []Percakapan {
	out := make([]Percakapan, 0, len(items))
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
		out = append(out, p)
	}
	return out
}

func headerIndex(header []string) (map[string]int, error) {
	idx := make(map[string]int, 5)
	for i, h := range header {
		switch normalizeHeader(h) {
		case "kanji":
			idx["kanji"] = i
		case "furigana", "hiragana/katakana", "hiragana / katakana", "kana", "hiragana", "katakana":
			idx["furigana"] = i
		case "terjemahan", "arti":
			idx["terjemahan"] = i
		case "buku":
			idx["buku"] = i
		case "lesson":
			idx["lesson"] = i
		}
	}
	if _, ok := idx["kanji"]; !ok {
		return nil, errors.New("header 'kanji' tidak ditemukan")
	}
	if _, ok := idx["furigana"]; !ok {
		return nil, errors.New("header 'furigana' (hiragana/katakana) tidak ditemukan")
	}
	if _, ok := idx["terjemahan"]; !ok {
		return nil, errors.New("header 'terjemahan' tidak ditemukan")
	}
	return idx, nil
}

func rowsFromRecords(records [][]string, idx map[string]int, delimiter string) []Row {
	rows := make([]Row, 0, len(records))
	for i, rec := range records {
		kanji := cell(rec, idx, "kanji")
		furigana := cell(rec, idx, "furigana")
		terjemahan := cell(rec, idx, "terjemahan")
		buku := cell(rec, idx, "buku")
		lesson := cell(rec, idx, "lesson")

		if kanji == "" && furigana == "" && terjemahan == "" && buku == "" && lesson == "" {
			continue
		}

		rows = append(rows, Row{
			Line:       i + 2, // 1 header + index 1-based
			Kanji:      kanji,
			Furigana:   furigana,
			Terjemahan: splitTerjemahan(terjemahan, delimiter),
			Buku:       buku,
			Lesson:     lesson,
		})
	}
	return rows
}

func cell(rec []string, idx map[string]int, key string) string {
	i, ok := idx[key]
	if !ok || i >= len(rec) {
		return ""
	}
	return strings.TrimSpace(rec[i])
}

func splitTerjemahan(raw, delimiter string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, delimiter)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}

func normalizeHeader(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	return strings.TrimPrefix(s, "\ufeff")
}

func stripBOM(data []byte) []byte {
	return bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
}
