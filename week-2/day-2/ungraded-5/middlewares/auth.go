package middlewares

import (
	"context"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
)

func JWTAuth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		fmt.Println(token)
		return nil, fmt.Errorf("Couldn't get token : %v", err)
	}
	
	claims := &jwt.StandardClaims{}
	tokenObj, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Couldn't get token : %v", err)
	}

	if !tokenObj.Valid {
		return nil, fmt.Errorf("Token is not valid")
	}

	ctx = context.WithValue(ctx, "user", claims.Subject)
	return ctx, nil
}
