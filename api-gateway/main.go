// Package main - API Gateway
//
// Punto de entrada HTTP/REST que expone los microservicios internos al exterior.
// Utiliza Gin como router HTTP y Go Micro como cliente para comunicarse
// con los servicios internos (Order, Product, User) vía gRPC.
//
// Responsabilidades:
//   - Exponer endpoints REST para clientes externos (frontend, mobile)
//   - Enrutar peticiones al servicio interno correspondiente
//   - Middleware de autenticación (validar token del UserService)
//   - Agregación de respuestas cuando sea necesario
//
// Endpoints planificados:
//   POST   /api/v1/auth/register      → UserService.Register
//   POST   /api/v1/auth/login          → UserService.Login
//   GET    /api/v1/users/:id           → UserService.GetUser
//   GET    /api/v1/products            → ProductService.ListProducts
//   GET    /api/v1/products/:id        → ProductService.GetProduct
//   POST   /api/v1/orders              → OrderService.CreateOrder
//   GET    /api/v1/orders              → OrderService.ListOrders
//   GET    /api/v1/orders/:id          → OrderService.GetOrder
//   PATCH  /api/v1/orders/:id/status   → OrderService.UpdateOrderStatus
//
// Implementación a cargo de: Persona 5
package main

import (
	"fmt"
	"log"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-micro.dev/v4"
	
	pb_order "ecommerce-microservices/proto/order"
	pb_product "ecommerce-microservices/proto/product"
	pb_user "ecommerce-microservices/proto/user"
)

var (
	userClient    pb_user.UserService
	productClient pb_product.ProductService
	orderClient   pb_order.OrderService
)

func main() {
	fmt.Println("=== API Gateway ===")
	fmt.Println("Gateway HTTP/REST → microservicios gRPC")
	fmt.Println("Puerto HTTP: 8080")

	// TODO (Persona 5): Implementar el gateway
	//
	// Pasos esperados:
	// 1. Crear servicio Go Micro (web.NewService o micro.NewService)
	service := micro.NewService(
		micro.Name("api-gateway.client"),
	)
	service.Init()
	// 2. Crear clientes Go Micro para cada servicio:
	//    - orderClient  = pb_order.NewOrderService("order-service", service.Client())
	//    - productClient = pb_product.NewProductService("product-service", service.Client())
	//    - userClient   = pb_user.NewUserService("user-service", service.Client())
	userClient = pb_user.NewUserService("user-service", service.Client())
	productClient = pb_product.NewProductService("product-service", service.Client())
	orderClient = pb_order.NewOrderService("order-service", service.Client())
	// 3. Configurar router Gin con los endpoints listados arriba
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handleRegister)
			auth.POST("/login", handleLogin)
		}

		products := v1.Group("/products")
		{
			products.GET("", handleListProducts)
			products.GET("/:id", handleGetProduct)
		}

		// Rutas Protegidas (requieren Token)
		protected := v1.Group("")
		protected.Use(AuthMiddleware())
		{
			orders := protected.Group("/orders")
			{
				orders.POST("", handleCreateOrder)
				orders.GET("", handleListOrders)
				orders.GET("/:id", handleGetOrder)
				orders.PATCH("/:id/status", handleUpdateOrderStatus)
			}
			
			protected.GET("/me/:id", handleGetUser)
		}
	}
	// 4. Iniciar en puerto 8080
	r.Run(":8080")
	log.Println("API Gateway iniciado")
}

// --- Middlewares ---

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || !strings.HasPrefix(token, "fake-jwt-token-for-") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido o inválido"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// --- Handlers para User Service ---

func handleRegister(c *gin.Context) {
	var req pb_user.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := userClient.Register(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func handleLogin(c *gin.Context) {
	var req pb_user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := userClient.Login(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func handleGetUser(c *gin.Context) {
	id := c.Param("id")
	res, err := userClient.GetUser(context.Background(), &pb_user.GetUserRequest{UserId: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// --- Handlers para Product Service ---

func handleListProducts(c *gin.Context) {
	category := c.Query("category")
	res, err := productClient.ListProducts(context.Background(), &pb_product.ListProductsRequest{Category: category})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func handleGetProduct(c *gin.Context) {
	id := c.Param("id")
	res, err := productClient.GetProduct(context.Background(), &pb_product.GetProductRequest{ProductId: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// --- Handlers para Order Service ---

func handleCreateOrder(c *gin.Context) {
	var req pb_order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := orderClient.CreateOrder(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func handleListOrders(c *gin.Context) {
	userId := c.Query("user_id")
	res, err := orderClient.ListOrders(context.Background(), &pb_order.ListOrdersRequest{UserId: userId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func handleGetOrder(c *gin.Context) {
	id := c.Param("id")
	res, err := orderClient.GetOrder(context.Background(), &pb_order.GetOrderRequest{OrderId: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func handleUpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status    string `json:"status"`
		NewStatus string `json:"new_status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status := body.Status
	if status == "" {
		status = body.NewStatus
	}
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'status' o 'new_status' es requerido"})
		return
	}

	res, err := orderClient.UpdateOrderStatus(context.Background(), &pb_order.UpdateOrderStatusRequest{
		OrderId:   id,
		NewStatus: status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

