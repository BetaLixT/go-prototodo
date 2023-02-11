// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: contracts/service.proto

package contracts

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	contracts "prototodo/pkg/domain/contracts"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TasksClient is the client API for Tasks service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TasksClient interface {
	// - Commands
	Create(ctx context.Context, in *contracts.CreateTaskCommand, opts ...grpc.CallOption) (*contracts.TaskEvent, error)
	Delete(ctx context.Context, in *contracts.DeleteTaskCommand, opts ...grpc.CallOption) (*contracts.TaskEvent, error)
	Update(ctx context.Context, in *contracts.UpdateTaskCommand, opts ...grpc.CallOption) (*contracts.TaskEvent, error)
	// Update existing task state to progress
	Progress(ctx context.Context, in *contracts.ProgressTaskCommand, opts ...grpc.CallOption) (*contracts.TaskEvent, error)
	// Update existing task to complete
	Complete(ctx context.Context, in *contracts.CompleteTaskCommand, opts ...grpc.CallOption) (*contracts.TaskEvent, error)
	// Query for existing tasks
	ListQuery(ctx context.Context, in *contracts.ListTasksQuery, opts ...grpc.CallOption) (*contracts.TaskEntityList, error)
}

type tasksClient struct {
	cc grpc.ClientConnInterface
}

func NewTasksClient(cc grpc.ClientConnInterface) TasksClient {
	return &tasksClient{cc}
}

func (c *tasksClient) Create(ctx context.Context, in *contracts.CreateTaskCommand, opts ...grpc.CallOption) (*contracts.TaskEvent, error) {
	out := new(contracts.TaskEvent)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksClient) Delete(ctx context.Context, in *contracts.DeleteTaskCommand, opts ...grpc.CallOption) (*contracts.TaskEvent, error) {
	out := new(contracts.TaskEvent)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksClient) Update(ctx context.Context, in *contracts.UpdateTaskCommand, opts ...grpc.CallOption) (*contracts.TaskEvent, error) {
	out := new(contracts.TaskEvent)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksClient) Progress(ctx context.Context, in *contracts.ProgressTaskCommand, opts ...grpc.CallOption) (*contracts.TaskEvent, error) {
	out := new(contracts.TaskEvent)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/Progress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksClient) Complete(ctx context.Context, in *contracts.CompleteTaskCommand, opts ...grpc.CallOption) (*contracts.TaskEvent, error) {
	out := new(contracts.TaskEvent)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/Complete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksClient) ListQuery(ctx context.Context, in *contracts.ListTasksQuery, opts ...grpc.CallOption) (*contracts.TaskEntityList, error) {
	out := new(contracts.TaskEntityList)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/ListQuery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TasksServer is the server API for Tasks service.
// All implementations must embed UnimplementedTasksServer
// for forward compatibility
type TasksServer interface {
	// - Commands
	Create(context.Context, *contracts.CreateTaskCommand) (*contracts.TaskEvent, error)
	Delete(context.Context, *contracts.DeleteTaskCommand) (*contracts.TaskEvent, error)
	Update(context.Context, *contracts.UpdateTaskCommand) (*contracts.TaskEvent, error)
	// Update existing task state to progress
	Progress(context.Context, *contracts.ProgressTaskCommand) (*contracts.TaskEvent, error)
	// Update existing task to complete
	Complete(context.Context, *contracts.CompleteTaskCommand) (*contracts.TaskEvent, error)
	// Query for existing tasks
	ListQuery(context.Context, *contracts.ListTasksQuery) (*contracts.TaskEntityList, error)
	mustEmbedUnimplementedTasksServer()
}

// UnimplementedTasksServer must be embedded to have forward compatible implementations.
type UnimplementedTasksServer struct {
}

func (UnimplementedTasksServer) Create(context.Context, *contracts.CreateTaskCommand) (*contracts.TaskEvent, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedTasksServer) Delete(context.Context, *contracts.DeleteTaskCommand) (*contracts.TaskEvent, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedTasksServer) Update(context.Context, *contracts.UpdateTaskCommand) (*contracts.TaskEvent, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedTasksServer) Progress(context.Context, *contracts.ProgressTaskCommand) (*contracts.TaskEvent, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Progress not implemented")
}
func (UnimplementedTasksServer) Complete(context.Context, *contracts.CompleteTaskCommand) (*contracts.TaskEvent, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Complete not implemented")
}
func (UnimplementedTasksServer) ListQuery(context.Context, *contracts.ListTasksQuery) (*contracts.TaskEntityList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListQuery not implemented")
}
func (UnimplementedTasksServer) mustEmbedUnimplementedTasksServer() {}

// UnsafeTasksServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TasksServer will
// result in compilation errors.
type UnsafeTasksServer interface {
	mustEmbedUnimplementedTasksServer()
}

func RegisterTasksServer(s grpc.ServiceRegistrar, srv TasksServer) {
	s.RegisterService(&Tasks_ServiceDesc, srv)
}

func _Tasks_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(contracts.CreateTaskCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).Create(ctx, req.(*contracts.CreateTaskCommand))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tasks_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(contracts.DeleteTaskCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).Delete(ctx, req.(*contracts.DeleteTaskCommand))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tasks_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(contracts.UpdateTaskCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).Update(ctx, req.(*contracts.UpdateTaskCommand))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tasks_Progress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(contracts.ProgressTaskCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).Progress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/Progress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).Progress(ctx, req.(*contracts.ProgressTaskCommand))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tasks_Complete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(contracts.CompleteTaskCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).Complete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/Complete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).Complete(ctx, req.(*contracts.CompleteTaskCommand))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tasks_ListQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(contracts.ListTasksQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).ListQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/ListQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).ListQuery(ctx, req.(*contracts.ListTasksQuery))
	}
	return interceptor(ctx, in, info, handler)
}

// Tasks_ServiceDesc is the grpc.ServiceDesc for Tasks service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Tasks_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tasks.Tasks",
	HandlerType: (*TasksServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Tasks_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Tasks_Delete_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Tasks_Update_Handler,
		},
		{
			MethodName: "Progress",
			Handler:    _Tasks_Progress_Handler,
		},
		{
			MethodName: "Complete",
			Handler:    _Tasks_Complete_Handler,
		},
		{
			MethodName: "ListQuery",
			Handler:    _Tasks_ListQuery_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "contracts/service.proto",
}

// QuotesClient is the client API for Quotes service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QuotesClient interface {
	// Get a quote
	Get(ctx context.Context, in *contracts.GetQuoteQuery, opts ...grpc.CallOption) (*contracts.QuoteData, error)
	Create(ctx context.Context, in *contracts.CreateQuoteCommand, opts ...grpc.CallOption) (*contracts.QuoteData, error)
}

type quotesClient struct {
	cc grpc.ClientConnInterface
}

func NewQuotesClient(cc grpc.ClientConnInterface) QuotesClient {
	return &quotesClient{cc}
}

func (c *quotesClient) Get(ctx context.Context, in *contracts.GetQuoteQuery, opts ...grpc.CallOption) (*contracts.QuoteData, error) {
	out := new(contracts.QuoteData)
	err := c.cc.Invoke(ctx, "/tasks.Quotes/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *quotesClient) Create(ctx context.Context, in *contracts.CreateQuoteCommand, opts ...grpc.CallOption) (*contracts.QuoteData, error) {
	out := new(contracts.QuoteData)
	err := c.cc.Invoke(ctx, "/tasks.Quotes/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QuotesServer is the server API for Quotes service.
// All implementations must embed UnimplementedQuotesServer
// for forward compatibility
type QuotesServer interface {
	// Get a quote
	Get(context.Context, *contracts.GetQuoteQuery) (*contracts.QuoteData, error)
	Create(context.Context, *contracts.CreateQuoteCommand) (*contracts.QuoteData, error)
	mustEmbedUnimplementedQuotesServer()
}

// UnimplementedQuotesServer must be embedded to have forward compatible implementations.
type UnimplementedQuotesServer struct {
}

func (UnimplementedQuotesServer) Get(context.Context, *contracts.GetQuoteQuery) (*contracts.QuoteData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedQuotesServer) Create(context.Context, *contracts.CreateQuoteCommand) (*contracts.QuoteData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedQuotesServer) mustEmbedUnimplementedQuotesServer() {}

// UnsafeQuotesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QuotesServer will
// result in compilation errors.
type UnsafeQuotesServer interface {
	mustEmbedUnimplementedQuotesServer()
}

func RegisterQuotesServer(s grpc.ServiceRegistrar, srv QuotesServer) {
	s.RegisterService(&Quotes_ServiceDesc, srv)
}

func _Quotes_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(contracts.GetQuoteQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuotesServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Quotes/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuotesServer).Get(ctx, req.(*contracts.GetQuoteQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _Quotes_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(contracts.CreateQuoteCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuotesServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Quotes/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuotesServer).Create(ctx, req.(*contracts.CreateQuoteCommand))
	}
	return interceptor(ctx, in, info, handler)
}

// Quotes_ServiceDesc is the grpc.ServiceDesc for Quotes service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Quotes_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tasks.Quotes",
	HandlerType: (*QuotesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Quotes_Get_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _Quotes_Create_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "contracts/service.proto",
}
