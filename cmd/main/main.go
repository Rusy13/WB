package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"WB/internal/infrastructure/database/postgres/database"
	"WB/internal/infrastructure/database/redis"
	"WB/internal/infrastructure/kafka"
	"WB/internal/infrastructure/kafka/consumer"
	"WB/internal/middleware"
	"WB/internal/order/delivery"
	serviceOrder "WB/internal/order/service"
	storageOrder "WB/internal/order/storage/database"
	"WB/internal/routes"
	"go.uber.org/zap"
)

func main() {
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

	cfg, err := kafka.NewConfig()
	if err != nil {
		logger.Fatalf("error in kafka config init: %v", err)
	}
	defer func() {
		err = cfg.Close()
		if err != nil {
			logger.Errorf("error in closing sync kafka producer: %v", err)
		}
	}()

	stOrder := storageOrder.New(db, redisConn, logger)
	svOrder := serviceOrder.New(stOrder, cfg.Producer)
	d := delivery.New(svOrder, logger)

	mw := middleware.New(logger)
	router := routes.GetRouter(d, mw)

	//кафка
	go func() {
		err = consumer.Run(ctx, cfg, stOrder, logger)
		if err != nil {
			logger.Errorf("error in consumer running")
		}
	}()

	port := os.Getenv("APP_PORT")
	addr := ":" + port
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	logger.Fatal(http.ListenAndServe(addr, router))
}
