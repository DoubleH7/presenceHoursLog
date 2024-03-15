package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbLogger log.Logger = *getDBlogger()

func ConnectDB() (*mongo.Client, error) {
	logger := getDBlogger()

	// loading database info from .env file
	godotenv.Load("../.env")
	db_uri := os.Getenv("DB_URI")

	logger.Print("Establishing database connection...")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db_uri))
	if err != nil {
		logger.Fatal("Database connection failure :", err.Error())
		return new(mongo.Client), err
	}
	logger.Println(" done!")

	return client, nil
}

func DisconnectDB(client *mongo.Client) error {
	logger := getDBlogger()
	if err := client.Disconnect(context.Background()); err != nil {
		logger.Fatal("database closure failed:", err)
	}
	logger.Println("connection closed")
	return nil
}

func getDBlogger() *log.Logger {
	file, err := openLogFile("./database/dblog.log")
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(file, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)

	return logger
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}
