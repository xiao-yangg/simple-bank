package gapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	db "github.com/xiao-yangg/simplebank/db/sqlc"
	"github.com/xiao-yangg/simplebank/pb"
	"github.com/xiao-yangg/simplebank/util"
	"github.com/xiao-yangg/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if violations := validateLoginUserRequest(req); len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password")
	}

	// access token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	// refresh token
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	refreshUUID, err := uuid.Parse(refreshPayload.ID)
	if err != nil {
		fmt.Println("error parsing uuid:", err)
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	mtdt := server.extractMetadata(ctx)

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshUUID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	rsp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiresAt.Time),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiresAt.Time),
		User:                  convertUser(user),
	}

	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}