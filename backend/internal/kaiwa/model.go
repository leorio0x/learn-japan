package kaiwa

// Dialog merepresentasikan satu tuturan kaiwa (pertanyaan atau jawaban).
// Bentuk kanji berisi kalimat lengkap, sedangkan hiragana/katakana opsional.
type Dialog struct {
	Kanji      string `json:"kanji"`
	Kana       string `json:"hiragana/katakana"`
	Terjemahan string `json:"terjemahan"`
}

// Contoh merepresentasikan satu contoh percakapan yang memuat daftar
// tuturan pertanyaan dan daftar tuturan jawaban.
type Contoh struct {
	Pertanyaan []Dialog `json:"pertanyaan"`
	Jawaban    []Dialog `json:"jawaban"`
}

// Pola merepresentasikan pola kalimat beserta contoh-contohnya.
type Pola struct {
	PolaPertanyaan string   `json:"pola_pertanyaan"`
	PolaJawaban    string   `json:"pola_jawaban"`
	Deskripsi      string   `json:"deskripsi"`
	Contoh         []Contoh `json:"contoh"`
}

// Kaiwa merepresentasikan satu entri percakapan.
type Kaiwa struct {
	ID     string `json:"id"`
	Pola   []Pola `json:"pola"`
	Lesson string `json:"lesson"`
	Buku   string `json:"buku"`
}

// Request adalah payload untuk membuat / memperbarui kaiwa.
type Request struct {
	Pola   []Pola `json:"pola"`
	Lesson string `json:"lesson"`
	Buku   string `json:"buku"`
}

// normalize memastikan seluruh slice tidak nil agar JSON tersimpan sebagai [].
func (c *Contoh) normalize() {
	if c.Pertanyaan == nil {
		c.Pertanyaan = []Dialog{}
	}
	if c.Jawaban == nil {
		c.Jawaban = []Dialog{}
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

func (k *Kaiwa) normalize() {
	if k.Pola == nil {
		k.Pola = []Pola{}
	}
	for i := range k.Pola {
		k.Pola[i].normalize()
	}
}
