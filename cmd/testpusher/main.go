package main

import (
	"WB_Intern/internal/model"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"math/rand"
	"os"
)

const (
	ClusterID = "test-cluster"
	PublisherID = "pub"

	NumberOfMessages = 1

)

func main() {
	sc, err := stan.Connect(ClusterID, PublisherID, stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Fatalln(err)
	}
	jsonObject, err := os.ReadFile("../model.json")
	if err != nil {
		log.Fatalln(err)
	}
	var order model.Order
	if err = json.Unmarshal(jsonObject, &order); err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < NumberOfMessages; i++ {
		order.OrderUID = fmt.Sprintf("%xtest", rand.Int())

		data, err := json.Marshal(order)
		if err != nil {
			log.Fatalln(err)
		}

		if err := sc.Publish("order", data); err != nil {
			log.Fatalln(err)
		}

		log.Printf("published order with order_uid=%s\n", order.OrderUID)
	}
}