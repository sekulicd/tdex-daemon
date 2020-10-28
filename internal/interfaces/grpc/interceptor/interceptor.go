package interceptor

import (
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/tdex-network/tdex-daemon/config"
	dbbadger "github.com/tdex-network/tdex-daemon/internal/infrastructure/storage/db/badger"
	"github.com/tdex-network/tdex-daemon/pkg/macaroons"
	"google.golang.org/grpc"
)

// UnaryInterceptor returns the unary interceptor
func UnaryInterceptor(
	dbManager *dbbadger.DbManager,
	macaroonService macaroons.Service,
) grpc.ServerOption {
	return grpc.UnaryInterceptor(
		middleware.ChainUnaryServer(
			unaryAuthHandler(config.RPCServerPermissions(), macaroonService),
			unaryLogger,
			unaryTransactionHandler(dbManager),
		),
	)
}

// StreamInterceptor returns the stream interceptor with a logrus log
func StreamInterceptor(
	dbManager *dbbadger.DbManager,
	macaroonService macaroons.Service,
) grpc.ServerOption {
	return grpc.StreamInterceptor(
		middleware.ChainStreamServer(
			streamAuthHandler(config.RPCServerPermissions(), macaroonService),
			streamLogger,
			streamTransactionHandler(dbManager),
		),
	)
}
