package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/say8hi/go-api-test/internal/models"
	"github.com/streadway/amqp"
)

func InitRabbitMQ() *amqp.Channel {
    amqpUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/",
      os.Getenv("RMQ_USER"),
      os.Getenv("RMQ_PASSWORD"),
      os.Getenv("RMQ_HOST"),
      os.Getenv("RMQ_PORT"),
    )
    conn, err := amqp.Dial(amqpUrl)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %s", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %s", err)
    }

    _, err = ch.QueueDeclare(
        "queue_from_datacollector", // name
        true,        // durable
        false,       // delete when unused
        false,       // exclusive
        false,       // no-wait
        nil,         // arguments
    )
    if err != nil {
        log.Fatalf("Failed to declare a queue: %s", err)
    }

    return ch
}

func ConsumeMessages(channel *amqp.Channel, queueName string) {
    msgs, err := channel.Consume(
        queueName, // queue
        "",        // consumer
        true,      // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
    if err != nil {
        log.Fatalf("Failed to register a consumer: %s", err)
    }

    forever := make(chan bool)

    go func() {
      for d := range msgs {
        log.Printf("Received a message from datacollector: %s", d.Body)
        
        var product models.Product
        if err := json.Unmarshal(d.Body, &product); err != nil {
            log.Printf("Error decoding JSON: %s", err)
            continue
        } 
      
      }
    }()

    <-forever
}
