#build
FROM golang:1.25.1-alpine AS build

WORKDIR /src

COPY paymentsvc/go.mod paymentsvc/go.sum /src/paymentsvc/

WORKDIR /src/paymentsvc

RUN go mod download

WORKDIR /src

COPY . .

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /src/paymentsvc

RUN CGO_ENABLED=0 go build -o /bin/paymentsvc_binary

RUN chmod +x /src/paymentsvc/entrypoint.sh

#run
FROM alpine:latest

WORKDIR /paymentsvc

COPY --from=build /bin/paymentsvc_binary /paymentsvc/

COPY --from=build /src/paymentsvc/migrations /paymentsvc/migrations

COPY --from=build /go/bin/migrate /usr/local/bin/migrate

COPY --from=build /src/paymentsvc/entrypoint.sh /paymentsvc/entrypoint.sh

ENTRYPOINT ["/paymentsvc/entrypoint.sh", "postgres", "app","app","payments_db"]
