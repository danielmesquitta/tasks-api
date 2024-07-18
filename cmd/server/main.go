package main

import (
	"github.com/danielmesquitta/tasks-api/internal/app/http"
)

// @title Tasks API
// @version 1.0
// @description This is a CRUD API for tasks.
// @contact.name Daniel Mesquita
// @contact.email danielmesquitta123@gmail.com
// @BasePath /api/v1
func main() {
	http.Start()
}
