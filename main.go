package main

import (
	"database/sql"
	"github.com/petershivachi/go_transact/api"
	db "github.com/petershivachi/go_transact/db/sqlc"
	"github.com/petershivachi/go_transact/util"
	"log"

	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database : ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server : ", err)
	}

}
