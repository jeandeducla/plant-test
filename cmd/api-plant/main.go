package main

import (
	"github.com/jeandeducla/api-plant/internal/plants"
	"github.com/jeandeducla/api-plant/internal/models"
	"github.com/jeandeducla/api-plant/internal/server"
)

func main() {
    config := NewConfig()

    // DB connection and init
    db, err := models.NewDB(config.dsn)
    if err != nil {
        panic(err)
    }

    // pure ORM layer
    plantsDB := plants.NewPlantsDB(db)

    // Business logic layer
    plantsService := plants.NewPlantsService(plantsDB)

    // http layer
    server, err := server.NewServer(plantsService)
    if err != nil {
        panic(err)
    }

    server.Router()
}
