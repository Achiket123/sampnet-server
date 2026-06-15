package auth

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	domain "server/internal/domain/auth"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	client   *redis.Client
	gormRepo domain.Repository
}

func NewRedisRepository(client *redis.Client, gormRepo domain.Repository) domain.Repository {
	return &redisRepository{
		client:   client,
		gormRepo: gormRepo,
	}
}

// Delegate non-token methods directly to GORM repository
func (r *redisRepository) Create(ctx context.Context, user *domain.User) error {
	return r.gormRepo.Create(ctx, user)
}

func (r *redisRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.gormRepo.GetByEmail(ctx, email)
}

func (r *redisRepository) GetByPhoneNumber(ctx context.Context, phone string) (*domain.User, error) {
	return r.gormRepo.GetByPhoneNumber(ctx, phone)
}

func (r *redisRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	return r.gormRepo.GetByID(ctx, id)
}

func (r *redisRepository) Update(ctx context.Context, user *domain.User) error {
	return r.gormRepo.Update(ctx, user)
}

func (r *redisRepository) CreateEmailVerification(ctx context.Context, ev *domain.EmailVerification) error {
	return r.gormRepo.CreateEmailVerification(ctx, ev)
}

func (r *redisRepository) GetEmailVerificationByToken(ctx context.Context, token string) (*domain.EmailVerification, error) {
	return r.gormRepo.GetEmailVerificationByToken(ctx, token)
}

func (r *redisRepository) GetActiveEmailVerificationByUserID(ctx context.Context, userID uint) (*domain.EmailVerification, error) {
	return r.gormRepo.GetActiveEmailVerificationByUserID(ctx, userID)
}

func (r *redisRepository) UpdateEmailVerification(ctx context.Context, ev *domain.EmailVerification) error {
	return r.gormRepo.UpdateEmailVerification(ctx, ev)
}

// Token methods implemented using Redis
func (r *redisRepository) SaveRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	key := "refresh:" + token.TokenHash
	userIDStr := strconv.FormatUint(uint64(token.UserID), 10)
	ttl := time.Until(time.Unix(token.ExpiresAt, 0))

	if ttl <= 0 {
		return errors.New("token has already expired")
	}

	// Set the refresh token key mapping to the userID string
	if err := r.client.Set(ctx, key, userIDStr, ttl).Err(); err != nil {
		return err
	}

	// SAdd the tokenHash to the user's tracking set so we can revoke all later if needed
	userSetKey := "user_tokens:" + userIDStr
	if err := r.client.SAdd(ctx, userSetKey, token.TokenHash).Err(); err != nil {
		return err
	}

	// Expire the user's tracking set after 30 days of inactivity to prevent memory leakage
	_ = r.client.Expire(ctx, userSetKey, 30*24*time.Hour)

	return nil
}

func (r *redisRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	key := "refresh:" + tokenHash
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New("invalid or expired refresh token")
	} else if err != nil {
		return nil, err
	}

	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// TTL of -2 means the key doesn't exist, TTL of -1 means no expiry.
	// Normally ttl should be positive. Let's calculate expiresAt.
	expiresAt := time.Now().Add(ttl).Unix()

	if strings.HasPrefix(val, "revoked:") {
		userIDStr := strings.TrimPrefix(val, "revoked:")
		userID, parseErr := strconv.ParseUint(userIDStr, 10, 64)
		if parseErr != nil {
			return nil, parseErr
		}
		return &domain.RefreshToken{
			UserID:    uint(userID),
			TokenHash: tokenHash,
			ExpiresAt: expiresAt,
			Revoked:   true,
		}, nil
	}

	userID, parseErr := strconv.ParseUint(val, 10, 64)
	if parseErr != nil {
		return nil, parseErr
	}

	return &domain.RefreshToken{
		UserID:    uint(userID),
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}, nil
}

func (r *redisRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	key := "refresh:" + tokenHash
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}

	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil || ttl <= 0 {
		return err
	}

	if !strings.HasPrefix(val, "revoked:") {
		revVal := "revoked:" + val
		return r.client.Set(ctx, key, revVal, ttl).Err()
	}

	return nil
}

func (r *redisRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uint) error {
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	userSetKey := "user_tokens:" + userIDStr

	hashes, err := r.client.SMembers(ctx, userSetKey).Result()
	if err != nil {
		return err
	}

	for _, hash := range hashes {
		key := "refresh:" + hash
		val, err := r.client.Get(ctx, key).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return err
		}

		ttl, err := r.client.TTL(ctx, key).Result()
		if err != nil || ttl <= 0 {
			continue
		}

		if !strings.HasPrefix(val, "revoked:") {
			revVal := "revoked:" + val
			if setErr := r.client.Set(ctx, key, revVal, ttl).Err(); setErr != nil {
				return setErr
			}
		}
	}

	return r.client.Del(ctx, userSetKey).Err()
}
