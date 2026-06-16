package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	domain "server/internal/domain/analytics"

	"github.com/redis/go-redis/v9"
)

type cachedRepository struct {
	client *redis.Client
	inner  domain.Repository
	ttl    time.Duration
}

// NewCachedRepository wraps an analytics Repository with a Redis caching layer.
func NewCachedRepository(client *redis.Client, inner domain.Repository, ttl time.Duration) domain.Repository {
	return &cachedRepository{
		client: client,
		inner:  inner,
		ttl:    ttl,
	}
}

func (r *cachedRepository) GetEmployeeAnalytics(ctx context.Context, userID uint, orgID uint, period string) (*domain.EmployeeAnalyticsSummary, error) {
	key := fmt.Sprintf("analytics:employee:%d:%d:%s", userID, orgID, period)

	cachedVal, err := r.client.Get(ctx, key).Result()
	if err == nil {
		var summary domain.EmployeeAnalyticsSummary
		if err := json.Unmarshal([]byte(cachedVal), &summary); err == nil {
			return &summary, nil
		}
		log.Printf("Error unmarshaling cached employee analytics: %v", err)
	} else if err != redis.Nil {
		log.Printf("Redis error getting cached employee analytics: %v", err)
	}

	summary, err := r.inner.GetEmployeeAnalytics(ctx, userID, orgID, period)
	if err != nil {
		return nil, err
	}

	summaryBytes, err := json.Marshal(summary)
	if err == nil {
		if err := r.client.Set(ctx, key, summaryBytes, r.ttl).Err(); err != nil {
			log.Printf("Redis error setting cached employee analytics: %v", err)
		}
	} else {
		log.Printf("Error marshaling employee analytics for cache: %v", err)
	}

	return summary, nil
}

func (r *cachedRepository) GetOrgEmployeeMonitor(ctx context.Context, orgID uint) ([]domain.EmployeeMonitorResponse, error) {
	key := fmt.Sprintf("analytics:org_monitor:%d", orgID)

	cachedVal, err := r.client.Get(ctx, key).Result()
	if err == nil {
		var monitor []domain.EmployeeMonitorResponse
		if err := json.Unmarshal([]byte(cachedVal), &monitor); err == nil {
			return monitor, nil
		}
		log.Printf("Error unmarshaling cached org employee monitor: %v", err)
	} else if err != redis.Nil {
		log.Printf("Redis error getting cached org employee monitor: %v", err)
	}

	monitor, err := r.inner.GetOrgEmployeeMonitor(ctx, orgID)
	if err != nil {
		return nil, err
	}

	monitorBytes, err := json.Marshal(monitor)
	if err == nil {
		if err := r.client.Set(ctx, key, monitorBytes, r.ttl).Err(); err != nil {
			log.Printf("Redis error setting cached org employee monitor: %v", err)
		}
	} else {
		log.Printf("Error marshaling org employee monitor for cache: %v", err)
	}

	return monitor, nil
}
