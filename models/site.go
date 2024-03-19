package models

type Site struct {
	ID       string `json:"id"`
	FullURL  string `json:"full_url"`
	ShortURL string `json:"short_url"`
}
