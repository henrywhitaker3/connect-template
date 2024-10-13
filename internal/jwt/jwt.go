package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/henrywhitaker3/connect-template/internal/crypto"
	"github.com/redis/rueidis"
)

var (
	ErrInvalidated = errors.New("jwt has been invalidated")
)

type Jwt struct {
	secret string
	redis  rueidis.Client
}

func New(secret string, redis rueidis.Client) *Jwt {
	return &Jwt{
		secret: secret,
		redis:  redis,
	}
}

type claims struct {
	jwt.StandardClaims
}

func (j *Jwt) New(expires time.Duration) (string, error) {
	exp := time.Now().Add(expires)

	claims := claims{
		jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (j *Jwt) Verify(ctx context.Context, token string) error {
	hash := crypto.Sum(token)

	// Check if the token has been invalidated first
	cmd := j.redis.B().Get().Key(j.invalidatedKey(hash)).Build()
	res := j.redis.Do(ctx, cmd)
	if err := res.Error(); err != nil {
		if !errors.Is(err, rueidis.Nil) {
			return err
		}
	} else {
		return ErrInvalidated
	}

	_, err := j.getClaims(token)
	if err != nil {
		return err
	}
	return nil
}

func (j *Jwt) getClaims(token string) (*claims, error) {
	claims := &claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (j *Jwt) Invalidate(ctx context.Context, token string) error {
	claims, err := j.getClaims(token)
	if err != nil {
		return err
	}

	expires := time.Unix(claims.ExpiresAt, 0)
	remaining := time.Until(expires)

	cmd := j.redis.B().Set().Key(j.invalidatedKey(crypto.Sum(token))).Value("true").Ex(remaining).Build()
	res := j.redis.Do(ctx, cmd)
	return res.Error()
}

func (j *Jwt) invalidatedKey(hash string) string {
	return fmt.Sprintf("tokens:invalidated:%s", hash)
}
