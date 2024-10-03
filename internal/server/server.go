package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/winQe/uniswap-fee-tracker/internal/api"
)

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Run() error {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		api.RegisterRoutes(v1, nil, nil)
	}

	serverAddr := fmt.Sprintf("0.0.0.0:%s", s.port)
	return router.Run(serverAddr)
}
