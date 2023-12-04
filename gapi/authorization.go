package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/xiao-yangg/simplebank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationType = "bearer"
)

func (server *Server) authorizedUser(ctx context.Context, accessibleRole []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader) // <authorization-type> <access token>
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization format")
	}

	
	authType := strings.ToLower(fields[0])
	if authType != authorizationType {
		return nil, fmt.Errorf("unsupported authorization type")
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	if !hasPermission(payload.Role, accessibleRole) {
		return nil, fmt.Errorf("permission denied")
	}

	return payload, nil
}

func hasPermission(userRole string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if userRole == role {
			return true
		}
	}
	return false
}