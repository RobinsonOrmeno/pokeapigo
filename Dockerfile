# Etapa 1: Build
FROM golang:1.22 AS builder
WORKDIR /app

# Copiar los archivos de tu aplicación
COPY go.mod ./
RUN go mod download

COPY . .

# Compilar la aplicación
RUN go build -o main .

# Etapa 2: Run
FROM ubuntu:latest
WORKDIR /app

# Instalar certificados raíz para TLS
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copiar binario desde la etapa de construcción
COPY --from=builder /app/main .

# Exponer el puerto en el que escucha tu API
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
