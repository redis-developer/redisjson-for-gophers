package domain

type MovieRaw struct {
	Title string   `json:"title"`
	Year  int      `json:"year"`
	Info  *InfoRaw `json:"info"`
}

type InfoRaw struct {
	RunningTime float32  `json:"running_time_secs"`
	ReleaseDate string   `json:"release_date"`
	Rating      float32  `json:"rating"`
	Plot        string   `json:"plot"`
	Genres      []string `json:"genres"`
	Actors      []string `json:"actors"`
	Directors   []string `json:"directors"`
}

type Movie struct {
	Title         string    `json:"title"`
	Year          int       `json:"year"`
	Plot          string    `json:"plot"`
	PlotEmbedding []float64 `json:"plotEmbedding"`
	RunningTime   float32   `json:"runningTime"`
	ReleaseDate   string    `json:"releaseDate"`
	Rating        float32   `json:"rating"`
	Genres        []string  `json:"genres"`
	Actors        []string  `json:"actors"`
	Directors     []string  `json:"directors"`
}
