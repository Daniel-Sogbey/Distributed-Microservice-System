package main

import (
	"context"
	"log"
	"time"

	eventsv1 "github.com/Daniel-Sogbey/micro-weekend/proto/events/v1"
	orderv1 "github.com/Daniel-Sogbey/micro-weekend/proto/order/v1"
	"google.golang.org/protobuf/proto"
)

type OrderServer struct {
	orderv1.UnimplementedOrderServiceServer
	repo     *Repo
	rabbitMQ *RabbitMQ
}

func (o *OrderServer) CreateOrder(ctx context.Context, in *orderv1.CreateOrderRequest) (*orderv1.Order, error) {
	order := Order{
		UserId:        in.UserId,
		AmountCents:   in.AmountCents,
		PaymentStatus: Payment_Pending,
	}

	err := o.repo.Create(ctx, &order)
	if err != nil {
		return nil, err
	}

	var status orderv1.Status
	switch order.PaymentStatus {
	case Payment_Pending:
		status = orderv1.Status_STATUS_PENDING
	case Payment_Paid:
		status = orderv1.Status_STATUS_PAID
	case Payment_Failed:
		status = orderv1.Status_STATUS_FAILED
	default:
		status = orderv1.Status_STATUS_UNSPECIFIED
	}

	orderEvent := &eventsv1.OrderEvent{
		OrderId:       order.Id,
		AmountCents:   order.AmountCents,
		UserId:        order.UserId,
		PaymentStatus: status,
		CreatedUnix:   order.CreatedUnix,
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body, err := proto.Marshal(orderEvent)
		if err != nil {
			log.Printf("Error marshalling order %v\n", err)
		}

		err = o.rabbitMQ.Publish(ctx, body, "orders", "order.created", EXCHANGE_FANOUT, true, false, false, false, nil)
		if err != nil {
			log.Printf("Error publishing order created event %v\n", err)
		}
	}()

	return &orderv1.Order{
		Id: order.Id, UserId: order.UserId, AmountCents: order.AmountCents, Status: status,
		CreatedUnix: order.CreatedUnix,
	}, nil
}

func (o *OrderServer) GetOrder(ctx context.Context, in *orderv1.GetOrderRequest) (*orderv1.Order, error) {
	order, err := o.repo.Get(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	var status orderv1.Status
	switch order.PaymentStatus {
	case Payment_Pending:
		status = orderv1.Status_STATUS_PENDING
	case Payment_Paid:
		status = orderv1.Status_STATUS_PAID
	case Payment_Failed:
		status = orderv1.Status_STATUS_FAILED
	default:
		status = orderv1.Status_STATUS_UNSPECIFIED
	}

	return &orderv1.Order{
		Id: order.Id, UserId: order.UserId, AmountCents: order.AmountCents, Status: status,
		CreatedUnix: order.CreatedUnix,
	}, nil
}
