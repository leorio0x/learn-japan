package kotoba

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// ErrNotFound dikembalikan ketika kotoba dengan id tertentu tidak ditemukan.
var ErrNotFound = errors.New("kotoba tidak ditemukan")

// Repository menyimpan data kotoba pada sebuah file JSON.
type Repository struct {
	path string
	mu   sync.RWMutex
}

// NewRepository membuat repository dengan path file JSON yang diberikan.
func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

func (r *Repository) readLocked() ([]Kotoba, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Kotoba{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []Kotoba{}, nil
	}
	var list []Kotoba
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []Kotoba{}
	}
	return list, nil
}

func (r *Repository) writeLocked(list []Kotoba) error {
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

// FindAll mengembalikan seluruh kotoba.
func (r *Repository) FindAll() ([]Kotoba, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.readLocked()
}

// FindByID mengembalikan kotoba berdasarkan id.
func (r *Repository) FindByID(id string) (Kotoba, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list, err := r.readLocked()
	if err != nil {
		return Kotoba{}, err
	}
	for _, k := range list {
		if k.ID == id {
			return k, nil
		}
	}
	return Kotoba{}, ErrNotFound
}

// Create menambahkan kotoba baru.
func (r *Repository) Create(k Kotoba) (Kotoba, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.readLocked()
	if err != nil {
		return Kotoba{}, err
	}
	list = append(list, k)
	if err := r.writeLocked(list); err != nil {
		return Kotoba{}, err
	}
	return k, nil
}

// CreateMany menambahkan beberapa kotoba sekaligus dalam satu penulisan.
func (r *Repository) CreateMany(items []Kotoba) error {
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

// Update memperbarui kotoba berdasarkan id.
func (r *Repository) Update(id string, k Kotoba) (Kotoba, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.readLocked()
	if err != nil {
		return Kotoba{}, err
	}
	for i := range list {
		if list[i].ID == id {
			k.ID = id
			list[i] = k
			if err := r.writeLocked(list); err != nil {
				return Kotoba{}, err
			}
			return k, nil
		}
	}
	return Kotoba{}, ErrNotFound
}

// Delete menghapus kotoba berdasarkan id.
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

// DeleteMany menghapus beberapa kotoba sekaligus berdasarkan daftar id.
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
	kept := make([]Kotoba, 0, len(list))
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
