package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

// [RO] Publicator de Evenimente (NATS JetStream)
//
// Această componentă este "Portavocea" sistemului.
// Ea strigă în rețea când se întâmplă ceva important (ex: "Am terminat de analizat un articol!").
// Alte servicii pot asculta aceste mesaje și pot reacționa.
type JetStreamEventPublisher struct {
	jetStreamContext nats.JetStreamContext
}

// [RO] Constructor Publicator
func NewJetStreamEventPublisher(url string) (*JetStreamEventPublisher, error) {
	// [RO] 1. Conectare la Magistrală
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("[RO] Eroare: Nu m-am putut conecta la NATS: %w", err)
	}

	// [RO] 2. Activare JetStream (Persistență)
	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("[RO] Eroare: Nu am putut obține contextul JetStream: %w", err)
	}

	// [RO] 3. Asigurare Stream (Canal de Comunicare)
	// Creăm canalul "TRUTHWEAVE_EVENTS" dacă nu există.
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "TRUTHWEAVE_EVENTS",
		Subjects: []string{"article.>"}, // Ascultă orice subiect care începe cu "article."
	})
	if err != nil {
		// Ignorăm eroarea dacă stream-ul există deja.
	}

	return &JetStreamEventPublisher{jetStreamContext: js}, nil
}

// [RO] Publică Eveniment
// Trimite un mesaj structurat către toți abonații.
func (publisher *JetStreamEventPublisher) PublishEvent(executionContext context.Context, eventPayload interface{}) error {
	data, err := json.Marshal(eventPayload)
	if err != nil {
		return fmt.Errorf("[RO] Eroare: Nu am putut serializa evenimentul: %w", err)
	}

	// [RO] Subiectul Mesajului
	// Definim clar despre ce e vorba.
	subject := "article.ingested"

	// [RO] Trimitere efectivă
	_, err = publisher.jetStreamContext.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("[RO] Eroare: Nu am putut publica în NATS: %w", err)
	}

	return nil
}
