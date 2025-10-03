package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	userv1 "github.com/Daniel-Sogbey/micro-weekend/proto/user/v1"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func mustEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}

	return v
}

func main() {
	dsn := mustEnv("DB_DSN", "postgres://app:app@postgres:5432/users_db?sslmode=disable")
	grpcAddr := mustEnv("GRPC_ADDR", ":50051")

	db, err := connectDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("database connection established successfully")

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()

	userv1.RegisterUserServiceServer(srv, &UserServer{repo: &Repo{DB: db}})
	log.Printf("usersvc listening at %s\n", grpcAddr)
	log.Fatal(srv.Serve(lis))
}

func connectDB(dsn string) (*sql.DB, error) {
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

		time.Sleep(3 * time.Second)
		log.Println("reconnecting to db after failed attempt...")
	}

}
