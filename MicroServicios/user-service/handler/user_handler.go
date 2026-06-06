// Package handler - User Service Handler
//
// Implementa la lógica de negocio del UserService definida en user.proto.
// Mantiene un almacén de usuarios en memoria.
//
// Implementación a cargo de: Persona 4
package handler

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "ecommerce-microservices/proto/user"
)

// UserHandler implementa la interfaz UserServiceHandler generada por protobuf.
type UserHandler struct {
	// TODO (Persona 4): Agregar campos necesarios
	// - users map[string]*pb.User  (almacén en memoria)
	mu        sync.RWMutex
	users     map[string]*pb.User
	passwords map[string]string
}

func NewUserHandler() *UserHandler {
	h := &UserHandler{
		users:     make(map[string]*pb.User),
		passwords: make(map[string]string),
	}
	
	h.users["1"] = &pb.User{
		Id:        "1",
		Name:      "Admin User",
		Email:     "admin@example.com",
		Role:      "ADMIN",
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	h.passwords["admin@example.com"] = "admin123"

	return h
}
// Register registra un nuevo usuario en el sistema.
// func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest, rsp *pb.RegisterResponse) error {
//     return nil
// }

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest, rsp *pb.RegisterResponse) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, u := range h.users {
		if u.Email == req.Email {
			rsp.Success = false
			rsp.Message = "El email ya está registrado"
			return nil
		}
	}

	id := fmt.Sprintf("%d", len(h.users)+1)
	newUser := &pb.User{
		Id:        id,
		Name:      req.Name,
		Email:     req.Email,
		Role:      "CUSTOMER",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	h.users[id] = newUser
	h.passwords[req.Email] = req.Password

	rsp.User = newUser
	rsp.Success = true
	rsp.Message = "Usuario registrado exitosamente"
	return nil
}


// Login autentica un usuario y devuelve un token.
// func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
//     return nil
// }

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, u := range h.users {
		if u.Email == req.Email {
			if h.passwords[req.Email] == req.Password {
				rsp.User = u
				rsp.Token = fmt.Sprintf("fake-jwt-token-for-%s", u.Id)
				rsp.Success = true
				rsp.Message = "Login exitoso"
				return nil
			}
			break
		}
	}

	rsp.Success = false
	rsp.Message = "Credenciales inválidas"
	return nil
}

// GetUser obtiene el perfil de un usuario por su ID.
// func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest, rsp *pb.GetUserResponse) error {
//     return nil
// }

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest, rsp *pb.GetUserResponse) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if user, ok := h.users[req.UserId]; ok {
		rsp.User = user
		rsp.Found = true
	} else {
		rsp.Found = false
	}
	return nil
}

// ListUsers lista todos los usuarios registrados.
// func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest, rsp *pb.ListUsersResponse) error {
//     return nil
// }

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest, rsp *pb.ListUsersResponse) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, user := range h.users {
		rsp.Users = append(rsp.Users, user)
	}
	return nil
}