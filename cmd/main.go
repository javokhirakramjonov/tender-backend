package main

import (
	"tender-backend/config"
	"tender-backend/db"
)

func main() {
	cnf := config.LoadConfig()

	db.ConnectDB(cnf.DB)
	defer db.CloseDB()
}
