# E-commerce Microservices Platform

##  Acceso al Sistema

Puedes acceder a las diferentes capas del sistema a través de los siguientes enlaces:

*   **Interfaz Web (Cliente):** http://localhost:3000
*   **API Gateway (Backend):** http://localhost:8080/api/v1
*   **Dashboard de Servicios:** http://localhost:8082

## 🛠 Tecnologías Utilizadas

*   **Lenguaje:** Go 1.25
*   **Framework:** Go Micro v4
*   **API Framework:** Gin Gonic (en el Gateway)
*   **Comunicación:** gRPC
*   **Service Discovery:** etcd
*   **Orquestación:** Docker & Docker Compose
*   **Frontend:** HTML5 / JavaScript (Servicio independiente via Nginx/Docker)
*   **Serialización:** Protocol Buffers (proto3)

## 🔗 Endpoints y Pruebas

El **API Gateway** expone la interfaz REST en `http://localhost:8080`.

### Productos (Acceso Público)
*   **Listar todos los productos:** `GET http://localhost:8080/api/v1/products`
*   **Ver producto específico:** `GET http://localhost:8080/api/v1/products/P1`

### Autenticación
*   **Registro:** `POST http://localhost:8080/api/v1/auth/register`
*   **Login:** `POST http://localhost:8080/api/v1/auth/login`

### Pedidos (Requiere Header `Authorization: fake-jwt-token-for-ID`)
*   **Listar pedidos:** `GET http://localhost:8080/api/v1/orders`
*   **Crear pedido:** `POST http://localhost:8080/api/v1/orders`
*   **Estado de pedido:** `PATCH http://localhost:8080/api/v1/orders/:id/status`