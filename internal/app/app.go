package app

import (
	"context"
	"log"

	"github.com/Vladimirmoscow84/Image_processor/internal/handlers"
	"github.com/Vladimirmoscow84/Image_processor/internal/queue_broker/kafka"
	"github.com/Vladimirmoscow84/Image_processor/internal/service"
	filestorage "github.com/Vladimirmoscow84/Image_processor/internal/storage/file_storage"
	"github.com/Vladimirmoscow84/Image_processor/internal/storage/postgres"
	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/ginext"
)

func Run() {
	cfg := config.New()
	err := cfg.LoadEnvFiles(".env")
	if err != nil {
		log.Fatalf("[app] error of loading cfg: %v", err)
	}
	cfg.EnableEnv("")

	databaseURI := cfg.GetString("DATABASE_URI")

	serverAddr := cfg.GetString("SERVER_ADDRESS")

	fileStorageRoot := cfg.GetString("FILE_STORAGE_ROOT")
	waterMarkPath := cfg.GetString("WATERMARK_PATH")

	kafkaBroker := cfg.GetString("KAFKA_BROKER")
	kafkaTopic := cfg.GetString("KAFKA_TOPIC")
	kafkaGroup := cfg.GetString("KAFKA_GROUP")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresStore, err := postgres.New(databaseURI)
	if err != nil {
		log.Fatalf("[app]failed to connect to PG DB: %v", err)
	}

	fileStorage, err := filestorage.New(fileStorageRoot, waterMarkPath)
	if err != nil {
		log.Fatalf("[app] failed to open file storage: %v", err)
	}

	kafkaCfg := &kafka.Config{
		Brokers: []string{kafkaBroker},
		Topic:   kafkaTopic,
		GroupID: kafkaGroup,
	}

	kafkaClient, err := kafka.NewKafkaClient(kafkaCfg)
	if err != nil {
		log.Fatalf("[app] failed to init kafka client: %v", err)
	}

	imageService, err := service.New(postgresStore, fileStorage, kafkaClient)
	if err != nil {
		log.Fatalf("[app] service init error: %v", err)
	}

	go imageService.StartKafkaConsumer(ctx)

	engine := ginext.New("release")
	router := handlers.New(engine, imageService, imageService, imageService, imageService)
	router.Routes()

	log.Printf("[app] server started on %s", serverAddr)
	err = engine.Run(serverAddr)
	if err != nil {
		log.Fatalf("[app] server failed: %v", err)
	}

}
