//go:build ignore

// Contains instructions for required generators, invoke by entering
// go generate gen.go

package main

//go:generate protoc --go_out=paths=source_relative:./pkg/domain/ proto/contracts/models.proto
//go:generate protoc --go-grpc_out=. --gocqrshttp_out=. proto/contracts/service.proto
//go:generate hndlrgen
//go:generate protoc -I=pkg/infra/impls/evcqrs/entities --go_out=. pkg/infra/impls/evcqrs/entities/data.proto
//go:generate cp pkg/app/server/contracts/service.http.json pkg/app/server/static/swagger/swagger.json

//go:generate wire ./pkg/app/server
