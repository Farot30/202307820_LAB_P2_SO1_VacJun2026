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

// Estructura para decodificar el JSON
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
	if err != nil {
		log.Fatalf("Error conectando a RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error abriendo canal RabbitMQ: %v", err)
	}
	defer ch.Close()

	mensajes, err := ch.Consume(
		"quinielas", "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Error registrando el consumidor: %v", err)
	}

	fmt.Println("Esperando mensajes de RabbitMQ...")

	// 3. Procesar mensajes
	for msg := range mensajes {
		var q Quiniela
		if err := json.Unmarshal(msg.Body, &q); err != nil {
			log.Printf("Error decodificando JSON: %v", err)
			continue
		}

		fmt.Printf("Procesando quiniela: %s vs %s\n", q.HomeTeam, q.AwayTeam)

		// A. Historial completo
		valores := map[string]interface{}{
			"home_team":  q.HomeTeam,
			"away_team":  q.AwayTeam,
			"home_goals": q.HomeGoals,
			"away_goals": q.AwayGoals,
			"username":   q.Username,
			"timestamp":  q.Timestamp,
		}
		valkeyClient.XAdd(ctx, &redis.XAddArgs{
			Stream: "predicciones_stream",
			Values: valores,
		})

		// B. Top de usuarios más activos
		valkeyClient.ZIncrBy(ctx, "top_usuarios", 1, q.Username)

		// C. Top de equipos con más victorias
		if q.HomeGoals > q.AwayGoals {
			valkeyClient.ZIncrBy(ctx, "top_victorias", 1, q.HomeTeam)
		} else if q.AwayGoals > q.HomeGoals {
			valkeyClient.ZIncrBy(ctx, "top_victorias", 1, q.AwayTeam)
		}

		// D. Moda de goles
		valkeyClient.ZIncrBy(ctx, "moda_goles_local", 1, fmt.Sprintf("%d", q.HomeGoals))
		valkeyClient.ZIncrBy(ctx, "moda_goles_visitante", 1, fmt.Sprintf("%d", q.AwayGoals))

		// E. Rangos (Mayor y Menor cantidad de goles)
		valkeyClient.ZAdd(ctx, "rango_goles_local", redis.Z{Score: float64(q.HomeGoals), Member: fmt.Sprintf("%d", q.HomeGoals)})
		valkeyClient.ZAdd(ctx, "rango_goles_visitante", redis.Z{Score: float64(q.AwayGoals), Member: fmt.Sprintf("%d", q.AwayGoals)})

		// F. Estadísticas específicas del equipo GTM (Guatemala)
		if q.HomeTeam == "GTM" {
			valkeyClient.Incr(ctx, "total_predicciones_GTM")
			valkeyClient.XAdd(ctx, &redis.XAddArgs{
				Stream: "stream_gtm_local",
				Values: map[string]interface{}{"goles": q.HomeGoals, "timestamp": q.Timestamp},
			})
		} else if q.AwayTeam == "GTM" {
			valkeyClient.Incr(ctx, "total_predicciones_GTM")
			valkeyClient.XAdd(ctx, &redis.XAddArgs{
				Stream: "stream_gtm_visitante",
				Values: map[string]interface{}{"goles": q.AwayGoals, "timestamp": q.Timestamp},
			})
		}
	}
}
