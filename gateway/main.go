package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	orderv1 "github.com/Daniel-Sogbey/micro-weekend/proto/order/v1"
	userv1 "github.com/Daniel-Sogbey/micro-weekend/proto/user/v1"
	"github.com/julienschmidt/httprouter"
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
	httpAddr := mustEnv("HTTP_ADDR", ":8080")
	usersvcAddr := mustEnv("USERSVC_ADDR", ":50051")
	ordersvcAddr := mustEnv("ORDERSVC_ADDR", ":50052")

	uconn, err := grpc.NewClient(usersvcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	oconn, err := grpc.NewClient(ordersvcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	uclient := userv1.NewUserServiceClient(uconn)
	oclient := orderv1.NewOrderServiceClient(oconn)

	mux := httprouter.New()

	mux.HandlerFunc(http.MethodPost, "/users", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		out, err := uclient.CreateUser(r.Context(), &userv1.CreateUserRequest{
			Name:  input.Name,
			Email: input.Email,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(out)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandlerFunc(http.MethodGet, "/users/:id", func(w http.ResponseWriter, r *http.Request) {
		userId := httprouter.ParamsFromContext(r.Context()).ByName("id")
		if userId == "" {
			http.Error(w, "missing user id", http.StatusBadRequest)
			return
		}

		out, err := uclient.GetUser(r.Context(), &userv1.GetUserRequest{Id: userId})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(out)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	mux.HandlerFunc(http.MethodPost, "/orders", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			UserId      string `json:"user_id"`
			AmountCents int64  `json:"amount_cents"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		out, err := oclient.CreateOrder(r.Context(), &orderv1.CreateOrderRequest{
			UserId:      input.UserId,
			AmountCents: input.AmountCents,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = json.NewEncoder(w).Encode(out)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	mux.HandlerFunc(http.MethodGet, "/orders/:id", func(w http.ResponseWriter, r *http.Request) {
		orderId := httprouter.ParamsFromContext(r.Context()).ByName("id")
		if orderId == "" {
			http.Error(w, "missing order id", http.StatusBadRequest)
			return
		}

		out, err := oclient.GetOrder(r.Context(), &orderv1.GetOrderRequest{Id: orderId})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(out)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	log.Printf("API gateway listening on port: %s \n", httpAddr)
	log.Fatal(srv.ListenAndServe())

}
