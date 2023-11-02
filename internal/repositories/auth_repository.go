package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-redis/redis"
)

type mysqlAuthRepository struct {
	conn  *sql.DB
	redis *redis.Client
}

func NewMysqlAuthRepository(conn *sql.DB, redis *redis.Client) AuthRepository {
	return &mysqlAuthRepository{
		conn:  conn,
		redis: redis,
	}
}

type AuthRepository interface {
	SaveOtpSecret(ctx context.Context, uuid string, secretKey string) (err error)
	GetOtpSecret(ctx context.Context, uuid string) (secretKey string, expiredAt time.Time, err error)
}
