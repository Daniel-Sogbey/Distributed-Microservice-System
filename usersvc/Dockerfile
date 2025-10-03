#build
FROM golang:1.25.1-alpine AS build

WORKDIR /src

# Cache deps
COPY usersvc/go.mod usersvc/go.sum /src/usersvc/
COPY proto/go.mod proto/go.sum /src/proto/

WORKDIR /src/usersvc
RUN go mod download

# Copy source
WORKDIR /src
COPY . .

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /src/usersvc
RUN CGO_ENABLED=0 go build -o /bin/usersvc_binary

RUN ["chmod", "+x", "/src/usersvc/entrypoint.sh"]

#run
FROM alpine:latest

WORKDIR /usersvc

COPY --from=build /bin/usersvc_binary /usersvc/

COPY --from=build /go/bin/migrate /usr/local/bin/migrate

COPY --from=build /src/usersvc//migrations /usersvc/migrations

COPY --from=build /src/usersvc/entrypoint.sh /usersvc/entrypoint.sh

ENV GRPC_ADDR=:50051

ENTRYPOINT [ "/usersvc/entrypoint.sh", "postgres", "app","app","users_db"]
