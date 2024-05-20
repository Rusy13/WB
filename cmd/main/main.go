package main

import (
	"context"
	"log"
	"net/http"
	//"os"

	"WB/internal/infrastructure/database/postgres/database"
	"WB/internal/infrastructure/database/redis"

	"WB/internal/middleware"
	"WB/internal/order/delivery"
	serviceOrder "WB/internal/order/service"
	storageOrder "WB/internal/order/storage/database"
	"WB/internal/routes"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
		log.Println("Error is-----------------------", err)
	}

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error in logger initialization: %v", err)
	}
	logger := zapLogger.Sugar()
	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Printf("error in logger sync: %v", err)
		}
	}()
	//err = godotenv.Load(".env.testing")
	//if err != nil {
	//	logger.Fatalf("error in getting env: %s", err)
	//}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := database.New(ctx)
	if err != nil {
		logger.Fatalf("error in database init: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			logger.Errorf("error in closing db")
		}
	}()

	redisConn, err := redis.Init()
	if err != nil {
		logger.Fatalf("error on connection to redis: %v", err)
	}
	defer func() {
		err = redisConn.Close()
		if err != nil {
			logger.Infof("error on redis close: %s", err.Error())
		}
	}()

	stOrder := storageOrder.New(db, redisConn, logger)
	svOrder := serviceOrder.New(stOrder)
	d := delivery.New(svOrder, logger)

	mw := middleware.New(logger)
	router := routes.GetRouter(d, mw)

	//port := os.Getenv("APP_PORT")
	port := "8000"
	addr := ":" + port
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	logger.Fatal(http.ListenAndServe(addr, router))
}
