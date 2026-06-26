package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

// Estructura para decodificar el JSON que viene de RabbitMQ
type Quiniela struct {
	HomeTeam  string `json:"home_team"`
	AwayTeam  string `json:"away_team"`
	HomeGoals int    `json:"home_goals"`
	AwayGoals int    `json:"away_goals"`
	Username  string `json:"username"`
	Timestamp string `json:"timestamp"`
}

func main() {
	ctx := context.Background()

	// 1. Conectar a Valkey (Redis)
	valkeyAddr := os.Getenv("VALKEY_ADDR")
	if valkeyAddr == "" {
		valkeyAddr = "localhost:6379"
	}
	valkeyClient := redis.NewClient(&redis.Options{
		Addr:     valkeyAddr,
		Password: "",
		DB:       0,
	})

	// Probamos conexión a Valkey
	if err := valkeyClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Error conectando a Valkey: %v", err)
	}
	fmt.Println("Conectado a Valkey exitosamente")

	// 2. Conectar a RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://admin:secret@localhost:5672/"
	}
	conn, err := amqp.Dial(rabbitURL)
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error abriendo canal RabbitMQ: %v", err)
	}
	defer ch.Close()

	// 3. Empezar a consumir los mensajes de la cola "quinielas"
	mensajes, err := ch.Consume(
		"quinielas", // nombre de la cola
		"",          // consumer
		true,        // auto-ack (le avisa a RabbitMQ que ya lo leímos para que lo borre)
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatalf("Error registrando el consumidor: %v", err)
	}

	fmt.Println("Esperando mensajes de RabbitMQ...")

	// 4. Bucle infinito para leer mensajes
	for msg := range mensajes {
		var q Quiniela
		if err := json.Unmarshal(msg.Body, &q); err != nil {
			log.Printf("Error decodificando JSON: %v", err)
			continue
		}

		fmt.Printf("Procesando quiniela: %s vs %s\n", q.HomeTeam, q.AwayTeam)

		// --- GUARDAR EN VALKEY ---

		// Guardamos la quiniela en un "Stream" de Redis (ideal para series de tiempo en Grafana)
		// Convertimos todo a un formato llave-valor para Redis
		valores := map[string]interface{}{
			"home_team":  q.HomeTeam,
			"away_team":  q.AwayTeam,
			"home_goals": q.HomeGoals,
			"away_goals": q.AwayGoals,
			"username":   q.Username,
			"timestamp":  q.Timestamp,
		}

		err = valkeyClient.XAdd(ctx, &redis.XAddArgs{
			Stream: "predicciones_stream", // Nombre del stream general
			Values: valores,
		}).Err()

		if err != nil {
			log.Printf("Error guardando en Valkey: %v", err)
		}

		// Si el equipo involucrado es GTM (Tu equipo asignado por el carnet 0),
		// aumentamos un contador específico
		if q.HomeTeam == "GTM" || q.AwayTeam == "GTM" {
			valkeyClient.Incr(ctx, "total_predicciones_GTM")
		}
	}
}
