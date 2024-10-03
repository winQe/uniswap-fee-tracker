package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/winQe/uniswap-fee-tracker/internal/api"
)

// Server represents the API server and route handlers
type Server struct {
	port            string
	txHandler       *api.TransactionHandler
	batchJobHandler *api.BatchJobHandler
}

// Server represents the API server and route handlers
func NewServer(port string, txHandler *api.TransactionHandler, batchJobHandler *api.BatchJobHandler) *Server {
	return &Server{
		port:            port,
		txHandler:       txHandler,
		batchJobHandler: batchJobHandler,
	}
}

// Run starts the HTTP server, registers API routes, and binds the server to the specified port.
func (s *Server) Run() error {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		api.RegisterRoutes(v1, s.txHandler, s.batchJobHandler)
	}

	serverAddr := fmt.Sprintf("0.0.0.0:%s", s.port)
	return router.Run(serverAddr)
}
