package main

import (
	"context"
	"os"
	"time"

	"oms/controllers"
	"oms/database"
	"oms/routes"
	"oms/utils"

	"github.com/joho/godotenv"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

func main() {
	log.Info("OMS Service Starting...")

	// Load environment variables (optional)
	_ = godotenv.Load()

	// Initialize i18n for internationalized error messages
	_ = i18n.Initialize(i18n.WithRootPath("./localization"))

	// Set default values for LocalStack if not provided
	bucketName := getEnvOrDefault("S3_BUCKET_NAME", "order-csv-bucket")
	s3Endpoint := getEnvOrDefault("AWS_S3_ENDPOINT", "http://localhost:4566")
	sqsEndpoint := getEnvOrDefault("AWS_SQS_ENDPOINT", "http://localhost:4566")
	awsRegion := getEnvOrDefault("AWS_REGION", "us-east-1")
	queueName := getEnvOrDefault("CREATE_BULK_ORDER_QUEUE_NAME", "CreateBulkOrder")

	// Kafka configuration
	kafkaBrokers := []string{getEnvOrDefault("KAFKA_BROKERS", "localhost:9092")}
	kafkaTopic := getEnvOrDefault("KAFKA_ORDER_TOPIC", "order.created")

	// MongoDB connection settings
	mongoURI := getEnvOrDefault("MONGODB_URI", "mongodb://myuser:mypassword@localhost:27018/mydb?authSource=admin")
	mongoDBName := getEnvOrDefault("MONGODB_DB_NAME", "mydb")

	// Initialize MongoDB connection
	mongoDB, err := database.NewDatabase(context.Background(), mongoURI, mongoDBName)
	if err != nil {
		log.WithError(err).Error("Failed to initialize MongoDB connection")
		return
	}

	// Initialize Order Repository
	orderRepo := database.NewOrderRepository(mongoDB)

	// Initialize IMS Client for validation
	imsBaseURL := getEnvOrDefault("IMS_BASE_URL", "http://localhost:8084")
	imsClient := utils.NewIMSClient(imsBaseURL)

	// Initialize S3 uploader
	s3Uploader, err := utils.NewS3Uploader(bucketName, s3Endpoint, awsRegion)
	if err != nil {
		log.WithError(err).Error("Failed to initialize S3 uploader")
		return
	}

	// Initialize SQS publisher
	sqsPublisher, err := utils.NewSQSPublisher(queueName, sqsEndpoint, awsRegion)
	if err != nil {
		log.WithError(err).Error("Failed to initialize SQS publisher")
		return
	}

	// Initialize Kafka producer
	kafkaProducer, err := utils.NewKafkaProducer(kafkaBrokers, kafkaTopic)
	if err != nil {
		log.WithError(err).Error("Failed to initialize Kafka producer")
		return
	}

	// Kafka consumer
	kafkaAvailable := true
	var kafkaConsumer *utils.OrderFinalizationConsumer

	if kafkaAvailable {
		kafkaConsumer, err = utils.NewOrderFinalizationConsumer(kafkaBrokers, kafkaTopic, orderRepo, imsClient)
		if err != nil {
			log.WithError(err).Error("Failed to initialize Kafka consumer")
			return
		}
	}

	// Initialize SQS consumer with default message handler
	defaultHandler, err := utils.NewDefaultMessageHandler(s3Endpoint, awsRegion, orderRepo, imsClient, kafkaProducer, s3Uploader)
	if err != nil {
		log.WithError(err).Error("Failed to initialize default message handler")
		return
	}
	sqsConsumer, err := utils.NewSQSConsumer(queueName, sqsEndpoint, awsRegion, defaultHandler)
	if err != nil {
		log.WithError(err).Error("Failed to initialize SQS consumer")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func() {
		_ = mongoDB.Close(context.Background())
		kafkaProducer.Close()
		if kafkaConsumer != nil {
			kafkaConsumer.Stop()
		}
	}()

	go func() {
		sqsConsumer.Start(ctx)
	}()

	if kafkaAvailable {
		go func() {
			kafkaConsumer.Start(ctx)
		}()
	}

	server := http.InitializeServer(
		":8086",
		10*time.Second,
		10*time.Second,
		70*time.Second,
		true, // enableRecovery
	)

	orderController := &controllers.OrderController{
		S3Uploader:   s3Uploader,
		SQSPublisher: sqsPublisher,
		OrderRepo:    orderRepo,
	}

	routes.RegisterOrderRoutes(server, orderController)

	log.Info("OMS Service initialized successfully")

	if err := server.StartServer("oms-service"); err != nil {
		log.WithError(err).Error("Failed to start OMS HTTP server")
		cancel()
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
