package search

import (
	"context"
	"errors"
	"server/internal/domain/search"
	"testing"
)

// mockRepository is a test double for search.Repository.
type mockRepository struct {
	called  bool
	results *search.SearchResults
	err     error
}

func (m *mockRepository) Search(_ context.Context, _ *search.SearchFilters) (*search.SearchResults, error) {
	m.called = true
	return m.results, m.err
}

func TestSearch_EmptyQuery_DoesNotCallRepo(t *testing.T) {
	repo := &mockRepository{}
	svc := NewService(repo)

	res, err := svc.Search(context.Background(), &search.SearchFilters{Query: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.called {
		t.Error("expected repository NOT to be called for empty query")
	}
	if res.TotalCount != 0 {
		t.Errorf("expected 0 results, got %d", res.TotalCount)
	}
}

func TestSearch_ShortQuery_DoesNotCallRepo(t *testing.T) {
	repo := &mockRepository{}
	svc := NewService(repo)

	res, err := svc.Search(context.Background(), &search.SearchFilters{Query: "a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.called {
		t.Error("expected repository NOT to be called for single-character query")
	}
	if len(res.Items) != 0 {
		t.Errorf("expected empty items, got %d", len(res.Items))
	}
}

func TestSearch_ValidQuery_DelegatesToRepo(t *testing.T) {
	expected := &search.SearchResults{
		Query:      "flutter",
		TotalCount: 2,
		Items: []search.SearchResultItem{
			{ID: 1, Type: "task", Title: "Flutter task"},
			{ID: 2, Type: "project", Title: "Flutter project"},
		},
	}
	repo := &mockRepository{results: expected}
	svc := NewService(repo)

	res, err := svc.Search(context.Background(), &search.SearchFilters{
		Query:          "flutter",
		OrganisationID: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !repo.called {
		t.Error("expected repository to be called for valid query")
	}
	if res.TotalCount != expected.TotalCount {
		t.Errorf("expected TotalCount %d, got %d", expected.TotalCount, res.TotalCount)
	}
}

func TestSearch_RepoError_PropagatesError(t *testing.T) {
	repoErr := errors.New("db connection failed")
	repo := &mockRepository{err: repoErr}
	svc := NewService(repo)

	_, err := svc.Search(context.Background(), &search.SearchFilters{
		Query:          "flutter",
		OrganisationID: 1,
	})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, repoErr) {
		t.Errorf("expected %v, got %v", repoErr, err)
	}
}

func TestSearch_DefaultLimit(t *testing.T) {
	var capturedFilters *search.SearchFilters
	repo := &captureRepo{captureFunc: func(f *search.SearchFilters) {
		capturedFilters = f
	}}
	svc := NewService(repo)

	_, _ = svc.Search(context.Background(), &search.SearchFilters{
		Query:          "go",
		OrganisationID: 1,
		Limit:          0, // zero → should default to 20
	})

	if capturedFilters == nil {
		t.Fatal("repository was not called")
	}
	if capturedFilters.Limit != 20 {
		t.Errorf("expected default limit 20, got %d", capturedFilters.Limit)
	}
}

func TestSearch_MaxLimit(t *testing.T) {
	var capturedFilters *search.SearchFilters
	repo := &captureRepo{captureFunc: func(f *search.SearchFilters) {
		capturedFilters = f
	}}
	svc := NewService(repo)

	_, _ = svc.Search(context.Background(), &search.SearchFilters{
		Query:          "go",
		OrganisationID: 1,
		Limit:          200, // should be capped to 50
	})

	if capturedFilters == nil {
		t.Fatal("repository was not called")
	}
	if capturedFilters.Limit != 50 {
		t.Errorf("expected max limit 50, got %d", capturedFilters.Limit)
	}
}

// captureRepo captures the filters it receives so tests can inspect them.
type captureRepo struct {
	captureFunc func(f *search.SearchFilters)
}

func (c *captureRepo) Search(ctx context.Context, f *search.SearchFilters) (*search.SearchResults, error) {
	c.captureFunc(f)
	return &search.SearchResults{Items: []search.SearchResultItem{}}, nil
}
