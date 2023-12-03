package gapi

import (
	"fmt"

	db "github.com/xiao-yangg/simplebank/db/sqlc"
	"github.com/xiao-yangg/simplebank/pb"
	"github.com/xiao-yangg/simplebank/token"
	"github.com/xiao-yangg/simplebank/util"
	"github.com/xiao-yangg/simplebank/worker"
)

// Server serves gRPC request
type Server struct {
	pb.UnimplementedSimpleBankServer
	config util.Config
	store db.Store
	tokenMaker token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a gRPC server
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey) // length 32
	if err != nil {
		return nil, fmt.Errorf("cannot create tokoen maker: %w", err)
	}

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}