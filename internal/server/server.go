package server

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jeandeducla/api-plant/internal/plants"
)

type Server struct {
    plantsService plants.ServiceInterface
}

func NewServer(plantsService plants.ServiceInterface) (*Server, error) {
    return &Server{
        plantsService: plantsService,
    }, nil
}

func (s *Server) Router() *gin.Engine {
    router := gin.Default()

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

    return router
}

func parseId(ctx *gin.Context) (uint, error) {
    param := ctx.Param("id")
    id, err := strconv.ParseUint(param, 0, 64)
    return uint(id), err
}

func matchError(err error) (int, error) {
    if errors.Is(err, plants.ErrEmptyResult) {
        return 404, err
    } else if err != nil {
        return 500, err
    }
    return 0, nil
}
