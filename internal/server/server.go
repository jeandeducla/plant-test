package server

import (
	"errors"
	"net/http"
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

func (s *Server) handleGetEnergyManagers(ctx *gin.Context) {
    res, err := s.plantsService.GetAllEnergyManagers()
    if err != nil {
        ctx.AbortWithStatus(500)
        return
    }
    ctx.JSON(http.StatusOK, res)
}

func (s *Server) handleGetEnergyManager(ctx *gin.Context) {
    id, err := parseId(ctx)
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }
    
    res, err := s.plantsService.GetEnergyManager(id)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.JSON(http.StatusOK, res)
}

func (s *Server) handleDeleteEnergyManager(ctx *gin.Context) {
    id, err := parseId(ctx)
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }

    err = s.plantsService.DeleteEnergyManager(id)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.String(http.StatusOK, "")
}

func (s *Server) handlePostEnergyManager(ctx *gin.Context) {
    var input plants.CreateEnergyManagerInput
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.String(http.StatusBadRequest, "")
        return
    }

    err := s.plantsService.CreateEnergyManager(input)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.String(http.StatusOK, "")
}

func (s *Server) handlePutEnergyManager(ctx *gin.Context) {
    id, err := parseId(ctx)
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }

    var input plants.UpdateEnergyManagerInput
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.String(http.StatusBadRequest, "")
        return
    }

    err = s.plantsService.UpdateEnergyManager(id, input)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.String(http.StatusOK, "")
}

func (s *Server) handleGetEnergyManagerPlants(ctx *gin.Context) {
    id, err := parseId(ctx)
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }

    res, err := s.plantsService.GetEnergyManagerPlants(id)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.JSON(http.StatusOK, res)
}

func (s *Server) handleGetPlants(ctx *gin.Context) {
    res, err := s.plantsService.GetAllPlants()
    if err != nil {
        ctx.AbortWithStatus(500)
        return
    }
    ctx.JSON(http.StatusOK, res)
}

func (s *Server) handlePostPlant(ctx *gin.Context) {
    var input plants.CreatePlantInput
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.String(http.StatusBadRequest, "")
        return
    }

    err := s.plantsService.CreatePlant(input)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.String(http.StatusOK, "")
}

func (s *Server) handleGetPlant(ctx *gin.Context) {
    id, err := parseId(ctx)
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }
    
    res, err := s.plantsService.GetPlant(id)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.JSON(http.StatusOK, res)
}

func (s *Server) handleDeletePlant(ctx *gin.Context) {
    id, err := parseId(ctx)
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }

    err = s.plantsService.DeletePlant(id)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.String(http.StatusOK, "")
}

func (s *Server) handlePutPlant(ctx *gin.Context) {
    id, err := parseId(ctx)
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }

    var input plants.UpdatePlantInput
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.String(http.StatusBadRequest, "")
        return
    }

    err = s.plantsService.UpdatePlant(id, input)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.String(http.StatusOK, "")
}
