package enum

//go:generate enumer -type=PointMediaType -json -trimprefix PointMedia -transform=snake_case

// PointMediaType представляет тип фото локации точки
type PointMediaType int8

const (
	// PointMediaTypeExterior - фото фасада/экстерьера
	PointMediaTypeExterior PointMediaType = iota
	// PointMediaTypeInterior - фото интерьера
	PointMediaTypeInterior
	// PointMediaTypeMap - карта проезда
	PointMediaTypeMap
	// PointMediaTypeMenu - фото меню/прайс-листа
	PointMediaTypeMenu
)
