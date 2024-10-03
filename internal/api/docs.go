// api/docs.go

// Package api API Server for Uniswap Fee Tracker
//
// Documentation for API endpoints.
//
//	Schemes: http
//	BasePath: /api/v1
//	Version: 1.0.0
//	Contact: Andhika Satriya Bhayangkara <adinbhayangkara@gmail.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Swagger: http://localhost:8080/swagger/index.html
package api

// ErrorResponse represents the structure of error responses.
// swagger:model
type ErrorResponse struct {
	Error string `json:"error"`
}
