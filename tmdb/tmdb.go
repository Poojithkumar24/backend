package tmdb

type TMDbMovie struct {
	Title        string `json:"title"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
}

type TMDbResponse struct {
	Results []TMDbMovie `json:"results"`
}
