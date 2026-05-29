package main

import (
	"log"
	"net/http"
	"os"

	"content-publisher/backend/internal/adapter"
	"content-publisher/backend/internal/service"
	"content-publisher/backend/internal/transport"
)

func main() {
	registry := adapter.NewRegistry(
		adapter.NewWeChatOfficialAccountAdapter(),
		adapter.NewZhihuAdapter(),
		adapter.NewBilibiliAdapter(),
		adapter.NewXiaohongshuAdapter(),
	)

	contentService := service.NewContentService(registry)
	handler := transport.NewHTTPHandler(contentService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler.Routes(),
	}

	log.Printf("content publisher backend listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
