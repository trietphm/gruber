# gruber

A simple API for gruber
===

## Install

- Golang
- Postgresql
- Setup project:
 - Install postgresql migration tool `goose`: `$ go get bitbucket.org/liamstask/goose/cmd/goose`
 - Install cassandra migration tool `migrate`:
 ```
	$ go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
	$ go build -tags 'cassandra' -o /usr/local/bin/migrate github.com/mattes/migrate/cli
 ```
 - Install `dep`: `$ go get -u github.com/golang/dep/cmd/dep`
 - Update project dependencies: `$ dep ensure`
 - Copy migration postgresql config example `migration/postgresql/db/dbconf.yml.example` to `migration/db/dbconf.yml` and update with properly config.
 - Copy migration cassandra config example `migration/cassandra/db/dbconf.yml.example` to `migration/db/dbconf.yml` and update with properly config.
 - Running migration: `$ goose up`
 - Copy app config example `config/dbconf.yml.example` to `migration/db/dbconf.yml` and update with properly config.

## Usage

- Build app: `go build -o gruber`
- Running:
 - With config file: `./gruber -config=config/configuration` (the later `configuration` is config file's name without extension `yaml`)
 - Or running with ENV variable: `./gruber` (see `config/configuration.yaml` for more information about ENV variables)
