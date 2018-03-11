#!/bin/sh

action=$1
path=$2

if [ "$action" = "up" ]; then
 echo "Migrating up..."
 migrate -database "cassandra://localhost:9042/gruber?username=gruber&password=password" -path migrations $1
elif [ "$action" = "create" ]; then
 echo "Create new migration"
 migrate -database "cassandra://localhost:9042/gruber?username=gruber&password=password" -path migrations $1 $2
elif [ "$action" = "down" ]; then
 echo "Migrating down..."
 migrate -database "cassandra://localhost:9042/gruber?username=gruber&password=password" -path migrations $1
elif [ "$action" = "version" ]; then
 migrate -database "cassandra://localhost:9042/gruber?username=gruber&password=password" -path migrations $1
fi
