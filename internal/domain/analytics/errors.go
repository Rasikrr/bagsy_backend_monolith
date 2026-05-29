package analytics

import "errors"

var (
	// ErrAccessDenied — у пользователя нет прав на запрашиваемую аналитику.
	ErrAccessDenied = errors.New("analytics: access denied")
	// ErrNotFound — запрашиваемый ресурс (сотрудник, локация) не найден.
	ErrNotFound = errors.New("analytics: not found")
	// ErrInvalidPeriod — невалидные параметры периода (from > to и т.п.).
	ErrInvalidPeriod = errors.New("analytics: invalid period")
)
