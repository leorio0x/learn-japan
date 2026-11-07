package latihan

// SoalBunpou merepresentasikan satu soal latihan bunpou.
// Format field mengikuti data/format/soal_bunpou.json.
type SoalBunpou struct {
	ID             string   `json:"id,omitempty"`
	Soal           string   `json:"soal"`
	TipeSoal       string   `json:"tipe_soal"`
	PilihanJawaban []string `json:"pilihan_jawaban,omitempty"`
	JawabanBenar   string   `json:"jawaban_benar"`
}

// Request adalah payload untuk membuat atau memperbarui soal latihan.
type Request struct {
	Lesson         string   `json:"lesson"`
	Soal           string   `json:"soal"`
	TipeSoal       string   `json:"tipe_soal"`
	PilihanJawaban []string `json:"pilihan_jawaban"`
	JawabanBenar   string   `json:"jawaban_benar"`
}
