# E-commerce Microservices Platform

> Plataforma de comercio electrónico construida sobre una arquitectura de **microservicios**
> en Go, con comunicación **gRPC**, descubrimiento dinámico vía **etcd** y orquestación con
> **Docker Compose**.

---

## Objetivo

Proveer una plataforma de e-commerce **modular, escalable y desacoplada**, en la que cada
capacidad de negocio (usuarios, productos, pedidos) se implementa como un microservicio
independiente. Los servicios se descubren entre sí dinámicamente y se comunican por gRPC,
mientras que un **API Gateway** expone una única interfaz REST hacia el exterior.

---

## Descripción del proyecto

El sistema permite a los usuarios **registrarse e iniciar sesión**, **consultar un catálogo de
productos** y **gestionar pedidos** con validación de stock en tiempo real. Al crear un pedido,
el `order-service` verifica y descuenta el inventario consultando al `product-service` mediante
gRPC, garantizando consistencia entre la disponibilidad de productos y las órdenes generadas.

Toda la comunicación externa pasa por el **API Gateway** (REST/HTTP), que traduce las peticiones
a llamadas gRPC hacia los microservicios. Una **interfaz web** facilita la interacción del
usuario final.

---

## Stack tecnológico

| Categoría                | Tecnología                  |
|--------------------------|-----------------------------|
| Lenguaje                 | Go 1.25                     |
| Framework microservicios | Go Micro v4                 |
| Framework API (Gateway)  | Gin Gonic                   |
| Comunicación             | gRPC                        |
| Serialización            | Protocol Buffers (proto3)   |
| Service Discovery        | etcd v3.5                   |
| Orquestación             | Docker & Docker Compose     |
| Frontend                 | HTML5 / JavaScript (Nginx)  |

---

## Arquitectura

El sistema se compone de seis contenedores que cooperan en una red interna (`ecommerce-net`):

| Componente          | Rol                                                              | Puerto host |
|---------------------|-----------------------------------------------------------------|-------------|
| **etcd**            | Registro y descubrimiento de servicios (Service Discovery)      | 2379 / 2380 |
| **user-service**    | Registro, login y perfiles de usuario (gRPC)                    | interno     |
| **product-service** | Catálogo y control de stock: `CheckStock` / `UpdateStock` (gRPC)| interno     |
| **order-service**   | Gestión de pedidos; **cliente gRPC** de product-service         | interno     |
| **api-gateway**     | Fachada REST → gRPC (Gin)                                        | 8080        |
| **web-frontend**    | Interfaz web del cliente (Nginx)                                | 3000        |

```
   Navegador / Cliente
          │  HTTP / REST
          ▼
   ┌──────────────┐      gRPC      ┌───────────────┐
   │ API Gateway  │ ─────────────► │ user-service  │
   │  (Gin :8080) │ ─────────────► │ order-service │ ──gRPC──► product-service
   │              │ ─────────────► │ product-svc   │  (CheckStock / UpdateStock)
   └──────────────┘                └───────────────┘
          │ descubrimiento de servicios
          ▼
       ┌──────┐
       │ etcd │   (cada servicio se registra al iniciar)
       └──────┘
```

---

## Requisitos previos

Única herramienta necesaria para el despliegue:

- **Docker Desktop** (incluye Docker Compose)

Verificación:

```bash
docker --version
docker compose version
```

> Docker Compose construye y ejecuta todo automáticamente — no se requiere instalar Go ni otras
> dependencias en la máquina anfitriona.
>
> *(Opcional, solo para desarrollo local sin contenedores: Go 1.25 y `protoc`.)*

---

## Estructura del proyecto

```
GoMicroTeam/
├── api-gateway/                 # Fachada REST (Gin) → gRPC
│   ├── main.go
│   └── Dockerfile
├── MicroServicios/
│   ├── user-service/            # Servicio de usuarios (gRPC)
│   │   ├── main.go
│   │   ├── handler/
│   │   └── Dockerfile
│   ├── product-service/         # Servicio de productos y stock (gRPC)
│   │   ├── main.go
│   │   ├── handler/
│   │   └── Dockerfile
│   └── order-service/           # Servicio de pedidos (gRPC)
│       ├── main.go
│       ├── handler/
│       └── Dockerfile
├── proto/                       # Contratos Protocol Buffers + código generado
│   ├── user/      (user.proto + *.pb.go)
│   ├── product/   (product.proto + *.pb.go)
│   └── order/     (order.proto + *.pb.go)
├── deploy/
│   └── docker-compose.yml       # Orquestación de los 6 contenedores
├── index.html                   # Interfaz web (servida por Nginx)
├── go.mod / go.sum              # Dependencias Go
├── regenerate_proto.sh          # Regenera código desde los .proto
├── Manual_De_Usuario.md         # Manual de usuario (entrega formal)
├── Plan_de_validacion.md        # Plan de pruebas y validación
└── README.md
```

---

## Guía de despliegue

### 1. Situarse en la carpeta de despliegue

```bash
cd GoMicroTeam/deploy
```

### 2. Construir y levantar el sistema

```bash
docker compose up --build -d
```

El primer arranque compila los servicios y puede tardar varios minutos.

### 3. Verificar el estado de los contenedores

```bash
docker compose ps
```

Los seis contenedores deben aparecer en estado `Up` / `running` (etcd como `healthy`).

### 4. Acceder al sistema

| Acceso              | URL                          |
|---------------------|------------------------------|
| Interfaz Web        | http://localhost:3000        |
| API Gateway (REST)  | http://localhost:8080/api/v1 |

### 5. Ver logs (opcional)

```bash
docker compose logs -f                 # todos los servicios
docker compose logs -f order-service   # un servicio concreto
```

### 6. Detener el sistema

```bash
docker compose down       # detener (conserva imágenes)
docker compose down -v    # detener y limpiar volúmenes
```

---

## Contratos y API (Endpoints principales)

Base URL: `http://localhost:8080/api/v1`

### Autenticación

| Método | Ruta              | Auth | Descripción                     |
|--------|-------------------|------|---------------------------------|
| POST   | `/auth/register`  | No   | Registrar nuevo usuario         |
| POST   | `/auth/login`     | No   | Iniciar sesión (devuelve token) |

### Productos (acceso público)

| Método | Ruta             | Auth | Descripción           |
|--------|------------------|------|-----------------------|
| GET    | `/products`      | No   | Listar productos      |
| GET    | `/products/:id`  | No   | Ver un producto       |

### Pedidos (requiere header `Authorization: fake-jwt-token-for-<ID>`)

| Método | Ruta                   | Auth | Descripción                  |
|--------|------------------------|------|------------------------------|
| POST   | `/orders`              | Sí   | Crear pedido                 |
| GET    | `/orders`              | Sí   | Listar pedidos               |
| GET    | `/orders/:id`          | Sí   | Ver un pedido                |
| PATCH  | `/orders/:id/status`   | Sí   | Cambiar estado de un pedido  |

### Usuario

| Método | Ruta        | Auth | Descripción            |
|--------|-------------|------|------------------------|
| GET    | `/me/:id`   | Sí   | Ver perfil de usuario  |

### Contratos gRPC (Protocol Buffers)

Definidos en `proto/`. Servicios y operaciones clave:

- **UserService** (`proto/user/user.proto`): `Register`, `Login`, `GetUser`, `ListUsers`
- **ProductService** (`proto/product/product.proto`): `ListProducts`, `GetProduct`,
  **`CheckStock`**, **`UpdateStock`**
- **OrderService** (`proto/order/order.proto`): `CreateOrder`, `ListOrders`, `GetOrder`,
  `UpdateOrderStatus`

> Flujo interno destacado: `OrderService.CreateOrder` invoca por gRPC a
> `ProductService.CheckStock` y `ProductService.UpdateStock` para validar y descontar inventario.

### Ejemplo rápido

```bash
# 1. Login → obtener token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'

# 2. Crear pedido (usar el token devuelto)
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: fake-jwt-token-for-1" \
  -d '{"user_id":"1","items":[{"product_id":"P1","product_name":"Laptop Gamer","quantity":1,"unit_price":1200.0}]}'
```

> **Usuario de prueba precargado:** `admin@example.com` / `admin123`.

---

## Documentación adicional

- 📘 **[Manual de Usuario](./Manual_De_Usuario.md)** — guía completa de uso, funcionalidades y
  flujo de trabajo para el cliente final.
