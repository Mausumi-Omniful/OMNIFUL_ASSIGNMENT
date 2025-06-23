package main

import (
	"context"
	"fmt"
	"os"
	"time"
    "oms/utils"
	"oms/controllers"
	"oms/database"
	"oms/routes"
	"github.com/joho/godotenv"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
)

func main() {
	fmt.Println("OMS Service Starting...")

	_ = godotenv.Load()
	_ = i18n.Initialize(i18n.WithRootPath("./localization"))


//    config
	bucketName := getEnvOrDefault("S3_BUCKET_NAME", "order-csv-bucket")
	s3Endpoint := getEnvOrDefault("AWS_S3_ENDPOINT", "http://localhost:4566")
	sqsEndpoint := getEnvOrDefault("AWS_SQS_ENDPOINT", "http://localhost:4566")
	awsRegion := getEnvOrDefault("AWS_REGION", "us-east-1")
	queueName := getEnvOrDefault("CREATE_BULK_ORDER_QUEUE_NAME", "CreateBulkOrder")

	kafkaBrokers := []string{getEnvOrDefault("KAFKA_BROKERS", "localhost:9092")}
	kafkaTopic := getEnvOrDefault("KAFKA_ORDER_TOPIC", "order.created")

	mongoURI := getEnvOrDefault("MONGODB_URI", "mongodb://myuser:mypassword@localhost:27018/mydb?authSource=admin")
	mongoDBName := getEnvOrDefault("MONGODB_DB_NAME", "mydb")




	// mongodb
	mongoDB, err := database.NewDatabase(context.Background(), mongoURI, mongoDBName)
	if err != nil {
		fmt.Printf("Mongo connect error: %v\n", err)
		return
	}

	orderRepo := database.NewOrderRepository(mongoDB)



	// imsurl
	imsBaseURL := getEnvOrDefault("IMS_BASE_URL", "http://localhost:8084")
	imsClient := utils.NewIMSClient(imsBaseURL)




	// s3upload
	s3Uploader, err := utils.NewS3Uploader(bucketName, s3Endpoint, awsRegion)
	if err != nil {
		fmt.Printf("S3 upload init error: %v\n", err)
		return
	}


	// sqspublisher
	sqsPublisher, err := utils.NewSQSPublisher(queueName, sqsEndpoint, awsRegion)
	if err != nil {
		fmt.Printf("SQS publish init error: %v\n", err)
		return
	}



	// kafkaproducer
	kafkaProducer, err := utils.NewKafkaProducer(kafkaBrokers, kafkaTopic)
	if err != nil {
		fmt.Printf("Kafka producer init error: %v\n", err)
		return
	}

	kafkaAvailable := true
	var kafkaConsumer *utils.OrderFinalizationConsumer

	if kafkaAvailable {
		kafkaConsumer, err = utils.NewOrderFinalizationConsumer(kafkaBrokers, kafkaTopic, orderRepo, imsClient)
		if err != nil {
			fmt.Printf("Kafka consumer init error: %v\n", err)
			return
		}
	}

	defaultHandler, err := utils.NewDefaultMessageHandler(s3Endpoint, awsRegion, orderRepo, imsClient, kafkaProducer, s3Uploader)
	if err != nil {
		fmt.Printf("Message handler error: %v\n", err)
		return
	}
	sqsConsumer, err := utils.NewSQSConsumer(queueName, sqsEndpoint, awsRegion, defaultHandler)
	if err != nil {
		fmt.Printf("SQS consumer error: %v\n", err)
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



	// server
	server := http.InitializeServer(
		":8086",
		10*time.Second,
		10*time.Second,
		70*time.Second,
		true,
	)

	orderController := &controllers.OrderController{
		S3Uploader:   s3Uploader,
		SQSPublisher: sqsPublisher,
		OrderRepo:    orderRepo,
	}

	routes.RegisterOrderRoutes(server, orderController)

	fmt.Println("OMS Service Ready")

	if err := server.StartServer("oms-service"); err != nil {
		fmt.Printf("HTTP server start error: %v\n", err)
		cancel()
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
