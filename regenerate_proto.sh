#!/bin/bash
set -e

# Install protoc and protoc plugins compatible with grpc 1.53
apk add --no-cache protobuf-dev
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

# Add to PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Regenerate proto files
protoc --go_out=. --go-grpc_out=. proto/order/order.proto
protoc --go_out=. --go-grpc_out=. proto/product/product.proto
protoc --go_out=. --go-grpc_out=. proto/user/user.proto

echo "Proto files regenerated successfully!"
