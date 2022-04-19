package main

import (
	"github.com/jeandeducla/api-plant/internal/models"
	"github.com/jeandeducla/api-plant/internal/server"
)

func main() {
    config := NewConfig()

    db, err := models.NewDB(config.dsn)
    if err != nil {
        panic(err)
    }

    server, err := server.NewServer(db)
    if err != nil {
        panic(err)
    }

    server.Router()
}
