package main

import (
	"tender-backend/db"
)

func main() {

	db.ConnectDB()
	defer db.CloseDB()
}
