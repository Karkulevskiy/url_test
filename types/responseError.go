package types

// Структура для описания ошибки во время рантайма
type ResponseError struct {
	Message string `json:"message"` // Сообщение об ошибке
	Status  int    `json:"status"`  // Код ошибки
}
