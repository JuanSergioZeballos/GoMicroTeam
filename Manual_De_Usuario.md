# Manual de Usuario
## E-commerce Microservices Platform — GoMicroTeam

> **Documento de entrega formal**
> Plataforma de comercio electrónico construida sobre una arquitectura de microservicios.
> Este manual describe las tecnologías empleadas, la instalación y puesta en marcha, las
> funcionalidades disponibles y el flujo de trabajo recomendado para el uso del sistema.

---

## Índice

1. [Descripción general](#1-descripción-general)
2. [Tecnologías utilizadas](#2-tecnologías-utilizadas)
3. [Arquitectura del sistema](#3-arquitectura-del-sistema)
4. [Requisitos previos](#4-requisitos-previos)
5. [Instalación y puesta en marcha](#5-instalación-y-puesta-en-marcha)
6. [Accesos al sistema](#6-accesos-al-sistema)
7. [Funcionalidades](#7-funcionalidades)
8. [Flujo de trabajo recomendado](#8-flujo-de-trabajo-recomendado)
9. [Datos de prueba precargados](#9-datos-de-prueba-precargados)
10. [Resolución de problemas](#10-resolución-de-problemas)
11. [Detener el sistema](#11-detener-el-sistema)

---

## 1. Descripción general

**E-commerce Microservices Platform** es un sistema de comercio electrónico que permite a los
usuarios registrarse, autenticarse, consultar un catálogo de productos y gestionar pedidos
verificando la disponibilidad de stock en tiempo real.

El sistema está dividido en **microservicios independientes** que se comunican entre sí mediante
**gRPC** y se descubren dinámicamente a través de un **registro de servicios (etcd)**. Todo el
acceso externo se realiza a través de un **API Gateway** que expone una interfaz REST, y una
**interfaz web** facilita la interacción al usuario final.

---

## 2. Tecnologías utilizadas

| Categoría             | Tecnología                          | Uso en el sistema                                |
|-----------------------|-------------------------------------|--------------------------------------------------|
| Lenguaje              | Go 1.25                             | Implementación de todos los microservicios       |
| Framework microservicios | Go Micro v4                      | Estructura de servicios, cliente/servidor        |
| Framework API         | Gin Gonic                           | API Gateway (rutas REST, middleware)             |
| Comunicación          | gRPC                                | Comunicación entre microservicios                |
| Serialización         | Protocol Buffers (proto3)           | Contratos de mensajes entre servicios            |
| Service Discovery     | etcd v3.5                           | Registro y descubrimiento dinámico de servicios  |
| Orquestación          | Docker & Docker Compose             | Empaquetado y despliegue de contenedores         |
| Frontend              | HTML5 / JavaScript (Nginx)          | Interfaz web del cliente                         |
| Seguridad             | Middleware de autenticación (token) | Protección de rutas de pedidos                   |

---

## 3. Arquitectura del sistema

El sistema se compone de seis contenedores que cooperan en una red interna:

| Componente        | Rol                                                                 | Puerto |
|-------------------|---------------------------------------------------------------------|--------|
| **etcd**          | Registro de servicios (Service Discovery)                           | 2379   |
| **user-service**  | Registro, login y gestión de perfiles de usuario (gRPC)             | interno|
| **product-service** | Catálogo de productos y control de stock (`CheckStock`/`UpdateStock`) | interno|
| **order-service** | Gestión de pedidos; consume product-service vía gRPC                 | interno|
| **api-gateway**   | Fachada REST que traduce HTTP → gRPC                                 | 8080   |
| **web-frontend**  | Interfaz web del cliente (Nginx)                                     | 3000   |

**Diagrama de flujo (resumen):**

```
   Navegador / Cliente
          │  HTTP/REST
          ▼
   ┌──────────────┐      gRPC      ┌──────────────┐
   │ API Gateway  │ ─────────────► │ user-service │
   │  (Gin :8080) │ ─────────────► │ order-service│ ──gRPC──► product-service
   └──────────────┘ ─────────────► │ product-svc  │   (CheckStock / UpdateStock)
          │                        └──────────────┘
          │ descubrimiento de servicios
          ▼
       ┌──────┐
       │ etcd │   (todos los servicios se registran aquí al iniciar)
       └──────┘
```

---

## 4. Requisitos previos

Antes de instalar, asegúrese de tener instalado:

- **Docker Desktop** (incluye Docker Compose).

Verifique la instalación:

```bash
docker --version
docker compose version
```

> No es necesario instalar Go ni ninguna otra dependencia: Docker Compose construye y ejecuta
> todo automáticamente.

---

## 5. Instalación y puesta en marcha

### Paso 1 — Obtener el proyecto

Descomprima/clone el proyecto y abra una terminal en la carpeta raíz `GoMicroTeam`.

### Paso 2 — Situarse en la carpeta de despliegue

```bash
cd GoMicroTeam/deploy
```

### Paso 3 — Construir y levantar el sistema

```bash
docker compose up --build -d
```

Este comando construye las imágenes y arranca los seis contenedores en segundo plano. El primer
arranque puede tardar varios minutos (compilación de los servicios).

### Paso 4 — Verificar que todo está activo

```bash
docker compose ps
```

Debe mostrar los seis contenedores en estado `Up` / `running` (etcd como `healthy`).

### Paso 5 — Acceder al sistema

Abra el navegador en **http://localhost:3000** (interfaz web).

---

## 6. Accesos al sistema

| Acceso                   | URL                                  | Descripción                          |
|--------------------------|--------------------------------------|--------------------------------------|
| Interfaz Web (Cliente)   | http://localhost:3000                | Punto de entrada para el usuario     |
| API Gateway (REST)       | http://localhost:8080/api/v1         | API para integraciones / pruebas     |

### Endpoints principales de la API

| Método | Ruta                              | Autenticación | Descripción                       |
|--------|-----------------------------------|---------------|-----------------------------------|
| POST   | `/api/v1/auth/register`           | No            | Registrar nuevo usuario           |
| POST   | `/api/v1/auth/login`              | No            | Iniciar sesión (devuelve token)   |
| GET    | `/api/v1/products`                | No            | Listar productos                  |
| GET    | `/api/v1/products/:id`            | No            | Ver un producto                   |
| POST   | `/api/v1/orders`                  | Sí (token)    | Crear pedido                      |
| GET    | `/api/v1/orders`                  | Sí (token)    | Listar pedidos                    |
| GET    | `/api/v1/orders/:id`              | Sí (token)    | Ver un pedido                     |
| PATCH  | `/api/v1/orders/:id/status`       | Sí (token)    | Cambiar estado de un pedido       |
| GET    | `/api/v1/me/:id`                  | Sí (token)    | Ver perfil de usuario             |

> **Autenticación:** las rutas protegidas requieren el header
> `Authorization: fake-jwt-token-for-<ID>`, obtenido al iniciar sesión.

---

## 7. Funcionalidades

La interfaz web (http://localhost:3000) organiza las funciones en las siguientes secciones:

### 7.1 🔐 Autenticación
- **Registro** de nuevos usuarios (nombre, email, contraseña). Los usuarios nuevos reciben el
  rol `CUSTOMER`.
- **Login** con email y contraseña. Devuelve un token de sesión necesario para operar con pedidos.

### 7.2 📦 Productos
- **Listado del catálogo** con nombre, descripción, precio, stock y categoría.
- **Consulta de un producto** específico por su identificador.

### 7.3 ➕ Crear Pedido
- Selección de productos y cantidades.
- Al crear el pedido, el sistema **verifica automáticamente el stock** (`CheckStock`) y, si hay
  disponibilidad, **descuenta las unidades** (`UpdateStock`) mediante comunicación gRPC interna.
- El pedido se crea en estado `PENDING`.

### 7.4 📋 Pedidos
- **Listado de pedidos** (todos o filtrados por usuario).
- **Consulta de un pedido** por su ID.
- **Cambio de estado** del pedido a lo largo de su ciclo de vida:
  `PENDING → CONFIRMED → SHIPPED → DELIVERED` (o `CANCELLED`).

### 7.5 👤 Mi Información
- Consulta del **perfil del usuario** autenticado (ID, nombre, email, rol).

---

## 8. Flujo de trabajo recomendado

Secuencia típica de uso de un cliente final:

1. **Acceder** a la interfaz web: http://localhost:3000
2. **Registrarse** (o iniciar sesión con un usuario existente) en la sección *Autenticación*.
   Al iniciar sesión se obtiene el token de sesión.
3. **Explorar el catálogo** en la sección *Productos* y anotar los identificadores y stock
   disponible de los artículos deseados.
4. **Crear un pedido** en la sección *Crear Pedido*, indicando productos y cantidades.
   El sistema valida el stock y lo descuenta automáticamente si la operación es válida.
5. **Revisar los pedidos** en la sección *Pedidos*: consultar el detalle y verificar el estado.
6. **Actualizar el estado** del pedido conforme avanza la gestión
   (`CONFIRMED`, `SHIPPED`, `DELIVERED`).
7. **Consultar el perfil** propio en *Mi Información* cuando se requiera.

> **Regla de negocio:** un pedido solo se crea si **todos** sus productos tienen stock
> suficiente. Si algún producto no tiene disponibilidad, el pedido se rechaza y no se descuenta
> stock de ningún artículo.

---

## 9. Datos de prueba precargados

El sistema incluye datos de ejemplo para realizar pruebas inmediatas.

**Usuario administrador:**

| Campo    | Valor               |
|----------|---------------------|
| Email    | `admin@example.com` |
| Password | `admin123`          |
| Rol      | `ADMIN`             |

**Productos:**

| ID | Nombre           | Precio  | Stock | Categoría    |
|----|------------------|---------|-------|--------------|
| P1 | Laptop Gamer     | 1200.00 | 10    | Electronics  |
| P2 | Mouse Wireless   | 40.00   | 50    | Peripherals  |
| P3 | Teclado Mecánico | 90.00   | 20    | Peripherals  |

**Pedidos de ejemplo:** `ORD-001`, `ORD-002` (usuario 1) y `ORD-003` (usuario 2).

> **Nota:** los datos se almacenan en memoria. Al reiniciar los contenedores, vuelven a su
> estado inicial.

---

## 10. Resolución de problemas

| Síntoma                          | Causa probable / Solución                                               |
|----------------------------------|-------------------------------------------------------------------------|
| La web no carga en :3000         | Verifique `docker compose ps`; reinicie con `docker compose up -d`.      |
| Error 401 al crear pedido        | Falta iniciar sesión; use el token devuelto por el login.               |
| Pedido rechazado por stock       | El producto no tiene unidades suficientes; reduzca la cantidad.         |
| Un contenedor no inicia          | Revise logs: `docker compose logs <servicio>`.                          |
| Puerto 8080 / 3000 ocupado       | Cierre el proceso que usa el puerto o cámbielo en `docker-compose.yml`. |

Ver logs en tiempo real de todos los servicios:

```bash
docker compose logs -f
```

---

## 11. Detener el sistema

```bash
# Detener los contenedores (conserva las imágenes)
docker compose down

# Detener y eliminar también los volúmenes (reinicio limpio)
docker compose down -v
```

---

*Fin del Manual de Usuario — GoMicroTeam E-commerce Microservices Platform.*
