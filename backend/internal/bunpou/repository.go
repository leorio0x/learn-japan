package bunpou

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// ErrNotFound dikembalikan ketika bunpou dengan id tertentu tidak ditemukan.
var ErrNotFound = errors.New("bunpou tidak ditemukan")

// Repository menyimpan data bunpou pada sebuah file JSON.
type Repository struct {
	path string
	mu   sync.RWMutex
}

// NewRepository membuat repository dengan path file JSON yang diberikan.
func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

func (r *Repository) readLocked() ([]Bunpou, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Bunpou{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []Bunpou{}, nil
	}
	var list []Bunpou
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []Bunpou{}
	}
	missingID := false
	for i := range list {
		list[i].normalize()
		if list[i].ID == "" {
			list[i].ID = newID()
			missingID = true
		}
	}
	// Data yang ditulis manual mengikuti format tidak memiliki id.
	// Lengkapi lalu simpan agar CRUD tetap dapat mengubah dan menghapus data.
	if missingID {
		if err := r.writeLocked(list); err != nil {
			return nil, err
		}
	}
	return list, nil
}

func (r *Repository) writeLocked(list []Bunpou) error {
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

// FindAll mengembalikan seluruh data bunpou.
func (r *Repository) FindAll() ([]Bunpou, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.readLocked()
}

// FindByID mengembalikan bunpou berdasarkan id.
func (r *Repository) FindByID(id string) (Bunpou, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.readLocked()
	if err != nil {
		return Bunpou{}, err
	}
	for _, b := range list {
		if b.ID == id {
			return b, nil
		}
	}
	return Bunpou{}, ErrNotFound
}

// Create menambahkan bunpou baru.
func (r *Repository) Create(b Bunpou) (Bunpou, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.readLocked()
	if err != nil {
		return Bunpou{}, err
	}
	list = append(list, b)
	if err := r.writeLocked(list); err != nil {
		return Bunpou{}, err
	}
	return b, nil
}

// Update memperbarui bunpou berdasarkan id.
func (r *Repository) Update(id string, b Bunpou) (Bunpou, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list, err := r.readLocked()
	if err != nil {
		return Bunpou{}, err
	}
	for i := range list {
		if list[i].ID == id {
			b.ID = id
			list[i] = b
			if err := r.writeLocked(list); err != nil {
				return Bunpou{}, err
			}
			return b, nil
		}
	}
	return Bunpou{}, ErrNotFound
}

// Delete menghapus bunpou berdasarkan id.
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

// DeleteMany menghapus beberapa bunpou sekaligus berdasarkan daftar id.
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
	kept := make([]Bunpou, 0, len(list))
	removed := 0
	for _, b := range list {
		if _, ok := want[b.ID]; ok {
			removed++
			continue
		}
		kept = append(kept, b)
	}
	if removed == 0 {
		return 0, nil
	}
	if err := r.writeLocked(kept); err != nil {
		return 0, err
	}
	return removed, nil
}
