package main

import (
	"fmt"
	"log"
	"time"

	"github.com/andryansyhh/auth-service/cmd/config"
	"github.com/andryansyhh/auth-service/pkg/handler"
	"github.com/andryansyhh/auth-service/pkg/middleware"
	"github.com/andryansyhh/auth-service/pkg/repository"
	"github.com/andryansyhh/auth-service/pkg/usecase"
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

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

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

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.StructuredLogger())
	r.Use(middleware.ErrorHandler())

	handler.RegisterRoutes(r)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Println("Auth service running on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
