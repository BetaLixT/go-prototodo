//go:build ignore

package main

//go:generate protoc --go-grpc_out=. --gocqrshttp_out=. contracts/service.proto
//go:generate protoc --go_out=paths=source_relative:./pkg/domain/ contracts/models.proto
//go:generate hndlrgen

//go:generate easyjson -all --lower_camel_case ./pkg/domain/contracts/models.pb.go

//go:generate wire ./pkg/app/rest
