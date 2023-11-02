package repository

import (
	"context"
	"time"
)

func (m *mysqlAuthRepository) SaveOtpSecret(ctx context.Context, uuid string, secretKey string) (err error) {
	query := `INSERT INTO banking.users_otp SET uuid=?, secret=?, created_at=?, expired_at=?`
	stmt, err := m.conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, uuid, secretKey, time.Now(), time.Now().Add(180*time.Second))
	if err != nil {
		return
	}
	return err
}
