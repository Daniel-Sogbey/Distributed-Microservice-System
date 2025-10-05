#build
FROM golang:1.25.1-alpine AS build

WORKDIR /src

COPY gateway/go.mod gateway/go.sum /src/gateway/
COPY proto/go.mod proto/go.sum /src/proto/

WORKDIR /src/gateway

RUN go mod download

WORKDIR /src

COPY . .

WORKDIR /src/gateway

RUN CGO_ENABLED=0 go build -o /bin/gateway_binary

#run
FROM alpine:latest

WORKDIR /gateway

COPY --from=build /bin/gateway_binary /gateway/gateway_binary

ENTRYPOINT ["/gateway/gateway_binary"]
