package main

import (
	"fmt"
	"log"

	"github.com/andryansyhh/auth-service/cmd/config"
	"github.com/andryansyhh/auth-service/internal/handler"
	"github.com/andryansyhh/auth-service/internal/middleware"
	"github.com/andryansyhh/auth-service/internal/repository"
	"github.com/andryansyhh/auth-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.Load()

	db, err := sqlx.Connect("postgres", cfg.DBDSN)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set goose dialect: %v", err)
	}
	if err := goose.Up(db.DB, "../db/migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Database migrated successfully")

	jwtManager := middleware.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration)
	repo := repository.NewRepository(db)
	uc := usecase.NewUserUsecase(repo, jwtManager)
	handler := handler.NewUserHandler(uc, jwtManager)

	r := gin.Default()
	handler.RegisterRoutes(r)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Println("Auth service running on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
