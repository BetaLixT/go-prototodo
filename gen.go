//go:build ignore

package main

//go:generate protoc --go_out=. --go-grpc_out=. --gen-gohttp=. contract.proto

//go:generate easyjson -all --lower_camel_case ./pkg/app/rest/dto/res/
//go:generate easyjson -all --lower_camel_case ./pkg/app/rest/dto/req/

//go:generate swag init -g cmd/server/main.go
//go:generate wire ./pkg/app/rest
