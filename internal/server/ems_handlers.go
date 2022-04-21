package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jeandeducla/api-plant/internal/plants"
)

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
