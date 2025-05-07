package main

import (
	"log"

	"github.com/DataDog/datadog-go/statsd"
)

func main() {
	client, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		log.Fatal(err)
	}

	err = client.Count("amantya.custom.metric", 1, nil, 1)
	if err != nil {
		log.Fatal("Failed to send metric:", err)
	}

	log.Println("Metric sent to Datadog agent successfully!")
}
