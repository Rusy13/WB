package main

import (
	"encoding/json"
	"log"

	"WB/internal/infrastructure/nats/dto"
	stan "github.com/nats-io/stan.go"
)

const (
	clusterID = "test-cluster"
	clientID  = "publisher"
	channel   = "order-channel"
	natsURL   = "nats://localhost:4222"
)

func main() {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
	}
	defer sc.Close()

	request := dto.HttpRequest{
		Method: "POST",
		URL:    "http://localhost:8000/order",
		Body:   `{"order_uid": "testtesttesttest4", "track_number": "WBILMTESTTRACK", "entry": "WBIL", "delivery_name": "Test Testov", "delivery_phone": "+9720000000", "delivery_zip": "2639809", "delivery_city": "Kiryat Mozkin", "delivery_address": "Ploshad Mira 15", "delivery_region": "Kraiot", "delivery_email": "test@gmail.com", "payment_transaction": "b563feb7b2b84b6test", "request_id": "", "currency": "USD", "provider": "wbpay", "amount": 1817, "payment_dt": 1637907727, "bank": "alpha", "delivery_cost": 1500, "goods_total": 317, "custom_fee": 1, "items": [{"chrt_id": 9934930, "track_number": "WBILMTESTTRACK", "price": 453, "rid": "ab4219087a764ae0btest", "name": "Mascaras", "sale": 30, "size": "0", "total_price": 317, "nm_id": 2389212, "brand": "Vivienne Sabo", "status": 202}], "locale": "en", "internal_signature": "", "customer_id": "test", "delivery_service": "meest", "shardkey": "9", "sm_id": 99, "date_created": "2021-11-26T06:22:19Z", "oof_shard": "1"}`,
	}

	data, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("Error marshaling request: %v", err)
	}

	if err := sc.Publish(channel, data); err != nil {
		log.Fatalf("Error publishing message: %v", err)
	}

	log.Println("Message published")
}
