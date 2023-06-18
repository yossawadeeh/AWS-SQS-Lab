package utils

import (
	"aws-sqs/constant"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type IocMessageResponse struct {
	Type                        string   `json:"type"`
	SpecVersion                 string   `json:"spec_version"`
	ID                          string   `json:"id"`
	Created                     string   `json:"created"`
	Modified                    string   `json:"modified"`
	Name                        string   `json:"name"`
	Description                 string   `json:"description"`
	Pattern                     string   `json:"pattern"`
	PatternType                 string   `json:"pattern_type"`
	ValidFrom                   string   `json:"valid_from"`
	Confidence                  uint     `json:"confidence"`
	Labels                      []string `json:"labels"`
	GovernmentScore             uint     `json:"x_government_score"`
	NationalSecurityAgencyScore uint     `json:"x_national_security_agency_score"`
	HealthCareScore             uint     `json:"x_health_care_score"`
	EnergyScore                 uint     `json:"x_energy_score"`
	TelecomScore                uint     `json:"x_telecom_score"`
	TransportationScore         uint     `json:"x_transportation_score"`
	FinancialScore              uint     `json:"x_financial_score"`
	OtherEnterpriseScore        uint     `json:"x_other_enterprise_score"`
}

func CreateSession() (ses *session.Session, err error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv(constant.AwsSqsRegion)), // AWS region
		Credentials: credentials.NewStaticCredentials(
			os.Getenv(constant.AwsSqsAccessKey),       // access key
			os.Getenv(constant.AwsSqsSecretAccessKey), // secret access key
			"",
		),
	})
	return sess, err
}

func SendMessageToQueue(messageBody string) error {
	fmt.Println("✨")

	// Create a new AWS session
	sess, err := CreateSession()
	if err != nil {
		return fmt.Errorf("❌ Failed to create session: %s", err)
	}

	// Create an SQS service client
	svc := sqs.New(sess)

	// Specify the queue URL
	queueURL := os.Getenv(constant.AwsSqsQueueUrl)

	// Send the message to the queue
	_, err = svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: &messageBody,
	})
	if err != nil {
		return fmt.Errorf("❌ Failed to send message to queue: %s", err)
	}

	fmt.Println("🧪 Send message to queue: ", messageBody)
	return nil
}

func ReceiveMessageIoc() {
	fmt.Println("✨")

	// Create a new AWS session
	sess, err := CreateSession()
	if err != nil {
		fmt.Println("❌ Failed to create session:", err)
		return
	}

	// Create an SQS service client
	svc := sqs.New(sess)

	// Specify the queue URL
	queueURL := os.Getenv(constant.AwsSqsQueueUrl)

	// Receive messages from the queue
	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &queueURL,
		MaxNumberOfMessages: aws.Int64(10), // Maximum number of messages to receive (adjust as needed)
		VisibilityTimeout:   aws.Int64(30), // Visibility timeout for the received messages (adjust as needed)
		WaitTimeSeconds:     aws.Int64(20), // Time to wait message
	})
	if err != nil {
		fmt.Println("❌ Failed to receive messages from queue:", err)
		return
	}

	fmt.Println("✅ result: ", result)
	fmt.Println("✅ count: ", len(result.Messages))
	// Process received messages
	for _, message := range result.Messages {
		// Access the message body
		body := aws.StringValue(message.Body)
		fmt.Println("✅ Received message:", body)

		// Remove outer double quotes
		if len(body) > 1 && body[0] == '"' && body[len(body)-1] == '"' {
			body = body[1 : len(body)-1]
		}
		body = strings.Replace(body, "\"\"", "\"", -1)
		body = strings.ReplaceAll(body, "\n", "")
		body = strings.ReplaceAll(body, "\r", "")
		fmt.Println("✅ fixed body message:", body)

		// save message to database
		jsonBytes := []byte(body)
		var iocData IocMessageResponse
		err := json.Unmarshal(jsonBytes, &iocData)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("✨ ", iocData)
		fmt.Println("✨ ", iocData.Name)

		// Delete the processed message from the queue
		_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      &queueURL,
			ReceiptHandle: message.ReceiptHandle,
		})
		if err != nil {
			fmt.Println("❌ Failed to delete message:", err)
			// Handle error if required
		}
		fmt.Println("🧪 Delete: ", message.ReceiptHandle)
	}
}

func MockIocMessage() (res string) {
	file, err := os.Open("data/mockIocData.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(content))
	return string(content)
}
