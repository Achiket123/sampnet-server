package search

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	domain "server/internal/domain/search"

	"github.com/redis/go-redis/v9"
)

type cachedRepository struct {
	client *redis.Client
	inner  domain.Repository
	ttl    time.Duration
}

// NewCachedRepository wraps a search Repository with a Redis caching layer.
func NewCachedRepository(client *redis.Client, inner domain.Repository, ttl time.Duration) domain.Repository {
	return &cachedRepository{
		client: client,
		inner:  inner,
		ttl:    ttl,
	}
}

// Search retrieves results from Redis cache or queries the database on miss.
func (r *cachedRepository) Search(ctx context.Context, filters *domain.SearchFilters) (*domain.SearchResults, error) {
	sortedTypes := make([]string, len(filters.Types))
	copy(sortedTypes, filters.Types)
	sort.Strings(sortedTypes)
	typesStr := strings.Join(sortedTypes, ",")

	queryStr := strings.TrimSpace(strings.ToLower(filters.Query))
	key := fmt.Sprintf("search:%d:%s:%s:%d:%d",
		filters.OrganisationID,
		queryStr,
		typesStr,
		filters.Limit,
		filters.Offset,
	)

	// Attempt to get from cache
	cachedVal, err := r.client.Get(ctx, key).Result()
	if err == nil {
		var results domain.SearchResults
		if err := json.Unmarshal([]byte(cachedVal), &results); err == nil {
			return &results, nil
		}
		log.Printf("Error unmarshaling search cache results: %v", err)
	} else if err != redis.Nil {
		log.Printf("Redis error getting search cache key: %v", err)
	}

	// Database lookup on miss or error
	results, err := r.inner.Search(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Write back to cache
	resultsBytes, err := json.Marshal(results)
	if err == nil {
		err = r.client.Set(ctx, key, resultsBytes, r.ttl).Err()
		if err != nil {
			log.Printf("Redis error setting search cache key: %v", err)
		}
	} else {
		log.Printf("Error marshaling search results for cache: %v", err)
	}

	return results, nil
}
