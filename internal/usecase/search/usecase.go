package search

import (
	"context"
	"strings"

	domain "server/internal/domain/search"
)

type service struct {
	repo domain.Repository
}

// NewService creates a new search UseCase backed by the given Repository.
func NewService(repo domain.Repository) domain.UseCase {
	return &service{repo: repo}
}

func (s *service) Search(ctx context.Context, filters *domain.SearchFilters) (*domain.SearchResults, error) {
	q := strings.TrimSpace(filters.Query)
	if len(q) < 2 {
		return &domain.SearchResults{
			Query: filters.Query,
			Items: []domain.SearchResultItem{},
		}, nil
	}

	if filters.Limit == 0 {
		filters.Limit = 20
	}
	if filters.Limit > 50 {
		filters.Limit = 50
	}

	return s.repo.Search(ctx, filters)
}