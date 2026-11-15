package main

import (
	"log"
	"net/http"
	"os"

	"golab/lesson09/internal/handler"
	"golab/lesson09/internal/repository"
	"golab/lesson09/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	repo := repository.NewMemory()
	svc := service.New(repo)
	h := handler.New(svc)
	root := handler.Chain(h.Routes(), handler.Logging, handler.Recover)

	log.Printf("Task API 서버 시작: http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, root))
}
