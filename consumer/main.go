package main

import (
	"consumer/config"
	"log"

	"github.com/andryansyhh/auth-service/pkg/domain/dto"
	"github.com/andryansyhh/auth-service/pkg/middleware"
	"github.com/andryansyhh/auth-service/pkg/repository"
	"github.com/andryansyhh/auth-service/pkg/usecase"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	log.Printf("Consumer Service starting on port %s", cfg.AppPort)

	conn, err := sqlx.Connect("postgres", cfg.DBDSN)
	if err != nil {
		log.Fatalf("failed to connect to consumer db: %v", err)
	}
	defer conn.Close()
	log.Println("Consumer successfully connected to its database!")

	jwtManager := middleware.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration)
	authRepo := repository.NewRepository(conn)
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
