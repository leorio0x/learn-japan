package kanji

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// ErrNotFound dikembalikan ketika kanji dengan id tertentu tidak ditemukan.
var ErrNotFound = errors.New("kanji tidak ditemukan")

// Repository menyimpan data kanji pada sebuah file JSON.
type Repository struct {
	path string
	mu   sync.RWMutex
}

// NewRepository membuat repository dengan path file JSON yang diberikan.
func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

func (r *Repository) readLocked() ([]Kanji, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Kanji{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []Kanji{}, nil
	}
	var list []Kanji
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []Kanji{}
	}
	return list, nil
}

func (r *Repository) writeLocked(list []Kanji) error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(r.path, data, 0o644)
}

// FindAll mengembalikan seluruh data kanji.
func (r *Repository) FindAll() ([]Kanji, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.readLocked()
}

// FindByID mengembalikan kanji berdasarkan id.
func (r *Repository) FindByID(id string) (Kanji, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list, err := r.readLocked()
	if err != nil {
		return Kanji{}, err
	}
	for _, k := range list {
		if k.ID == id {
			return k, nil
		}
	}
	return Kanji{}, ErrNotFound
}

// Create menambahkan kanji baru.
func (r *Repository) Create(k Kanji) (Kanji, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.readLocked()
	if err != nil {
		return Kanji{}, err
	}
	list = append(list, k)
	if err := r.writeLocked(list); err != nil {
		return Kanji{}, err
	}
	return k, nil
}

// CreateMany menambahkan beberapa kanji sekaligus dalam satu penulisan.
func (r *Repository) CreateMany(items []Kanji) error {
	if len(items) == 0 {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.readLocked()
	if err != nil {
		return err
	}
	list = append(list, items...)
	return r.writeLocked(list)
}

// Update memperbarui kanji berdasarkan id.
func (r *Repository) Update(id string, k Kanji) (Kanji, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.readLocked()
	if err != nil {
		return Kanji{}, err
	}
	for i := range list {
		if list[i].ID == id {
			k.ID = id
			list[i] = k
			if err := r.writeLocked(list); err != nil {
				return Kanji{}, err
			}
			return k, nil
		}
	}
	return Kanji{}, ErrNotFound
}

// Delete menghapus kanji berdasarkan id.
func (r *Repository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.readLocked()
	if err != nil {
		return err
	}
	for i := range list {
		if list[i].ID == id {
			list = append(list[:i], list[i+1:]...)
			return r.writeLocked(list)
		}
	}
	return ErrNotFound
}

// DeleteMany menghapus beberapa kanji sekaligus berdasarkan daftar id.
// Mengembalikan jumlah data yang benar-benar terhapus.
func (r *Repository) DeleteMany(ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	want := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		want[id] = struct{}{}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.readLocked()
	if err != nil {
		return 0, err
	}
	kept := make([]Kanji, 0, len(list))
	removed := 0
	for _, k := range list {
		if _, ok := want[k.ID]; ok {
			removed++
			continue
		}
		kept = append(kept, k)
	}
	if removed == 0 {
		return 0, nil
	}
	if err := r.writeLocked(kept); err != nil {
		return 0, err
	}
	return removed, nil
}
