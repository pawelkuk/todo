package periodictask

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer because the With API mutates the value.
type QueryFilter struct {
	ID          *int64
	Title       *string
	Description *string
	Schedule    *string
}
