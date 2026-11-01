package dokkai

// LevelDokkai merepresentasikan satu dokkai berlevel (N5, N4, dst).
// Format field mengikuti data/format/dokkai-level.json.
type LevelDokkai struct {
	ID         string `json:"id"`
	Level      string `json:"level"`
	Judul      string `json:"judul"`
	Paragraf   string `json:"paragraf"`
	Terjemahan string `json:"terjemahan"`
}
