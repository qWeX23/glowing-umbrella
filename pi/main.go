package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shirou/gopsutil/host"
)

func main() {
	// Loop forever
	for {
		// Call the stubbed out function
		sendToRabbit()

		// Wait for five minutes
		time.Sleep(5 * time.Minute)
	}
}

func sendToRabbit() {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://username:password@192.168.1.16:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"myqueue", // Name of the queue
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Get the current CPU temperature
	cpuTemp, err := host.SensorsTemperatures()
	if err != nil {
		log.Fatalf("Failed to get CPU temperature: %v", err)
	}
	// create a sample object
	data := map[string]interface{}{
		"time":  time.Now().UTC().Unix(),
		"temps": cpuTemp,
	}

	// convert object to JSON string
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// print JSON string
	fmt.Println(string(jsonData))
	// Send the CPU temperature to the queue
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(string(jsonData)),
		},
	)
	if err != nil {
		log.Fatalf("Failed to send message to queue: %v", err)
	}

	fmt.Println("Message sent to queue")
}
