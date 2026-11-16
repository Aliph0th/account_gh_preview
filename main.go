package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	botService "tokens-bot/internal/bot"
	"tokens-bot/internal/database"
	"tokens-bot/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set in .env file")
	}

	db, err := database.New("/app/data/database.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	b, err := bot.New(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	botHandler := botService.New(db)
	botHandler.RegisterHandlers(b)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	httpServer := server.New(db)
	httpServer.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	wg.Go(func() {
		log.Println("Bot started polling...")
		b.Start(ctx)
		log.Println("Bot stopped")
	})

	wg.Go(func() {
		log.Println("🌐 HTTP server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
		}
		log.Println("HTTP server stopped")
	})

	<-ctx.Done()
	fmt.Println("\n⏹️  Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	}

	wg.Wait()

	fmt.Println("✅ Shutdown complete")
}
