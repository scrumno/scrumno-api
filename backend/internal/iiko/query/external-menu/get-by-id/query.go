package get_by_id

type Query struct {
	ExternalMenuID  string
	OrganizationIDs []string
	PriceCategoryID *string
	Version         *int32
	Language        *string
	StartRevision   *int64
}
