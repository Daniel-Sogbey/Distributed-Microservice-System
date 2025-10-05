package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	orderv1 "github.com/Daniel-Sogbey/micro-weekend/proto/order/v1"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

func mustEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return v
	}
	return def
}

func main() {
	dsn := mustEnv("DB_DSN", "postgres://app:app@postgres:5432/orders_db?sslmode=disable")
	grpcAddr := mustEnv("GRPC_ADDR", ":50052")
	rabbitMQURL := mustEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672")

	amqpConn, err := connectToRabbitMQ(rabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}
	defer amqpConn.Close()

	channel, err := amqpConn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer channel.Close()

	db, err := connectToDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(srv, &OrderServer{
		repo:     &Repo{DB: db},
		rabbitMQ: &RabbitMQ{conn: amqpConn, channel: channel},
	})
	log.Printf("ordersvc listening at %s\n", grpcAddr)
	log.Fatal(srv.Serve(lis))
}

func connectToRabbitMQ(url string) (*amqp.Connection, error) {
	var amqpConn *amqp.Connection
	var err error
	for {
		amqpConn, err = amqp.Dial(url)

		if err == nil {
			break
		}

		log.Printf("failed to connect to rabbitmq. Retrying in 3 seconds...%v ->%s\n", err, url)
		time.Sleep(3 * time.Second)
	}

	return amqpConn, err
}

func connectToDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err = db.PingContext(ctx)
		cancel()

		if err == nil {
			return db, nil
		}

		log.Println("failed to connect to postgres. Retrying in 3 seconds...")
		time.Sleep(3 * time.Second)
	}
}
