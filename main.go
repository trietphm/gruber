package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/trietphm/gruber/app/handler"
	"github.com/trietphm/gruber/config"
	"github.com/trietphm/gruber/database"
)

func main() {
	var configPath = flag.String("config", "", "Set config file path")
	flag.Parse()

	conf, err := config.ReadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	dbPg, err := database.OpenPostgresqlDB(conf.Postgresql)
	if err != nil {
		panic(err)
	}

	dbCass, err := database.OpenCassandraDB(conf.Cassandra)
	if err != nil {
		panic(err)
	}

	dbRedis, err := database.OpenRedisDB(conf.Redis)
	if err != nil {
		panic(err)
	}

	engine, err := handler.NewEngine(dbPg, dbCass, dbRedis)
	if err != nil {
		panic(err)
	}

	if err := engine.Run(":" + strconv.Itoa(conf.App.Port)); err != nil {
		fmt.Println("Serve fail:", err)
	}
}
