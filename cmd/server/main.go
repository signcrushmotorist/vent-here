package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/signcrushmotorist/vent-here/internal/application/auth"
	"github.com/signcrushmotorist/vent-here/internal/application/user"
	"github.com/signcrushmotorist/vent-here/internal/infrastructure/database"
	"github.com/signcrushmotorist/vent-here/internal/infrastructure/persistence"
	"github.com/signcrushmotorist/vent-here/internal/interfaces/http/handlers"
	"github.com/signcrushmotorist/vent-here/internal/interfaces/http/middleware"
)

func main() {

	db := database.Connect("localhost", "5432", "postgres", "wildan", "anonconfess")
	defer db.Close()
	mux := http.NewServeMux()

	userRepo := persistence.NewUserRepository(db)
	authService := auth.NewAuthServie(userRepo)
	userService := user.NewUserService(userRepo)
	authHandler := handlers.NewAuthHandler(authService, userService)
	sessionStore := middleware.NewSessionStore()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server is running"))
	})
	mux.HandleFunc("/register", authHandler.RegisterHandler)
	mux.HandleFunc("/login", authHandler.LoginHandler(sessionStore))
	mux.HandleFunc("/change-alias", middleware.RequireAuth(sessionStore, authHandler.ChangeAliasHandler))
	mux.HandleFunc("/protected", middleware.RequireAuth(sessionStore, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("you are authenticated"))
	}),
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down server")
}
