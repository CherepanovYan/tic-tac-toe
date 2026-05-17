package main

import (
	"fmt"
	"log"
	"net/http"
	"tic-tac-toe/internal/config"
	"time"

	"tic-tac-toe/internal/domain"
	"tic-tac-toe/internal/handler"
	"tic-tac-toe/internal/handler/game"
	"tic-tac-toe/internal/repository"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	repo := repository.NewRepository()
	domain := domain.NewDomain(repo)
	gameHandler := game.NewGameHendler(domain, domain)

	mux := http.NewServeMux()
	handler.RegisterRouter(mux, gameHandler)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
		// ReadTimeout — максимальное время на чтение всего запроса (заголовки + тело)
		ReadTimeout: 10 * time.Second,
		// WriteTimeout — максимальное время на запись ответа клиенту
		WriteTimeout: 10 * time.Second,
		// IdleTimeout — максимальное время ожидания следующего запроса при keep-alive соединении
		IdleTimeout: 60 * time.Second,
	}
	fmt.Println("Server started at http://localhost:" + cfg.Port)
	log.Fatal(server.ListenAndServe())
}
