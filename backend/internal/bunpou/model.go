package bunpou

// Kalimat merepresentasikan satu kalimat bunpou beserta terjemahannya.
// Field kalimat dapat memuat format <kanji:kana:arti>.
type Kalimat struct {
	Kalimat    string `json:"kalimat"`
	Terjemahan string `json:"terjemahan"`
}

// Contoh merepresentasikan satu contoh pemakaian pola bunpou.
// Pertanyaan bersifat opsional (boleh null) untuk contoh berupa pernyataan.
type Contoh struct {
	Pertanyaan *Kalimat `json:"pertanyaan"`
	Jawaban    *Kalimat `json:"jawaban"`
}

// Pola merepresentasikan satu pola bunpou beserta contoh-contohnya.
type Pola struct {
	NoBunpou  string   `json:"no_bunpou"`
	NamaPola  string   `json:"nama_pola"`
	Rumus     string   `json:"rumus"`
	Deskripsi string   `json:"deskripsi"`
	Contoh    []Contoh `json:"contoh"`
}

// Bunpou merepresentasikan satu entri materi bunpou per lesson.
type Bunpou struct {
	ID     string `json:"id,omitempty"`
	Lesson string `json:"lesson"`
	Buku   string `json:"buku"`
	Pola   []Pola `json:"pola"`
}

// Request adalah payload untuk membuat / memperbarui bunpou.
type Request struct {
	Lesson string `json:"lesson"`
	Buku   string `json:"buku"`
	Pola   []Pola `json:"pola"`
}

func (c *Contoh) normalize() {
	if c.Jawaban == nil {
		c.Jawaban = &Kalimat{}
	}
}

func (p *Pola) normalize() {
	if p.Contoh == nil {
		p.Contoh = []Contoh{}
	}
	for i := range p.Contoh {
		p.Contoh[i].normalize()
	}
}

func (b *Bunpou) normalize() {
	if b.Pola == nil {
		b.Pola = []Pola{}
	}
	for i := range b.Pola {
		b.Pola[i].normalize()
	}
}
