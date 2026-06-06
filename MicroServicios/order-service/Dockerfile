# Dockerfile Multietapa para Microservicios Go
FROM golang:1.21-alpine AS builder

# Instalar git para dependencias
RUN apk add --no-cache git

WORKDIR /app

# Copiar go.mod y go.sum (si existe)
COPY go.mod ./
# Si el archivo go.sum no existe aún, el comando tidy lo creará
RUN go mod tidy

# Copiar el código fuente completo
COPY . .

# Compilar el binario (se pasa la ruta del main como argumento en el compose)
ARG SERVICE_PATH
RUN go build -o main ${SERVICE_PATH}

# Etapa final: Imagen ligera
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/main .

# Exponer puertos dinámicos si es necesario (Go Micro maneja esto)
CMD ["./main"]
