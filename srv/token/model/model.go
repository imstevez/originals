package model

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

type TokenModel struct {
	Redis *redis.Client
}

const (
	cancelTokenScope = `CancelToken`
	freshTokenScope  = `FreshToken`
)

var ErrKeyNotExist = errors.New("key is not exist")

// CancelToken
func (mdl *TokenModel) CancelToken(token string, expiration time.Duration) error {
	key := cancelTokenScope + ":" + token
	if err := mdl.Redis.Set(key, token, expiration).Err(); err != nil {
		return err
	}
	return nil
}

// IsTokenCanceled
func (mdl *TokenModel) IsTokenCanceled(token string) (bool, error) {
	key := cancelTokenScope + ":" + token
	val, err := mdl.Redis.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, err
}

// GetRefreshToken
func (mdl *TokenModel) GetFreshToken(token string) (string, error) {
	key := freshTokenScope + ":" + token
	val, err := mdl.Redis.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrKeyNotExist
		}
		return "", err
	}
	return val, err
}

// SetRefreshToken
func (mdl *TokenModel) SetRefreshToken(oldToken, newToken string, expiration time.Duration) error {
	key := freshTokenScope + ":" + oldToken
	if err := mdl.Redis.Set(key, newToken, expiration).Err(); err != nil {
		return err
	}
	return nil
}
