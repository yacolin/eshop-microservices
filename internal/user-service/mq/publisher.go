package mq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"eshop-microservices/pkg/mq"
)

const sourceService = "user-service"

type Publisher struct {
	client *mq.Client
}

func NewPublisher(client *mq.Client) *Publisher {
	return &Publisher{client: client}
}

func (p *Publisher) PublishUserCreated(id, username, email string) {
	evt := mq.UserCreatedEvent{
		ID:       id,
		Username: username,
		Email:    email,
	}
	body := mq.Event{
		Type:      "user.created",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "user.created", body); err != nil {
		log.Printf("publish user.created: %v", err)
	}
}

func (p *Publisher) PublishUserUpdated(id, username, email string) {
	evt := mq.UserUpdatedEvent{
		ID:       id,
		Username: username,
		Email:    email,
	}
	body := mq.Event{
		Type:      "user.updated",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "user.updated", body); err != nil {
		log.Printf("publish user.updated: %v", err)
	}
}

func (p *Publisher) PublishUserDeleted(id string) {
	evt := mq.UserDeletedEvent{ID: id}
	body := mq.Event{
		Type:      "user.deleted",
		Data:      mustMarshal(evt),
		Timestamp: time.Now().Format(time.RFC3339),
		Source:    sourceService,
	}
	if err := p.client.Publish(context.Background(), "user.deleted", body); err != nil {
		log.Printf("publish user.deleted: %v", err)
	}
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
