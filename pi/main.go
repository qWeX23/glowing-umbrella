package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shirou/gopsutil/host"
)

func main() {
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

	// Send the CPU temperature to the queue
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("Current CPU temperature: %v\n", cpuTemp[0].Temperature)),
		},
	)
	if err != nil {
		log.Fatalf("Failed to send message to queue: %v", err)
	}

	fmt.Println("Message sent to queue")
}
