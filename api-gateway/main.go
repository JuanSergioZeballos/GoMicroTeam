// Package main - API Gateway
package main

import (
	"fmt"
	"log"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"go-micro.dev/v4"
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	_ "github.com/go-micro/plugins/v4/client/grpc"
	_ "github.com/go-micro/plugins/v4/server/grpc"
	_ "github.com/go-micro/plugins/v4/transport/grpc"
	
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

	service := micro.NewService()
	service.Init(
		micro.Name("api-gateway.client"),
	)
	userClient = pb_user.NewUserService("user-service", service.Client())
	productClient = pb_product.NewProductService("product-service", service.Client())
	orderClient = pb_order.NewOrderService("order-service", service.Client())
	
	r := gin.Default()
	
	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	
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
	r.Run(":8080")
	log.Println("API Gateway iniciado")
}

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
