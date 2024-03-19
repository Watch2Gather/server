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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

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
	ID       uuid.UUID
	Username string
	Email    string
}

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

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

func CreateRefreshToken(ctx context.Context) (string, error) {
	claims := RefreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 8784).Unix(), // 8784 = 1 year
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
		return "", errInvalidToken
	}
	return CreateAccessToken(ctx, data)
}

func ParseToken(ctx context.Context, tokenString string) (AccessTokenClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, AccessTokenClaims{})
	if err != nil {
		return AccessTokenClaims{}, fmt.Errorf("jwt.ParseUnverified: %e", err)
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return AccessTokenClaims{}, errInvalidToken
	}

	return *claims, nil
}

func TokenInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	slog.Debug(fmt.Sprint("Method name: ", info.FullMethod))
	if info.FullMethod == gen.UserService_LoginUser_FullMethodName {
		return handler(ctx, req)
	}
	if info.FullMethod == gen.UserService_RegisterUser_FullMethodName {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	if len(md["authorization"]) < 1 {
		return nil, errInvalidToken
	}
	tokenString := strings.TrimPrefix(md["authorization"][0], "Bearer ")

	if !valid(tokenString) {
		return nil, errInvalidToken
	}

	m, err := handler(ctx, req)
	if err != nil {
		slog.Error("RPC failed with error: %v", err)
	}
	return m, err
}

func valid(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(accessTokenKey), nil
	})
	if err != nil {
		slog.Error(fmt.Sprintf("jwt.Parse: %e", err))
		return false
	}
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return false
	}

	now := time.Now().Unix()
	if claims.NotBefore > now {
		return false
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return false
	}

	return token.Valid
}
