package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	eventsv1 "github.com/Daniel-Sogbey/micro-weekend/proto/events/v1"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type PaymentService struct {
	rabbitmq *RabbitMQ
	repo     *Repo
}

func (p PaymentService) consumeOrderCreated() error {
	err := p.rabbitmq.ExchangeDeclare("orders", EXCHANGE_FANOUT, true, false, false, false, nil)
	if err != nil {
		return err
	}

	q, err := p.rabbitmq.QueueDeclare("order_created", true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = p.rabbitmq.QueueBind(q, "order.created", "orders", false, nil)
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)

	return p.rabbitmq.Consume(context.Background(), "order_created_consumer", q.Name, "orders", "order.created",
		EXCHANGE_FANOUT, true, false, false, false, nil, func(b []byte) error {
			log.Printf(">>>>>>>>>>>> [x] Received message %s\n", string(b))

			var orderEvent eventsv1.OrderEvent

			if err := proto.Unmarshal(b, &orderEvent); err != nil {
				log.Println("unmarshal:", err)
				return err
			}

			payment := &Payment{
				OrderId:     orderEvent.OrderId,
				Status:      string(orderEvent.PaymentStatus),
				AmountCents: orderEvent.AmountCents,
			}

			return p.publishOrderProcessed(payment)
		})
}

func (p PaymentService) publishOrderProcessed(payment *Payment) error {
	log.Printf(">>>>>>>>>>>> [x] Processing order record %v", payment)

	time.Sleep(2 * time.Second)

	payment.Status = "paid"

	return p.repo.Create(context.Background(), payment)
}

func mustEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	dsn := mustEnv("DB_URL", "postgres://app:app@postgres:5432/payments_db?sslmode=disable")
	rabbitMQURL := mustEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")

	conn, err := connectToRabbitMQ(rabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer channel.Close()

	db, err := connectToDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	psvc := PaymentService{
		rabbitmq: &RabbitMQ{conn: conn, channel: channel},
		repo:     &Repo{DB: db},
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)

	go func() {
		if err := psvc.consumeOrderCreated(); err != nil {
			stop()
		}
	}()

	<-ctx.Done()
	log.Printf(">>>>>>>> [X] Received signal to shutdown server gracefully...Signal: %v\n", ctx.Err())
}

func connectToRabbitMQ(url string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}

		log.Printf("Failed to establish connection to rabbitmq. Retrying in 3 seconds...Err:%v\n", err)
		time.Sleep(3 * time.Second)
	}

	return conn, err
}

func connectToDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := db.PingContext(ctx)
		cancel()

		if err == nil {
			return db, nil
		}

		log.Printf("Failed to establish connection to postgres. Retrying in 3 seconds...Err: %v\n", err)
		time.Sleep(3 * time.Second)
	}
}
