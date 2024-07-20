package grpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	clientToken "github.com/thetherington/jobber-common/client-token"
	gatewayToken "github.com/thetherington/jobber-common/gateway-token"
	"github.com/thetherington/jobber-gateway/internal/adapter/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// create a interceptor that extracts the username and email from the session store http request cookie
// create the gRPC headers username and email of the logged in user
func CreateContextForwarderInterceptor(session *scs.SessionManager) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		payload, ok := session.Get(ctx, "currentUser").(*clientToken.Payload)
		if ok {
			ctx = metadata.AppendToOutgoingContext(ctx,
				"username", payload.Username,
			)

			ctx = metadata.AppendToOutgoingContext(ctx,
				"email", payload.Email,
			)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// create a interceptor that creates and inserts a token into the gRPC authrization header. none token expiry
// that validates the request came from the gateway to the microservice.
func CreateTokenForwarderInterceptor(serviceName string) grpc.UnaryClientInterceptor {
	tokenMaker, err := gatewayToken.NewGatewayJWTMaker(config.Config.Tokens.GATEWAY)
	if err != nil {
		slog.With("error", err).Error(fmt.Sprintf("failed to create token maker for %s", serviceName))
		panic(err)
	}

	token, _, err := tokenMaker.CreateToken(serviceName)
	if err != nil {
		slog.With("error", err).Error(fmt.Sprintf("failed to generate gateway token for %s", serviceName))
		panic(err)
	}

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx,
			"authorization", fmt.Sprintf("Bearer %s", token),
		)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// create a interceptor that creates and inserts a token into the gRPC authrization header. none token expiry
// that validates the request came from the gateway to the microservice.
func CreateStreamTokenForwarderInterceptor(serviceName string) grpc.StreamClientInterceptor {
	tokenMaker, err := gatewayToken.NewGatewayJWTMaker(config.Config.Tokens.GATEWAY)
	if err != nil {
		slog.With("error", err).Error(fmt.Sprintf("failed to create token maker for %s", serviceName))
		panic(err)
	}

	token, _, err := tokenMaker.CreateToken(serviceName)
	if err != nil {
		slog.With("error", err).Error(fmt.Sprintf("failed to generate gateway token for %s", serviceName))
		panic(err)
	}

	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx,
			"authorization", fmt.Sprintf("Bearer %s", token),
		)

		return streamer(ctx, desc, cc, method, opts...)
	}
}
