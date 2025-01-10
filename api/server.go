package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/petershivachi/go_transact/db/sqlc"
)

// Server serves HTTP requests for our banking enigne
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a HTTP engine and sets up routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	server.router = router
	return server

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// Start runs HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
