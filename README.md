# gruber [![Build Status](https://travis-ci.org/trietphm/gruber.svg?branch=master)](https://travis-ci.org/trietphm/gruber)

A simple API for gruber
===

## Install

- Golang
- Postgresql. Create user and database gruber.
- Cassandra. Create user and keyspace gruber:
```
CREATE KEYSPACE gruber WITH replication =  {'class': 'SimpleStrategy', 'replication_factor': '1'}  AND durable_writes = true;
```
- Setup project:
 - Install postgresql migration tool `goose`: `$ go get bitbucket.org/liamstask/goose/cmd/goose`
 - Install cassandra migration tool `migrate`:
 ```
	$ go get -u -d github.com/mattes/migrate/cli github.com/gocql/gocql
	$ go build -tags 'cassandra' -o $GOPATH/bin/migrate github.com/mattes/migrate/cli
 ```
 - Install `dep`: `$ go get -u github.com/golang/dep/cmd/dep`
 - Update project dependencies: `$ dep ensure`
 - Create database in Postgresql and Cassandra before running migration.
 - Copy migration postgresql config example `migration/postgresql/db/dbconf.yml.example` to `migration/db/dbconf.yml` and update with properly config.
 - Update migration cassandra config `migration/cassandra/migrate.sh`.
 - Running migration Postgresql:
 ```
 $ cd migration/postgresql
 $ goose up
 ```
 - Running migration Cassandra:
 ```
 $ cd migration/cassandra
 $ ./migrate.sh up
 ```
 - Copy app config example `config/dbconf.yml.example` to `migration/db/dbconf.yml` and update with properly config.

## Usage

- Build app: `go build -o gruber`
- Running:
 - With config file: `./gruber -config=config/configuration` (the later `configuration` is config file's name without extension `yaml`)
 - Or running with ENV variable: `./gruber` (see `config/configuration.yaml.example` for more information about ENV variables)
