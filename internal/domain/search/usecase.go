package search

import "context"

// UseCase defines the business logic contract for search.
type UseCase interface {
	Search(ctx context.Context, filters *SearchFilters) (*SearchResults, error)
}
