package repository

import (
	"context"
	"log"
	"time"
)

func (m *mysqlAuthRepository) GetOtpSecret(ctx context.Context, uuid string) (secretKey string, expiredAt time.Time, err error) {
	query := `SELECT secret, expired_at FROM banking.users_otp WHERE uuid = ?`

	rows, err := m.conn.Query(query, uuid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&secretKey, &expiredAt)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return secretKey, expiredAt, nil
}
