package search

import "context"

// Repository defines the data access contract for search.
type Repository interface {
	Search(ctx context.Context, filters *SearchFilters) (*SearchResults, error)
}