package media

//go:generate easyjson -all models.go

type uploadRequest struct {
	Filename  string `json:"filename"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
	Purpose   string `json:"purpose" enums:"avatars,organizations,locations,services,service-categories"`
}

type uploadResponse struct {
	AssetID      string            `json:"asset_id"`
	UploadURL    string            `json:"upload_url"`
	UploadFields map[string]string `json:"upload_fields"`
}
