package utils

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
	"service-datacollector/models"
	"strconv"
)

func publishToRabbitMQ(products []models.Product) {
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
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"queue_from_datacollector", // name
		true,                       // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	for _, product := range products {
		body, err := json.Marshal(product)
		if err != nil {
			log.Printf("Error encoding product: %s", err)
			continue
		}

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		if err != nil {
			log.Fatalf("Failed to publish a message: %s", err)
		}
	}
}

func FetchData() {
	url := "https://petstore.swagger.io/v2/pet/findByStatus?status=available"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching data: %s", err)
	}
	defer resp.Body.Close()

	var pets []models.PetResponse
	if err := json.NewDecoder(resp.Body).Decode(&pets); err != nil {
		log.Fatalf("Error decoding response: %s", err)
	}

	var products []models.Product
	for _, pet := range pets {
		product := models.Product{
			Name:        pet.Name,
			Description: pet.Status,
			Categories: []models.Category{{
				Name:        pet.Category.Name,
				Description: strconv.FormatInt(pet.Category.ID, 10),
			}},
		}
		products = append(products, product)
	}
	publishToRabbitMQ(products)
	log.Printf("Processed %d products", len(products))
}
