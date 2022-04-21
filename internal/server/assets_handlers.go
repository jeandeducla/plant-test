package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jeandeducla/api-plant/internal/plants"
)

func (s *Server) handleGetPlantAssets(ctx *gin.Context) {
    id, err := parseId(ctx, "id")
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
    id, err := parseId(ctx, "id")
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

func (s *Server) handleGetPlantAsset(ctx *gin.Context) {
    plant_id, err := parseId(ctx, "id")
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }
    asset_id, err := parseId(ctx, "asset_id")
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }

    res, err := s.plantsService.GetPlantAsset(plant_id, asset_id)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.JSON(http.StatusOK, res)
}

func (s *Server) handleDeletePlantAsset(ctx *gin.Context) {
    plant_id, err := parseId(ctx, "id")
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }
    asset_id, err := parseId(ctx, "asset_id")
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }

    err = s.plantsService.DeletePlantAsset(plant_id, asset_id)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.String(http.StatusOK, "")
}

func (s *Server) handlePutPlantAsset(ctx *gin.Context) {
    plant_id, err := parseId(ctx, "id")
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }
    asset_id, err := parseId(ctx, "asset_id")
    if err != nil {
        ctx.AbortWithStatus(404)
        return
    }

    var input plants.UpdateAssetInput
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.String(http.StatusBadRequest, "")
        return
    }

    err = s.plantsService.UpdatePlantAsset(plant_id, asset_id, input)
    status, err := matchError(err)
    if err != nil {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.String(http.StatusOK, "")
}
