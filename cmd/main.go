package main

import (
	"context"
	"os"
	"os/signal"
	"terminer"
	"terminer/internal/handler"
	"terminer/internal/repository"
	"terminer/internal/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// TODO add logger to all services and repos
	// TODO add statistics service
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})
	if err := initConfig(); err != nil {
		logger.Fatalf("Failed to init config: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error loading .env file")
	}
	db, err := repository.NewPostgres(&repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.database"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logger.Fatalf("Failed to init db: %s", err.Error())
	}
	repos := repository.NewRepository(db, logger)
	services := service.NewService(repos, logger)
	handlers := handler.NewHandler(services)

	server := new(terminer.Server)
	
	go func ()  {
		if err := server.Run(viper.GetString("server.port"), handlers.InitRoutes()); err != nil {
			logger.Fatalf("Failed to start server: %s", err.Error())
		}
	} ()
	logger.Print("Server started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Fatalf("Failed to shutdown server: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logger.Fatalf("Failed to close db: %s", err.Error())
	}
	
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
