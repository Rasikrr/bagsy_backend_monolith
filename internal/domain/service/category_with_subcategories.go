package service

type CategoryWithSubcategories struct {
	Category      *Category
	Subcategories []*Subcategory
}
