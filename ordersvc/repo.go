package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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
	return [...]string{"unspecified", "pending", "paid", "failed"}[p]
}

func (p *PaymentStatus) Scan(v interface{}) error {
	var status PaymentStatus
	var err error
	switch val := v.(type) {
	case string:
		status, err = parsePaymentStatus(val)
		*p = status
	case []byte:
		status, err = parsePaymentStatus(string(val))
		*p = status
	default:
		return fmt.Errorf("failed to scan PaymentStatus. Err:%v", err)
	}

	return err
}

func parsePaymentStatus(v string) (PaymentStatus, error) {
	switch strings.ToLower(v) {
	case "unspecified":
		return Payment_Unspecified, nil
	case "pending":
		return Payment_Pending, nil
	case "paid":
		return Payment_Paid, nil
	case "failed":
		return Payment_Failed, nil
	}

	return Payment_Unspecified, fmt.Errorf("invalid payment status string %s", v)
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

	args := []any{order.UserId, order.AmountCents, order.PaymentStatus.String(), time.Now().Unix()}

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&order.Id)
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
