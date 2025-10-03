#build
FROM golang:1.25.1-alpine AS build

WORKDIR /src

# Cache deps
COPY ordersvc/go.mod ordersvc/go.sum /src/ordersvc/
COPY proto/go.mod proto/go.sum /src/proto/

WORKDIR /src/ordersvc
RUN go mod download

# Copy source
WORKDIR /src
COPY . .

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /src/ordersvc
RUN CGO_ENABLED=0 go build -o /bin/ordersvc_binary

RUN ["chmod", "+x", "/src/ordersvc/entrypoint.sh"]

#run
FROM alpine:latest

WORKDIR /ordersvc

COPY --from=build /bin/ordersvc_binary /ordersvc/

COPY --from=build /go/bin/migrate /usr/local/bin/migrate

COPY --from=build /src/ordersvc//migrations /ordersvc/migrations

COPY --from=build /src/ordersvc/entrypoint.sh /ordersvc/entrypoint.sh

ENV GRPC_ADDR=:50051

ENTRYPOINT [ "/ordersvc/entrypoint.sh", "postgres", "app","app","orders_db"]
