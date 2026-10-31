package kaiwa

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// ErrNotFound dikembalikan ketika kaiwa dengan id tertentu tidak ditemukan.
var ErrNotFound = errors.New("kaiwa tidak ditemukan")

// Repository menyimpan data kaiwa pada sebuah file JSON.
type Repository struct {
	path string
	mu   sync.RWMutex
}

// NewRepository membuat repository dengan path file JSON yang diberikan.
func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

func (r *Repository) readLocked() ([]Kaiwa, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Kaiwa{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []Kaiwa{}, nil
	}
	var list []Kaiwa
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []Kaiwa{}
	}
	for i := range list {
		list[i].normalize()
	}
	return list, nil
}

func (r *Repository) writeLocked(list []Kaiwa) error {
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

// FindAll mengembalikan seluruh data kaiwa.
func (r *Repository) FindAll() ([]Kaiwa, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.readLocked()
}

// FindByID mengembalikan kaiwa berdasarkan id.
func (r *Repository) FindByID(id string) (Kaiwa, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list, err := r.readLocked()
	if err != nil {
		return Kaiwa{}, err
	}
	for _, k := range list {
		if k.ID == id {
			return k, nil
		}
	}
	return Kaiwa{}, ErrNotFound
}

// Create menambahkan kaiwa baru.
func (r *Repository) Create(k Kaiwa) (Kaiwa, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.readLocked()
	if err != nil {
		return Kaiwa{}, err
	}
	list = append(list, k)
	if err := r.writeLocked(list); err != nil {
		return Kaiwa{}, err
	}
	return k, nil
}

// Update memperbarui kaiwa berdasarkan id.
func (r *Repository) Update(id string, k Kaiwa) (Kaiwa, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.readLocked()
	if err != nil {
		return Kaiwa{}, err
	}
	for i := range list {
		if list[i].ID == id {
			k.ID = id
			list[i] = k
			if err := r.writeLocked(list); err != nil {
				return Kaiwa{}, err
			}
			return k, nil
		}
	}
	return Kaiwa{}, ErrNotFound
}

// Delete menghapus kaiwa berdasarkan id.
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

// DeleteMany menghapus beberapa kaiwa sekaligus berdasarkan daftar id.
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
	kept := make([]Kaiwa, 0, len(list))
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
