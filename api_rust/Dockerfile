# Usamos una imagen oficial de Rust para compilar
FROM rust:latest AS builder
WORKDIR /usr/src/api_rust

RUN apt-get update && apt-get install -y protobuf-compiler

# 1. Copiamos la carpeta proto desde la raíz del proyecto
COPY proto /usr/src/proto

# 2. Copiamos los archivos de configuración de la API especificando la carpeta
COPY api_rust/Cargo.toml api_rust/Cargo.lock ./
COPY api_rust/build.rs ./
COPY api_rust/src ./src

# Compilamos la aplicación para producción
RUN cargo build --release

# Usamos una imagen más ligera de Debian para correr la app final
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y libssl3 && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/src/api_rust/target/release/api_rust /usr/local/bin/api_rust

# Exponemos el puerto
EXPOSE 8080
CMD ["api_rust"]