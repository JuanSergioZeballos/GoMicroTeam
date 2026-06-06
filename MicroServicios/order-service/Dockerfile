# ============================================================
# Dockerfile genérico para microservicios Go
# ============================================================
# Multi-stage build:
#   1. Etapa de compilación: usa golang:1.21-alpine
#   2. Etapa de ejecución: usa alpine:3.18 (imagen mínima)
# ============================================================

# --- Etapa 1: Compilación ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias primero (cache de Docker)
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true

# Copiar todo el código fuente
COPY . .

# Compilar el servicio
RUN CGO_ENABLED=0 GOOS=linux go build -o /service ./order-service/

# --- Etapa 2: Ejecución ---
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /service .

EXPOSE 8080

CMD ["./service"]
