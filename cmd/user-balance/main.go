package main

import (
	"log/slog"
	"net/http"
	"os"

	_ "github.com/irbgeo/user-balance/docs"
	"github.com/irbgeo/user-balance/internal/api"
	"github.com/irbgeo/user-balance/internal/service"
	"github.com/irbgeo/user-balance/internal/storage"
	"github.com/irbgeo/user-balance/internal/storage/driver"
)

// @title			Balance Service API
// @description	This is an example API for managing user balances.
func main() {
	slog.Info("configuration")

	port := os.Getenv("PORT")
	if port == "" {
		slog.Error("not found server port")
		return
	}

	redisAddress := os.Getenv("REDIS_ADDRESS")
	if port == "" {
		slog.Error("not found redis address")
		return
	}

	redisUser := os.Getenv("REDIS_USER")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDriver, err := driver.Redis(driver.RedisOpts{
		Address:  redisAddress,
		Username: redisUser,
		Password: redisPassword,
	})
	if err != nil {
		slog.Error("failed connect to redis: %w", err)
		return
	}

	str := storage.New(redisDriver)
	svc := service.New(str)

	api.Routes(svc)

	slog.Info("Turn on")

	if err := http.ListenAndServe(port, nil); err != http.ErrServerClosed {
		slog.Error("server: %w", err)
		return
	}

	slog.Info("Goodbye")
}
