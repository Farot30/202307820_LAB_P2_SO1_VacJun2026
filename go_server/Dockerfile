FROM golang:1.21 as builder
WORKDIR /app

# Copiamos los módulos
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el código
COPY . .

# Compilamos el binario
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Imagen final super ligera (Alpine)
FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 50051
CMD ["./server"]