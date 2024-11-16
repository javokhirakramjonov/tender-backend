package main

import (
	"log"
	"tender-backend/api"
	"tender-backend/config"
	"tender-backend/db"
	"tender-backend/internal/http/handlers"
)

func main() {
	config.LoadConfig()

	db.ConnectDB()
	defer db.CloseDB()

	h := handlers.NewHttpHandler(db.DB)

	r := api.NewGinRouter(h)

	err := r.Run(config.GlobalConfig.AppPort)
	if err != nil {
		log.Fatal(err)
	}
}
