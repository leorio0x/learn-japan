package dokkai

// Dokkai merepresentasikan satu cerita pendek (dokkai).
// Format field mengikuti data/format/dokkai.json.
type Dokkai struct {
	Lesson     string `json:"lesson"`
	Judul      string `json:"judul"`
	Paragraf   string `json:"paragraf"`
	Terjemahan string `json:"terjemahan"`
}
