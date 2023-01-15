//go:build ignore

package main

//go:generate protoc --go_out=. --go-grpc_out=. --gocqrshttp_out=. contract.proto

//go:generate easyjson -all --lower_camel_case ./pkg/app/server/contracts/contract.pb.go

//go:generate wire ./pkg/app/rest
