package main

import (
	"consumer/config"
	"log"
	"time"

	"github.com/andryansyhh/auth-service/pkg/domain/dto"
	"github.com/andryansyhh/auth-service/pkg/middleware"
	"github.com/andryansyhh/auth-service/pkg/repository"
	"github.com/andryansyhh/auth-service/pkg/usecase"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Fatal error (panic): %v", r)
		}
	}()

	cfg := config.Load()

	db, err := sqlx.Connect("postgres", cfg.DBDSN)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)
	log.Println("Consumer successfully connected to its database!")

	jwtManager := middleware.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration)
	authRepo := repository.NewRepository(db)
	authUsecase := usecase.NewUserUsecase(authRepo, jwtManager)

	err = authUsecase.Register(dto.RegisterRequest{
		Username: "test_from_consumer",
		Password: "password",
	})

	if err != nil {
		log.Printf("Error during registration: %v", err)
	} else {
		log.Print("Success register user via consumer")
	}
}
