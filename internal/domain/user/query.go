package user

import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"

type Filter struct {
	NetworkCode *string
	PointCode   *string
	Roles       []Role
	PhoneSearch *string // Частичный или полный поиск по номеру телефона (ILIKE %value%)
	Limit       uint64
	Offset      uint64
	OrderBy     string
	SortOrder   enum.SortOrder
}

// CustomerFilter - фильтр для получения клиентов (users с role='user'), обслуживавшихся в точках
// Поля для авторизации устанавливаются на уровне сервиса в зависимости от роли пользователя
type CustomerFilter struct {
	PhoneSearch *string // Частичный или полный поиск по номеру телефона клиента (ILIKE %value%)
	// Поля авторизации (устанавливаются сервисом):
	MasterPhone *string  // Для staff - устанавливается автоматически. Для manager+ - из запроса после валидации
	PointCodes  []string // Для manager/net_manager/self_owner - клиенты из этих точек
	Limit       uint64
	Offset      uint64
	OrderBy     string
	SortOrder   enum.SortOrder
}
