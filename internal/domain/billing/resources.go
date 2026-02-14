package billing

type Resource string

// Quotas (Quantitative limits)

const (
	ResourceMaxLocations Resource = "max_locations"
	ResourceMaxEmployees Resource = "max_employees"
)

// Features (Boolean flags)

const (
	FeatureAnalytics Resource = "analytics"
)
