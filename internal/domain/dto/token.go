package dto

// AccessTokenPayload - результат парсинга access токена
// Это DTO (Data Transfer Object) - структура для передачи данных между слоями
// Находится в domain, так как описывает доменные данные и является контрактом между слоями
type AccessTokenPayload struct {
	Phone       string
	Role        string // Строка, не enum.Role - конвертация в enum происходит на уровне Service
	PointCode   string
	NetworkCode string
}

type AuthTokenPayload struct {
	Phone       string `json:"phone"`
	PointCode   string `json:"point_code"`
	NetworkCode string `json:"network_code"`
}
