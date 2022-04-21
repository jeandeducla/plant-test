package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jeandeducla/api-plant/internal/plants"
)

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
    id, err := parseId(ctx, "id")
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
    id, err := parseId(ctx, "id")
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
    id, err := parseId(ctx, "id")
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
