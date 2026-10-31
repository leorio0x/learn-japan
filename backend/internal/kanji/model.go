package kanji

// Kanji merepresentasikan satu entri kanji.
// Satu kanji dapat memiliki satu atau lebih arti (terjemahan).
type Kanji struct {
	ID               string       `json:"id"`
	Kanji            string       `json:"kanji"`
	Furigana         string       `json:"furigana"`
	Terjemahan       []string     `json:"terjemahan"`
	ContohPercakapan []Percakapan `json:"contoh_percakapan"`
	Buku             string       `json:"buku"`
	Lesson           string       `json:"lesson"`
}

// Percakapan merepresentasikan satu contoh percakapan pada kanji.
type Percakapan struct {
	Pertanyaan     string `json:"pertanyaan"`
	ArtiPertanyaan string `json:"arti_pertanyaan"`
	Jawaban        string `json:"jawaban"`
	ArtiJawaban    string `json:"arti_jawaban"`
}

// Request adalah payload untuk membuat / memperbarui kanji.
type Request struct {
	Kanji            string       `json:"kanji"`
	Furigana         string       `json:"furigana"`
	Terjemahan       []string     `json:"terjemahan"`
	ContohPercakapan []Percakapan `json:"contoh_percakapan"`
	Buku             string       `json:"buku"`
	Lesson           string       `json:"lesson"`
}
