package main

import (
	"github.com/Yukio0315/reservation-backend/src/db"
	"github.com/Yukio0315/reservation-backend/src/server"
)

func main() {
	db.Init()
	server.Init()
}
