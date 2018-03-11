#!/bin/sh
# Config cassandra
host="localhost"
port="9042"
username="gruber"
password="password"
keyspace="gruber"

action=$1
path=$2

if [ "$action" = "up" ]; then
 echo "Migrating up..."
 migrate -database "cassandra://$host:$port/$keyspace?username=$username&password=$password" -path migrations $1
elif [ "$action" = "create" ]; then
 echo "Create new migration"
 migrate -database "cassandra://$host:$port/$keyspace?username=$username&password=$password" -path migrations $1 $2
elif [ "$action" = "down" ]; then
 echo "Migrating down..."
 migrate -database "cassandra://$host:$port/$keyspace?username=$username&password=$password" -path migrations $1
elif [ "$action" = "version" ]; then
 migrate -database "cassandra://$host:$port/$keyspace?username=$username&password=$password" -path migrations $1
fi
