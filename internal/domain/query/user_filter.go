package query

type UserFilter struct {
	NetworkCode *string
	PointCode   *string
	Roles       []string
	Phones      []string
}
