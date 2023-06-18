package main

import (
	"aws-sqs/utils"
	"fmt"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("config/env/.env"); err != nil {
		fmt.Println("Error loading .env file")
	}

	// utils.SendMessageToQueue("Msg 1")
	// utils.SendMessageToQueue("Msg 2")
	// utils.SendMessageToQueue("Msg 3")
	// utils.SendMessageToQueue(utils.MockIocMessage())

	// utils.ReceiveMessageIoc()

	utils.ReceiveMessageIoc()
	for i := 0; i < 3; i++ {
		utils.ReceiveMessageIoc()
		time.Sleep(3 * time.Second)
	}
}
