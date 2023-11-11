package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/xiao-yangg/simplebank/api"
	db "github.com/xiao-yangg/simplebank/db/sqlc"
	"github.com/xiao-yangg/simplebank/db/util"
)

func main() {
	config, err := util.LoadConfig(".") // app.env currently same dir as main.go
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connection, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(connection)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}