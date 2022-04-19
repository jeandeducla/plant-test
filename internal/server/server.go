package server

import (
	"errors"
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
    router.POST("/ems", s.handleCreateEnergyManager)
    router.GET("/ems/:id", s.handleReadEnergyManager)
    router.DELETE("/ems/:id", s.handleDeleteEnergyManager)

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
    if errors.Is(err, plants.ErrEmptyResult) {
        ctx.AbortWithStatus(404)
        return
    } else if err != nil {
        ctx.AbortWithStatus(500)
        return
    }
    ctx.JSON(http.StatusOK, res)
}

func (s *Server) handleDeleteEnergyManager(ctx *gin.Context) {
    param := ctx.Param("id")
    id, err := strconv.ParseUint(param, 0, 64)
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }

    err = s.plantsService.DeleteEnergyManager(uint(id))
    if errors.Is(err, plants.ErrEmptyResult) {
        ctx.AbortWithStatus(404)
        return
    } else if err != nil {
        ctx.AbortWithStatus(500)
        return
    }
    ctx.String(http.StatusOK, "")
}

func (s *Server) handleCreateEnergyManager(ctx *gin.Context) {
    var input plants.CreateEnergyManagerInput
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.String(http.StatusBadRequest, "")
        return
    }

    err := s.plantsService.CreateEnergyManager(input)
    if errors.Is(err, plants.ErrEmptyResult) {
        ctx.AbortWithStatus(404)
        return
    } else if err != nil {
        ctx.AbortWithStatus(500)
        return
    }
    ctx.String(http.StatusOK, "")
}
