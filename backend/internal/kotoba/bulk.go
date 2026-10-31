package kotoba

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"learn-japan/backend/internal/bulk"
)

var seqPattern = regexp.MustCompile(`^ktb-(\d+)$`)

// Bulk memvalidasi seluruh baris lalu menyimpannya dengan id yang dilanjutkan
// dari data existing (data/rules.md).
func (s *Service) Bulk(rows []bulk.Row) ([]Kotoba, error) {
	if len(rows) == 0 {
		return nil, errors.New("tidak ada data yang bisa diimpor")
	}

	requests := make([]Request, len(rows))
	for i, r := range rows {
		requests[i] = Request{
			Kanji:            r.Kanji,
			Furigana:         r.Furigana,
			Terjemahan:       r.Terjemahan,
			ContohPercakapan: toPercakapanList(r.ContohPercakapan),
			Buku:             r.Buku,
			Lesson:           r.Lesson,
		}
		if msg := validate(requests[i]); msg != "" {
			return nil, &bulk.ValidationError{Line: r.Line, Msg: msg}
		}
	}

	existing, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	next := nextSequence(existing)
	used := make(map[string]struct{}, len(existing))
	for _, k := range existing {
		used[k.ID] = struct{}{}
	}

	created := make([]Kotoba, 0, len(rows))
	for _, req := range requests {
		k := toKotoba(req)
		for {
			next++
			candidate := fmt.Sprintf("ktb-%03d", next)
			if _, taken := used[candidate]; !taken {
				k.ID = candidate
				used[candidate] = struct{}{}
				break
			}
		}
		created = append(created, k)
	}

	if err := s.repo.CreateMany(created); err != nil {
		return nil, err
	}
	return created, nil
}

func toPercakapanList(items []bulk.Percakapan) []Percakapan {
	out := make([]Percakapan, 0, len(items))
	for _, item := range items {
		out = append(out, Percakapan{
			Pertanyaan:     item.Pertanyaan,
			ArtiPertanyaan: item.ArtiPertanyaan,
			Jawaban:        item.Jawaban,
			ArtiJawaban:    item.ArtiJawaban,
		})
	}
	return out
}

func nextSequence(list []Kotoba) int {
	max := 0
	for _, k := range list {
		if m := seqPattern.FindStringSubmatch(k.ID); m != nil {
			if n, err := strconv.Atoi(m[1]); err == nil && n > max {
				max = n
			}
		}
	}
	return max
}
