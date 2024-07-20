package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	token "github.com/thetherington/jobber-common/gateway-token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type key int

const (
	CtxUsernameKey key = iota
	CtxEmailKey
)

var (
	customFunc recovery.RecoveryHandlerFunc

	TokenIds    = []string{"auth", "users", "gig", "search", "buyer", "message", "order", "review"}
	MetaKeys    = []string{"username", "email"}
	ContextKeys = []key{CtxUsernameKey, CtxEmailKey}
)

func CreateTokenReceiverInterceptor(secret, serviceName string) grpc.UnaryServerInterceptor {
	tokenMaker, err := token.NewGatewayJWTMaker(secret)
	if err != nil {
		slog.With("error", err).Error(fmt.Sprintf("failed to create token maker for %s", serviceName))
		panic(err)
	}

	return auth.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) {
		// get the token from the header "authorization" and "bearer" label
		token, err := auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		// validate the token is signed by api gateway secret
		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid gateway token, request not coming from gateway: %v", err)
		}

		// check the id in the claims is one of the slice values that is hard coded.
		if !slices.Contains(TokenIds, payload.ID) {
			return nil, status.Error(codes.Unauthenticated, "invalid request, request payload invalid")
		}

		return ctx, nil
	})
}

func CreateMetadataReceiverInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			for i, k := range MetaKeys {
				if v := md.Get(k); len(v) > 0 {
					ctx = context.WithValue(ctx, ContextKeys[i], v[0])
				}
			}
		}

		return handler(ctx, req)
	}
}

func CreatePanicRecoveryInterceptor() grpc.UnaryServerInterceptor {
	customFunc = func(p any) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}

	opts := []recovery.Option{
		recovery.WithRecoveryHandler(customFunc),
	}

	return recovery.UnaryServerInterceptor(opts...)
}

func CreateStreamPanicRecoveryInterceptor() grpc.StreamServerInterceptor {
	customFunc = func(p any) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}

	opts := []recovery.Option{
		recovery.WithRecoveryHandler(customFunc),
	}

	return recovery.StreamServerInterceptor(opts...)
}

func CreateStreamTokenReceiverInterceptor(secret, serviceName string) grpc.StreamServerInterceptor {
	tokenMaker, err := token.NewGatewayJWTMaker(secret)
	if err != nil {
		slog.With("error", err).Error(fmt.Sprintf("failed to create token maker for %s", serviceName))
		panic(err)
	}

	return auth.StreamServerInterceptor(func(ctx context.Context) (context.Context, error) {
		// get the token from the header "authorization" and "bearer" label
		token, err := auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		// validate the token is signed by api gateway secret
		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid gateway token, request not coming from gateway: %v", err)
		}

		// check the id in the claims is one of the slice values that is hard coded.
		if !slices.Contains(TokenIds, payload.ID) {
			return nil, status.Error(codes.Unauthenticated, "invalid request, request payload invalid")
		}

		return ctx, nil
	})
}
