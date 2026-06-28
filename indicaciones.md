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

## activar la maquina de grafana
kubectl apply -f grafana.yaml

./virtctl console grafana-vm

## para ver valkey y sus llaves
kubectl run -i --tty --rm valkey-client --image=valkey/valkey:latest --restart=Never -- valkey-cli -h valkey-service.default -p 6379

## trabajar con Ngrok

curl -O https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-v3-stable-linux-amd64.tgz
tar -xvzf ngrok-v3-stable-linux-amd64.tgz

./ngrok config add-authtoken

./ngrok http 5000