package nats

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"WB/internal/infrastructure/nats/dto"
	"WB/internal/order/service"
	stan "github.com/nats-io/stan.go"
)

type NatsHandler struct {
	sc      stan.Conn
	service service.OrderService
}

func NewNatsHandler(sc stan.Conn, service service.OrderService) *NatsHandler {
	return &NatsHandler{sc: sc, service: service}
}

func (nh *NatsHandler) Subscribe(channel string) error {
	_, err := nh.sc.Subscribe(channel, func(msg *stan.Msg) {
		go nh.handleMessage(msg)
	}, stan.DurableName("order-durable"))
	return err
}

func (nh *NatsHandler) handleMessage(msg *stan.Msg) {
	var request dto.HttpRequest
	if err := json.Unmarshal(msg.Data, &request); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	req, err := http.NewRequest(request.Method, request.URL, bytes.NewBuffer([]byte(request.Body)))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing HTTP request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Response status: %s, body: %s", resp.Status, string(body))
}

func (nh *NatsHandler) Close() {
	nh.sc.Close()
}
