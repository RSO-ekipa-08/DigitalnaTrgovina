package auth

import (
	"context"

	pb "authentication/src/gen/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client pb.AuthServiceClient
}

func NewAuthClient(miniAuthAddr string) (*AuthClient, error) {
	conn, err := grpc.Dial(miniAuthAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewAuthServiceClient(conn)
	return &AuthClient{client: client}, nil
}

func (ac *AuthClient) Login(ctx context.Context, redirectURL string) (string, error) {
	resp, err := ac.client.Login(ctx, &pb.LoginRequest{
		RedirectUrl: redirectURL,
	})
	if err != nil {
		return "", err
	}
	return resp.AuthUrl, nil
}

func (ac *AuthClient) Verify(ctx context.Context, code string) (*pb.VerifyResponse, error) {
	return ac.client.Verify(ctx, &pb.VerifyRequest{
		Code: code,
	})
}

func (ac *AuthClient) Logout(ctx context.Context, returnURL string) (string, error) {
	resp, err := ac.client.Logout(ctx, &pb.LogoutRequest{
		ReturnUrl: returnURL,
	})
	if err != nil {
		return "", err
	}
	return resp.LogoutUrl, nil
}

func (ac *AuthClient) VerifyToken(ctx context.Context, token string) (*pb.VerifyTokenResponse, error) {
	return ac.client.VerifyToken(ctx, &pb.VerifyTokenRequest{
		Token: token,
	})
}
