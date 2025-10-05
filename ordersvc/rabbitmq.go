package main

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Exchange int8

const (
	EXCHANGE_DIRECT Exchange = iota
	EXCHANGE_TOPICS
	EXCHANGE_HEADERS
	EXCHANGE_FANOUT
)

func (e Exchange) String() string {
	return [...]string{"direct", "topic", "headers", "fanout"}[e]
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func (r RabbitMQ) Publish(ctx context.Context, body []byte, exName, routingKey string, kind Exchange, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	err := r.ExchangeDeclare(exName, kind, durable, autoDelete, internal, noWait, args)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(ctx, exName, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

func (r RabbitMQ) ExchangeDeclare(name string, kind Exchange, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return r.channel.ExchangeDeclare(name, kind.String(), durable, autoDelete, internal, noWait, args)
}
