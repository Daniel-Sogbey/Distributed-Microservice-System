package main

import (
	"context"
	"database/sql"
	"time"
)

type PaymentStatus int8

const (
	Payment_Unspecified PaymentStatus = iota
	Payment_Pending
	Payment_Paid
	Payment_Failed
)

func (p PaymentStatus) String() string {
	return [...]string{"payment_unspecified", "payment_pending", "payment_paid", "payment_failed"}[p]
}

type Order struct {
	Id            string
	UserId        string
	AmountCents   int64
	PaymentStatus PaymentStatus
	CreatedUnix   int64
}

type Repo struct {
	DB *sql.DB
}

func (r *Repo) Create(ctx context.Context, order *Order) error {
	query := `INSERT INTO orders (user_id, amount_cents, status, created_unix) VALUES ($1,$2,$3,$4) RETURNING id`

	args := []any{order.UserId, order.AmountCents, order.PaymentStatus, time.Now().Unix()}

	return r.DB.QueryRowContext(ctx, query, args...).Scan(order.Id)
}

func (r *Repo) Get(ctx context.Context, id string) (*Order, error) {
	query := `SELECT id, user_id, amount_cents, status, created_unix FROM orders WHERE id=$1`

	var order Order

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&order.Id,
		&order.UserId,
		&order.AmountCents,
		&order.PaymentStatus,
		&order.CreatedUnix,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}
