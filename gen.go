//go:build ignore

package prototodo

//go:generate protoc -I=pkg/infra/db --go_out=. pkg/infra/db/data.proto

//go:generate easyjson -all --lower_camel_case ./pkg/app/rest/dto/res/
//go:generate easyjson -all --lower_camel_case ./pkg/app/rest/dto/req/

//go:generate swag init -g cmd/server/main.go
//go:generate wire ./pkg/app/rest
