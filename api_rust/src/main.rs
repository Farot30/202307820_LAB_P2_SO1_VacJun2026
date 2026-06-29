use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder};
use serde::{Deserialize, Serialize};

// Importamos el código gRPC auto-generado
pub mod quinielapb {
    tonic::include_proto!("quinielapb"); 
}
use quinielapb::quiniela_service_client::QuinielaServiceClient;
use quinielapb::QuinielaRequest;
use std::env;

#[derive(Debug, Serialize, Deserialize)]
struct Quiniela {
    home_team: String,
    away_team: String,
    home_goals: u32,
    away_goals: u32,
    username: String,
    timestamp: String,
}

#[get("/")]
async fn health_check() -> impl Responder {
    HttpResponse::Ok().body("API Rust funcionando correctamente")
}

#[post("/grpc-202307820")]
async fn recibir_prediccion(datos: web::Json<Quiniela>) -> impl Responder {
    let go_url = env::var("GO_SERVER_URL").unwrap_or_else(|_| "http://127.0.0.1:50051".to_string());
    
    let mut client = match QuinielaServiceClient::connect(go_url).await {
        Ok(c) => c,
        Err(e) => {
            eprintln!("Error conectando a Go: {}", e);
            return HttpResponse::InternalServerError().finish();
        }
    };

    let request = tonic::Request::new(QuinielaRequest {
        home_team: datos.home_team.clone(),
        away_team: datos.away_team.clone(),
        home_goals: datos.home_goals as i32,
        away_goals: datos.away_goals as i32,
        username: datos.username.clone(),
        timestamp: datos.timestamp.clone(),
    });

    match client.enviar_quiniela(request).await {
        Ok(response) => {
            println!("Respuesta de Go: {}", response.into_inner().message);
            HttpResponse::Ok().json(serde_json::json!({"status": "success"}))
        },
        Err(e) => {
            eprintln!("Error en gRPC: {}", e);
            HttpResponse::InternalServerError().finish()
        }
    }
}

#[post("/dapr-202307820")]
async fn recibir_prediccion_dapr(datos: web::Json<Quiniela>) -> impl Responder {
    // "rabbitmq-pubsub" es el nombre del componente y "quinielas" es la cola
    let dapr_url = "http://localhost:3500/v1.0/publish/rabbitmq-pubsub/quinielas?metadata.rawPayload=true";
    
    let client = reqwest::Client::new();
    
    // Hacemos el POST directo al sidecar de Dapr
    match client.post(dapr_url).json(&datos.into_inner()).send().await {
        Ok(response) => {
            if response.status().is_success() {
                println!("Mensaje enviado exitosamente mediante Dapr a RabbitMQ");
                HttpResponse::Ok().json(serde_json::json!({"status": "success", "method": "dapr"}))
            } else {
                eprintln!("Dapr devolvió un error: {:?}", response.status());
                HttpResponse::InternalServerError().finish()
            }
        },
        Err(e) => {
            eprintln!("Error conectando con Dapr: {}", e);
            HttpResponse::InternalServerError().finish()
        }
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let puerto = 8080;
    println!("API Rust corriendo en http://0.0.0.0:{}", puerto);

    HttpServer::new(|| {
        App::new()
            .service(health_check)
            .service(recibir_prediccion)
            .service(recibir_prediccion_dapr)
    })
    .bind(("0.0.0.0", puerto))?
    .run()
    .await
}