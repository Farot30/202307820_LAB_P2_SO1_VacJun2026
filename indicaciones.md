## Python
.\.venv\Scripts\activate
locust -f locustfile.py

## Rust
cargo run

## locust
http://localhost:8089

en host http://localhost:8080

numero de usuario: 2
spawn rate: 1

## api rust
docker build -f api_rust/Dockerfile -t mi-api-rust .

## go server
docker build -t mi-go-server .

## go consumer
docker build -t mi-go-consumer .

## inicializar instancia de VM

nohup ./zot serve config.json > zot.log 2>&1 &

## probar en la terminal lo de go consumer logs
kubectl logs go-consumer-xxxx-xxxx
o
kubectl logs -l app=go-consumer