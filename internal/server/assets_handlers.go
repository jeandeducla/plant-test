package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jeandeducla/api-plant/internal/plants"
)

func (s *Server) handleGetPlantAssets(ctx *gin.Context) {
    id, err := parseId(ctx)
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }

    res, err := s.plantsService.GetPlantAssets(id)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.JSON(http.StatusOK, res)
}

func (s *Server) handlePostAsset(ctx *gin.Context) {
    id, err := parseId(ctx)
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }

    var input plants.CreateAssetInput
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.String(http.StatusBadRequest, "")
        return
    }

    err = s.plantsService.CreateAsset(id, input)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.String(http.StatusOK, "")
}
