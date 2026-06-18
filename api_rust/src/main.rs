use actix_web::{post, web, App, HttpResponse, HttpServer, Responder};
use serde::{Deserialize, Serialize};

// 1. Definimos la estructura de la quiniela basada en el JSON del proyecto
#[derive(Debug, Serialize, Deserialize)]
struct Quiniela {
    home_team: String,
    away_team: String,
    home_goals: u32,
    away_goals: u32,
    username: String,
    timestamp: String,
}

// 2. Definimos la ruta exigida. ¡CAMBIA EL NÚMERO POR TU CARNET!
#[post("/grpc-202307820")]
async fn recibir_prediccion(datos: web::Json<Quiniela>) -> impl Responder {
    // Por ahora, solo imprimimos el JSON en la consola de Rust
    println!("Predicción recibida de: {}", datos.username);
    println!("   Partido: {} {} - {} {}", datos.home_team, datos.home_goals, datos.away_goals, datos.away_team);
    
    // Respondemos a Locust que todo salió bien (HTTP 200)
    HttpResponse::Ok().json(serde_json::json!({
        "status": "success",
        "message": "Predicción procesada por la API en Rust"
    }))
}

// 3. Función principal que levanta el servidor
#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let puerto = 8080;
    println!("Servidor Rust levantado en http://0.0.0.0:{}", puerto);
    println!("Escuchando en la ruta /grpc-202307820"); // Cambia por tu carnet

    HttpServer::new(|| {
        App::new()
            .service(recibir_prediccion)
    })
    .bind(("0.0.0.0", puerto))?
    .run()
    .await
}