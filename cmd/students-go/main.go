package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MdSadiqMd/students-go/internal/config"
	"github.com/MdSadiqMd/students-go/internal/http/handlers/student"
	"github.com/MdSadiqMd/students-go/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("Failed to create storage: ", err) // No need to start server if database is not created
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	router.HandleFunc("POST /api/v1/student", student.New(storage))
	router.HandleFunc("GET /api/v1/student/{id}", student.GetById(storage))

	server := http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}

	slog.Info("Server is running", "address", cfg.HTTPServer.Address)

	done := make(chan os.Signal, 1)                                    // we are using a channel to handle the shutdown
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // Notify the done channel when the type of signals given gets triggered

	go func() { // we do this using do routines for graceful shutdown of the server
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	}()

	<-done

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal("Failed to shutdown server: ", err)
	}

	slog.Info("Server shut down")
}
