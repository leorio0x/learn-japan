package latihan

// Soal merepresentasikan satu soal latihan bunpou yang sudah dinormalisasi.
// Dokkai dipecah per pertanyaan agar dapat ditampilkan satu per satu.
type Soal struct {
	ID             string   `json:"id,omitempty"`
	Lesson         string   `json:"lesson"`
	Kategori       string   `json:"kategori"`
	Instruksi      string   `json:"instruksi,omitempty"`
	Soal           string   `json:"soal"`
	PilihanJawaban []string `json:"pilihan_jawaban,omitempty"`
	JawabanBenar   string   `json:"jawaban_benar,omitempty"`
	Jawaban        string   `json:"jawaban,omitempty"`
	Dokkai         string   `json:"dokkai,omitempty"`
	No             string   `json:"no,omitempty"`
}

// Request adalah payload untuk membuat atau memperbarui soal latihan.
type Request struct {
	Lesson         string   `json:"lesson"`
	Kategori       string   `json:"kategori"`
	Instruksi      string   `json:"instruksi"`
	Soal           string   `json:"soal"`
	PilihanJawaban []string `json:"pilihan_jawaban"`
	JawabanBenar   string   `json:"jawaban_benar"`
	Jawaban        string   `json:"jawaban"`
	Dokkai         string   `json:"dokkai"`
	No             string   `json:"no"`
}

type rawJtest struct {
	Lesson    string `json:"lesson,omitempty"`
	Kategori  string `json:"kategori,omitempty"`
	Instruksi string `json:"instruksi,omitempty"`
	Soal      string `json:"soal"`
	Jawaban   string `json:"jawaban"`
}

type rawPilihan struct {
	Lesson         string   `json:"lesson,omitempty"`
	Kategori       string   `json:"kategori,omitempty"`
	Soal           string   `json:"soal"`
	PilihanJawaban []string `json:"pilihan_jawaban"`
	JawabanBenar   string   `json:"jawaban_benar"`
}

type rawDokkaiSoal struct {
	No             string   `json:"no,omitempty"`
	Pertanyaan     string   `json:"pertanyaan"`
	PilihanJawaban []string `json:"pilihan_jawaban"`
	JawabanBenar   string   `json:"jawaban_benar"`
}

type rawDokkai struct {
	Lesson   string          `json:"lesson,omitempty"`
	Kategori string          `json:"kategori,omitempty"`
	Dokkai   string          `json:"dokkai"`
	Soal     []rawDokkaiSoal `json:"soal"`
}
