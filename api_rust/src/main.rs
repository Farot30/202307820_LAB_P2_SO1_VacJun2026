use actix_web::{post, web, App, HttpResponse, HttpServer, Responder};
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

// TU CARNET AQUÍ
#[post("/grpc-202307820")]
async fn recibir_prediccion(datos: web::Json<Quiniela>) -> impl Responder {

    let go_url = env::var("GO_SERVER_URL").unwrap_or_else(|_| "http://127.0.0.1:50051".to_string());
    // 1. Conectarnos al servidor Go
    let mut client = match QuinielaServiceClient::connect(go_url).await {
        Ok(c) => c,
        Err(e) => {
            eprintln!("Error conectando a Go: {}", e);
            return HttpResponse::InternalServerError().finish();
        }
    };

    // 2. Armar el paquete gRPC
    let request = tonic::Request::new(QuinielaRequest {
        home_team: datos.home_team.clone(),
        away_team: datos.away_team.clone(),
        home_goals: datos.home_goals as i32,
        away_goals: datos.away_goals as i32,
        username: datos.username.clone(),
        timestamp: datos.timestamp.clone(),
    });

    // 3. Enviar a Go y esperar respuesta
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

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let puerto = 8080;
    println!("API Rust corriendo en http://0.0.0.0:{}", puerto);

    HttpServer::new(|| {
        App::new().service(recibir_prediccion)
    })
    .bind(("0.0.0.0", puerto))?
    .run()
    .await
}