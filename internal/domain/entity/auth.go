package entity

type Auth struct {
	AccessToken  string
	RefreshToken string
}

type PayloadParams struct {
	Phone  string
	Role   string
	Active bool
}
