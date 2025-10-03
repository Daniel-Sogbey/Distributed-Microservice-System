package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	orderv1 "github.com/Daniel-Sogbey/micro-weekend/proto/order/v1"
)

type OrderServer struct {
	orderv1.UnimplementedOrderServiceServer
	repo     *Repo
	rabbitMQ *RabbitMQ
}

func (o *OrderServer) CreateOrder(ctx context.Context, in *orderv1.CreateOrderRequest) (*orderv1.Order, error) {
	order := Order{
		UserId:      in.UserId,
		AmountCents: in.AmountCents,
	}

	err := o.repo.Create(ctx, &order)
	if err != nil {
		return nil, err
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error marshalling order %v\n", err)
		}

		err = o.rabbitMQ.Publish(ctx, body, "orders", "", EXCHANGE_FANOUT, true, false, false, false, nil)
		if err != nil {
			log.Printf("Error publishing order created event %v\n", err)
		}

	}()

	return &orderv1.Order{Id: order.Id, UserId: order.UserId, AmountCents: order.AmountCents, Status: order.PaymentStatus.String(), CreatedUnix: order.CreatedUnix}, nil
}

func (o *OrderServer) GetOrder(ctx context.Context, in *orderv1.GetOrderRequest) (*orderv1.Order, error) {
	order, err := o.repo.Get(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &orderv1.Order{Id: order.Id, UserId: order.UserId, AmountCents: order.AmountCents, Status: order.PaymentStatus.String(), CreatedUnix: order.CreatedUnix}, nil
}
