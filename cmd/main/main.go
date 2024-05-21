package main

import (
	"WB/internal/infrastructure/nats"
	"context"
	"github.com/nats-io/stan.go"
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
	err = stOrder.RestoreCacheFromDB() // Вызов функции восстановления кэша из БД
	if err != nil {
		logger.Fatalf("error restoring cache from DB: %v", err)
	}
	svOrder := serviceOrder.New(stOrder)
	d := delivery.New(svOrder, logger)

	mw := middleware.New(logger)
	router := routes.GetRouter(d, mw)

	//////////////////////////////////////////////////NATS//////////////////////////////////////////////
	natsConn, err := stan.Connect("test-cluster", "order-service", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		logger.Fatalf("error connecting to NATS Streaming: %v", err)
	}
	defer natsConn.Close()

	natsHandler := nats.NewNatsHandler(natsConn, svOrder)
	err = natsHandler.Subscribe("order-channel")
	if err != nil {
		logger.Fatalf("error subscribing to NATS channel: %v", err)
	}
	defer natsHandler.Close()
	//////////////////////////////////////////////////NATS//////////////////////////////////////////////
	//port := os.Getenv("APP_PORT")
	port := "8000"
	addr := ":" + port
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	logger.Fatal(http.ListenAndServe(addr, router))
}
