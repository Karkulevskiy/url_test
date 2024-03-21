package types

// Тип Site, описывает каждый сайт, который был отправлен на сервер
type Site struct {
	ID       string `json:"id"`        // ID каждого сайта
	FullURL  string `json:"full_url"`  // Полный URL сайта, до сокращения
	ShortURL string `json:"short_url"` // Сокращенный вариант URL
}
