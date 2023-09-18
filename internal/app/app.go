package app

import (
	"context"
	"fio/internal/config"
	delivery "fio/internal/delivery/http"
	"fio/internal/delivery/kafka"
	"fio/internal/repository"
	"fio/internal/server"
	"fio/internal/service"
	"fio/pkg/cache"
	"fio/pkg/database/postgres"
	"fio/pkg/database/redis"
	"fio/pkg/profiler"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	queue "fio/pkg/queue/kafka"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

const (
	agifyResource       = "https://api.agify.io"
	genderizeResource   = "https://api.genderize.io"
	nationalizeResource = "https://api.nationalize.io"
	timeout             = 5 * time.Second
)

// @title Effective-Mobile Trainee Assignment
// @version 2.0
// @description тех. задание с отбора на go разарботчика Effective-Mobile

// @host localhost:8079
// @BasePath /
func Run(configDir string) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error occurred while loading zapLogger: %s\n", err.Error())
		return
	}
	defer zapLogger.Sync() //nolint:errcheck
	logger := zapLogger.Sugar()

	cfg, err := config.InitConfig(configDir)
	if err != nil {
		logger.Errorf("Error occurred while loading config: %s\n", err.Error())
		return
	}

	db, err := postgres.NewPostgresqlDB(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username,
		cfg.Postgres.DBName, cfg.Postgres.Password, cfg.Postgres.SSLMode)
	if err != nil {
		logger.Errorf("Error occurred while loading Postgres DB: %s\n", err.Error())
		return
	}

	rdb, err := redis.NewRedisClient(cfg.Redis.Address, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.Errorf("Error occurred while loading Redis DB: %s\n", err.Error())
		return
	}
	cache := cache.NewRedisCache(rdb)
	repos := repository.NewRepository(db)
	prof := profiler.NewNameProfiler(agifyResource, genderizeResource, nationalizeResource, timeout)
	services := service.NewService(service.Dependencies{
		Repos:        repos,
		Cache:        cache,
		CacheTTL:     cfg.CacheTTL,
		NameProfiler: prof,
	})

	validate := validator.New()

	httpHandler := delivery.NewHandler(services, validate, logger)

	mux := httpHandler.InitRoutes()

	consumerGroup, err := queue.InitKafkaConsumerGroup(cfg.KafkaEndpoints, cfg.Kafka.GroupID, cfg.Kafka.ClientID,
		cfg.Kafka.TLSEnable, cfg.Kafka.ReturnSucceses, sarama.RequiredAcks(cfg.Kafka.RequiredAcks))
	if err != nil {
		logger.Errorf("Error occurred while init Kafka Consumer Group: %s\n", err.Error())
		return
	}

	syncProducer, err := queue.InitKafkaSyncProducer(cfg.KafkaEndpoints, cfg.Kafka.ClientID, cfg.Kafka.TLSEnable,
		cfg.Kafka.ReturnSucceses, sarama.RequiredAcks(cfg.Kafka.RequiredAcks))
	if err != nil {
		logger.Errorf("Error occurred while init Kafka Sync Producer: %s\n", err.Error())
		return
	}

	messageHandler := kafka.NewMessageHander(services, validate, logger, syncProducer)

	consumeContext := context.Background()
	go messageHandler.ConsumeLoop(cfg.KafkaTopics, consumeContext, consumerGroup)
	defer consumeContext.Done()

	srv := server.NewServer(cfg, mux)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err = srv.Run(); err != nil {
			logger.Errorf("Failed to start server: %s\n", err.Error())
		}
	}()

	logger.Info("Application is running")

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	logger.Info("Application is shutting down")

	if err = srv.Shutdown(ctx); err != nil {
		logger.Error(err.Error())
	}

	if err = db.Close(); err != nil {
		logger.Error(err.Error())
	}
}
