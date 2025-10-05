#!/bin/sh

set -e

host="$1"
user="$2"
password="$3"
db_name="$4"

echo "host => $host : user => $user : db_name => $db_name"

echo "Running migration..."
until migrate -path "/paymentsvc/migrations" -database "postgres://$user:$password@$host/$db_name?sslmode=disable" up; do
    echo "Migration failed. Retrying in 3 seconds ..."
    sleep 3
done


exec /paymentsvc/paymentsvc_binary
