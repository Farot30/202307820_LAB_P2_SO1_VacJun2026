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

kubectl logs -l app=go-consumer

## ver servicios en cloud shell y otras cosas

kubectl get services

kubectl get vmi

kubectl get pods

kubectl get svc

kubectl get hpa

kubectl apply -f apps.yaml

## entrar a la maquina de grafana

./virtctl console grafana-vm

tail -f /var/log/cloud-init-output.log

## para ver valkey y sus llaves
kubectl run -i --tty --rm valkey-client --image=valkey/valkey:latest --restart=Never -- valkey-cli -h valkey-service.default -p 6379

## trabajar con Ngrok

curl -O https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-v3-stable-linux-amd64.tgz
tar -xvzf ngrok-v3-stable-linux-amd64.tgz

dayroom-outlying-magnifier.ngrok-free.dev

./ngrok config add-authtoken

./ngrok http 5000

## hacer push 

docker push dayroom-outlying-magnifier.ngrok-free.dev/mi-api-rust:latest

docker push dayroom-outlying-magnifier.ngrok-free.dev/mi-go-consumer:v3

docker push dayroom-outlying-magnifier.ngrok-free.dev/mi-go-server:latest

kubectl get gateway

## limpiar valkey
kubectl run limpiador-db --rm -it --image=redis:alpine -- redis-cli -h valkey-service FLUSHALL

kubectl delete pod limpiador-db --force --grace-period=0

kubectl run explorador-db --rm -it --image=redis:alpine -- redis-cli -h valkey-service

KEYS *

XREVRANGE predicciones_stream + - COUNT 3

ZREVRANGE top_usuarios 0 5 WITHSCORES


kubectl describe pod go-server-7dcd64dcbb-7rlcq   | grep -i "image"


curl -i -X POST http://34.130.237.145/grpc-202307820 -H "Content-Type: application/json" -d '{
  "home_team": "GTM",
  "away_team": "MEX",
  "home_goals": 3,
  "away_goals": 5,
  "username": "tester_dapr",
  "timestamp": "2023-10-25T12:00:00Z"
}'


