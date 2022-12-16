//go:build ignore

package prototodo

//go:generate protoc -I=. --go_out=./pkg/app contract.proto

//easyjson -all --lower_camel_case ./pkg/app/rest/dto/res/
//easyjson -all --lower_camel_case ./pkg/app/rest/dto/req/

//swag init -g cmd/server/main.go
//wire ./pkg/app/rest
