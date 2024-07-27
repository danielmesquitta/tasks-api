package main

import "github.com/danielmesquitta/tasks-api/internal/app/http"

// @title Tasks API
// @version 1.0
// @description This is a CRUD API for tasks.
// @contact.name Daniel Mesquita
// @contact.email danielmesquitta123@gmail.com
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @securityDefinitions.basic BasicAuth
func main() {
	http.Start()
}
