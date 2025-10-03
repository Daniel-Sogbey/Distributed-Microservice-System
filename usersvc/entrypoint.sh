#!/bin/sh
#
set -e

host="$1"
user="$2"
password="$3"
db_name="$4"

echo "host => $host : user => $user : db_name => $db_name"

# until pg_isready -h "$host" -U "$user";do
#     echo "Postgres is anavailable - sleeping for 2 seconds"
#     sleep 2
# done

echo "running migrations"
until migrate -path ./migrations -database "postgres://$user:$password@$host/$db_name?sslmode=disable" up; do
    echo "migrations failed - retrying in 3s..."
    sleep 3
done

echo "starting users gRPC service"
exec ./usersvc_binary
