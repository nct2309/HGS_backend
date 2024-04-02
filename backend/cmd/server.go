package cmd

import (
	"github.com/gin-gonic/gin"
)

func NewServer() Server {
	return server{
		router: gin.New(),
	}
}

type Server interface {
	Start()
}

type server struct {
	router *gin.Engine
}
