package sharedkernel

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/Watch2Gather/server/proto/gen"
)

var (
	accessTokenKey  = os.Getenv("ACCESS_TOKEN_KEY")
	refreshTokenKey = os.Getenv("REFRESH_TOKEN_KEY")
)

type AccessTokenClaims struct {
	Username string
	Email    string
	jwt.StandardClaims
}
type RefreshTokenClaims struct {
	jwt.StandardClaims
}

type UserData struct {
	Username string
	Email    string
	ID       uuid.UUID
}

func CreateAccessToken(ctx context.Context, data UserData) (string, error) {
	claims := AccessTokenClaims{
		Username: data.Username,
		Email:    data.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Add(time.Minute * -5).Unix(),
			Id:        data.ID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(accessTokenKey))
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %e", err)
	}
	return tokenString, nil
}

func CreateRefreshToken(ctx context.Context, id uuid.UUID) (string, error) {
	claims := RefreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 8784).Unix(),
			Id:        id.String(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Add(time.Minute * -5).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(refreshTokenKey))
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %e", err)
	}
	return tokenString, nil
}

func RefreshAccessToken(ctx context.Context, tokenString string, data UserData) (string, error) {
	if !valid(tokenString) {
		return "", ErrInvalidToken
	}
	return CreateAccessToken(ctx, data)
}

func ParseToken(ctx context.Context, tokenString string) (AccessTokenClaims, error) {
	claims := AccessTokenClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, &claims)
	if err != nil {
		return AccessTokenClaims{}, fmt.Errorf("jwt.ParseUnverified: %e", err)
	}

	return claims, nil
}

func GetToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMissingMetadata
	}

	if len(md["authorization"]) < 1 {
		return "", ErrInvalidToken
	}
	tokenString := strings.TrimPrefix(md["authorization"][0], "Bearer ")

	return tokenString, nil
}

func TokenInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	slog.Debug(fmt.Sprint("Method name: ", info.FullMethod))
	if info.FullMethod == gen.UserService_LoginUser_FullMethodName {
		return handler(ctx, req)
	}
	if info.FullMethod == gen.UserService_RegisterUser_FullMethodName {
		return handler(ctx, req)
	}

	tokenString, err := GetToken(ctx)
	if err != nil {
		slog.Error("token.GetToken: %e", err)
		return nil, err
	}

	if !valid(tokenString) {
		return nil, ErrInvalidToken
	}

	m, err := handler(ctx, req)
	if err != nil {
		slog.Error("RPC failed with error: %v", err)
		return nil, err
	}
	return m, nil
}

func valid(tokenString string) bool {
	claims := AccessTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessTokenKey), nil
	})
	if err != nil {
		slog.Error(fmt.Sprintf("jwt.Parse: %e", err))
		return false
	}

	slog.Debug("valid", "token", token)

	return token.Valid
}
