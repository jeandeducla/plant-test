package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jeandeducla/api-plant/internal/plants"
)

type Server struct {
    plantsService *plants.Service
}

func NewServer(plantsService *plants.Service) (*Server, error) {
    return &Server{
        plantsService: plantsService,
    }, nil
}

func (s *Server) Router() {
    router := gin.Default()

    router.GET("/ems", s.handleReadEnergyManagers)
    router.GET("/ems/:id", s.handleReadEnergyManager)

    router.Run()
}

func (s *Server) handleReadEnergyManagers(ctx *gin.Context) {
    res, err := s.plantsService.GetAllEnergyManagers()
    if err != nil {
        ctx.AbortWithStatus(500)
        return
    }
    ctx.JSON(http.StatusOK, res)
}

func (s *Server) handleReadEnergyManager(ctx *gin.Context) {
    param := ctx.Param("id")
    id, err := strconv.ParseUint(param, 0, 64)
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }
    
    res, err := s.plantsService.GetEnergyManager(uint(id))
    if err != nil {
        ctx.AbortWithStatus(500)
        return
    }
    if res == nil {
        ctx.AbortWithStatus(404)
        return
    }
    ctx.JSON(http.StatusOK, res)
}
