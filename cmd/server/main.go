package main

import (
	// "fmt"
	// "os"
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/bizzysGitHub/Golang-Chat-Angular/internal/chat"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	// Root context that cancels on Ctrl+C / SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Create the chat manager and run it
	mgr := chat.NewManager()
	go mgr.Run(ctx)
	mux := http.NewServeMux()

	// Health root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	// WebSocket endpoint
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWS(mgr, w, r)
	})

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server
	go func() {
		log.Println("HTTP listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until signal
	<-ctx.Done()
	log.Println("shutting down...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx) // stop accepting new conns; let in-flight finish
	log.Println("bye")
}

/*
func checkForJwtSig() {
	// secret := os.Getenv("JWT_HS256_SECRET")
	// if secret == "" {
	// 	fmt.Println("❌ No secret loaded. Did you create .env and run go get github.com/joho/godotenv ?")
	// } else {
	// 	fmt.Println("✅ Secret loaded from .env:", secret[:8], "...") // just show first few chars
	// }
}
*/
