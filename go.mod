module ecommerce-microservices

go 1.21

// Las dependencias se agregarán al ejecutar `go mod tidy` tras implementar los servicios.
// Dependencias principales esperadas:
//   go-micro.dev/v4          — framework de microservicios
//   google.golang.org/protobuf — soporte protobuf
//   github.com/gin-gonic/gin  — HTTP router para el API Gateway
