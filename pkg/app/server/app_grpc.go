package server

import (
	"context"
	"fmt"
	"net"
	"prototodo/pkg/app/server/common"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (a *app) startGRPC(portStr string) {
	s := grpc.NewServer(

		// interceptor that
		// - handles panics
		// - extracts trace info from request
		// and replaces the grpc context with a golang context
		// that has traceinfo
		grpc.UnaryInterceptor(func(
			c context.Context,
			req interface{},
			info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler,
		) (resp interface{}, err error) {
			md, ok := metadata.FromIncomingContext(c)
			if !ok {
				return nil, fmt.Errorf("empty context")
			}

			temp := md["traceparent"]
			traceparent := ""
			if len(temp) > 0 {
				traceparent = temp[0]
			}

			ctx := a.ctxf.Create(traceparent)
			resp, err = handler(ctx, req)
			return
		}),
	)

	a.registerGRPCHandlers(s)
	a.registerCloser(s.GracefulStop)

	port, err := strconv.Atoi(portStr)
	if err != nil {
		a.lgr.Warn(
			"unable to parse provided port, setting port to default",
			zap.String("portConfig", portStr),
		)
		port = common.GRPCDefaultPort
	}
	if port < 0 {
		a.lgr.Warn(
			"port was configured was invalid, setting port to default",
		)
		port = common.GRPCDefaultPort
	}

	a.lgr.Info("grpc listening", zap.Int("port", port))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
