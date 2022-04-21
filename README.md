# api-plant

## Requirements
- `docker`
- `docker-compose`

## Run

First run:
```$xslt
    $ docker-compose up --build
```
which will pop the api server (listening on port 8080) and a postgresql.

Then you can `curl` your localhost like that:
```$xslt
    $ curl localhost:8080/ems
```

## Code structure

This service has three layers:
- the http layer in `./internal/server/server.go`: all what is http related
- the business layer in `./internal/plants/plants_service.go`: all the business logic happens here
- the DB or ORM layer in `./internal/plants/plants_db.go`: all that is pure database related is here

## Routes

The server router has this shape:
```
...
    router.GET("/ems", s.handleGetEnergyManagers)
    router.POST("/ems", s.handlePostEnergyManager)
    router.GET("/ems/:id", s.handleGetEnergyManager)
    router.DELETE("/ems/:id", s.handleDeleteEnergyManager)
    router.PUT("/ems/:id", s.handlePutEnergyManager)

    router.GET("/ems/:id/plants", s.handleGetEnergyManagerPlants)

    router.GET("/plants", s.handleGetPlants)
    router.POST("/plants", s.handlePostPlant)
    router.GET("/plants/:id", s.handleGetPlant)
    router.DELETE("/plants/:id", s.handleDeletePlant)
    router.PUT("/plants/:id", s.handlePutPlant)

    router.GET("/plants/:id/assets", s.handleGetPlantAssets)
    router.POST("/plants/:id/assets", s.handlePostAsset)
    router.GET("/plants/:id/assets/:asset_id", s.handleGetPlantAsset)
    router.DELETE("/plants/:id/assets/:asset_id", s.handleDeletePlantAsset)
    router.PUT("/plants/:id/assets/:asset_id", s.handlePutPlantAsset)
...

```

For example to create a new Energy Manager you could do:
```$xslt
    $ curl -X POST -d '{"name": "jack", "surname": "chirak"}' localhost:8080/ems
```
To see the assets of a specific plant:
```$xslt
    $ curl localhost:8080/plants/1/assets
```


## Test

To run the http layer tests, run:
```$xslt
    $ docker-compose -f docker-compose.test.yaml run test-server
```

To run the business layer tests, run:
```$xslt
    $ docker-compose -f docker-compose.test.yaml run test-plants
```
