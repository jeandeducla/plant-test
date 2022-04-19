package server

import (
	"net/http"

    "gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

type Server struct {
    db *gorm.DB
}

func NewServer(db *gorm.DB) (*Server, error) {
    return &Server{
        db: db,
    }, nil
}

func (s *Server) Router() {
    router := gin.Default()

    router.GET("/", func(ctx *gin.Context) {
        ctx.JSON(http.StatusOK, gin.H{"data": "coucou"})
    })

    router.Run()
}
