package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

type JWT struct {
	JWT JWTVal `json:"content"`
}

type JWTVal struct {
	Link string `json:"authToken"`
}

func init() {
	//InfoLogger.Println("Starting the application...")
	//InfoLogger.Println("Something noteworthy happened")
	//WarningLogger.Println("There is something you should know about")
	//ErrorLogger.Println("Something went wrong")

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func goDotEnvVariable(key string) string {

	// load .env file
	errLoad := godotenv.Load(".env")

	if errLoad != nil {
		ErrorLogger.Println("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	ourSign, err := return_sign()
	if err != nil {
		ErrorLogger.Println(err)
	} else {
		ourToken, err := return_token(ourSign)
		if err != nil {
			ErrorLogger.Println(err)
		} else {
			getPrice(ourToken)
			getDescr(ourToken)
		}
	}

}
