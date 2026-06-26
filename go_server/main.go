package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"

	// Importamos la carpeta pb que generamos en el paso anterior
	pb "go_server/pb"
)

// Estructura de nuestro servidor
type server struct {
	pb.UnimplementedQuinielaServiceServer
	rabbitChannel *amqp.Channel
}

// Implementación de la función que definimos en el .proto
func (s *server) EnviarQuiniela(ctx context.Context, req *pb.QuinielaRequest) (*pb.QuinielaResponse, error) {
	log.Printf("Recibido de Rust (gRPC): %s vs %s (Goles: %d-%d)",
		req.HomeTeam, req.AwayTeam, req.HomeGoals, req.AwayGoals)

	// Convertimos la quiniela a JSON para RabbitMQ
	cuerpoMensaje, err := json.Marshal(req)
	if err != nil {
		return &pb.QuinielaResponse{Success: false, Message: "Error interno"}, err
	}

	// Enviamos el mensaje a la cola "quinielas" en RabbitMQ
	err = s.rabbitChannel.PublishWithContext(ctx,
		"",          // exchange
		"quinielas", // nombre de la cola
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        cuerpoMensaje,
		})

	if err != nil {
		return &pb.QuinielaResponse{Success: false, Message: "Fallo envío a RabbitMQ"}, err
	}

	return &pb.QuinielaResponse{Success: true, Message: "Guardado en RabbitMQ exitosamente"}, nil
}

func main() {
	// <-- 2. CAMBIO: Leemos la URL desde Kubernetes, si no existe, usamos localhost
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://admin:secret@localhost:5672/"
	}

	// Conectar a RabbitMQ usando la variable dinámica
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Fallo al conectar a RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Fallo al abrir un canal: %v", err)
	}
	defer ch.Close()

	// Asegurarnos de que la cola exista
	_, err = ch.QueueDeclare("quinielas", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Fallo al declarar la cola: %v", err)
	}

	// 2. Levantar el servidor gRPC en el puerto 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Fallo al escuchar puerto: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterQuinielaServiceServer(s, &server{rabbitChannel: ch})

	log.Printf("Servidor Go escuchando peticiones gRPC en el puerto 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo en el servidor gRPC: %v", err)
	}
}
