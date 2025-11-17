package app

import (
	"log"

	"github.com/Vladimirmoscow84/Image_processor/internal/queue_broker/kafka"
	"github.com/Vladimirmoscow84/Image_processor/internal/service"
	filestorage "github.com/Vladimirmoscow84/Image_processor/internal/storage/file_storage"
	"github.com/Vladimirmoscow84/Image_processor/internal/storage/postgres"
	"github.com/wb-go/wbf/config"
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

	fileSorageRoot := cfg.GetString("FILE_STORAGE_ROOT")
	waterMarkPath := cfg.GetString("WATERMARK_PATH")

	kafkaBroker := cfg.GetString("KAFKA_BROKER")
	kafkaTopic := cfg.GetString("KAFKA_TOPIC")

	postgresStore, err := postgres.New(databaseURI)
	if err != nil {
		log.Fatalf("[app]failed to connect to PG DB: %v", err)
	}
	fileStorage, err := filestorage.New(fileSorageRoot, waterMarkPath)
	if err != nil {
		log.Fatal("[app] failed to open file storage: %v", err)
	}

	producer, err := kafka.NewProducer(kafkaBroker, kafkaTopic)
	if err != nil {
		log.Fatalf("[app] kafka producer failed: %v", err)
	}

	imageService, err := service.New(postgresStore, fileStorage, producer)

}
