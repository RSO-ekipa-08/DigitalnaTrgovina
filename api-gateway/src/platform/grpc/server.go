package grpc

import (
	pb "authentication/src/gen/proto"
	"authentication/src/platform/authenticator"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	auth *authenticator.Authenticator
}

func NewServer(auth *authenticator.Authenticator) *Server {
	return &Server{auth: auth}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	state, err := generateRandomState()
	if err != nil {
		return nil, err
	}

	authURL := s.auth.AuthCodeURL(state)
	return &pb.LoginResponse{AuthUrl: authURL}, nil
}

func (s *Server) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	token, err := s.auth.Exchange(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	idToken, err := s.auth.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, err
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return nil, err
	}

	profileJSON, err := json.Marshal(profile)
	if err != nil {
		return nil, err
	}

	return &pb.VerifyResponse{
		AccessToken: token.AccessToken,
		IdToken:     token.Extra("id_token").(string),
		Profile:     string(profileJSON),
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	logoutURL := fmt.Sprintf("https://%s/v2/logout?returnTo=%s&client_id=%s",
		os.Getenv("AUTH0_DOMAIN"),
		url.QueryEscape(req.ReturnUrl),
		os.Getenv("AUTH0_CLIENT_ID"))

	return &pb.LogoutResponse{LogoutUrl: logoutURL}, nil
}

func (s *Server) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	// Parse and verify the token
	token, err := s.auth.Provider.Verifier(&oidc.Config{
		ClientID: s.auth.ClientID,
	}).Verify(ctx, req.Token)

	if err != nil {
		return &pb.VerifyTokenResponse{
			IsValid: false,
		}, nil
	}

	// Extract claims
	var claims map[string]interface{}
	if err := token.Claims(&claims); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse claims: %v", err)
	}

	// Convert claims to string map
	stringClaims := make(map[string]string)
	for k, v := range claims {
		if str, ok := v.(string); ok {
			stringClaims[k] = str
		} else {
			// Convert non-string values to string representation
			stringClaims[k] = fmt.Sprintf("%v", v)
		}
	}

	return &pb.VerifyTokenResponse{
		IsValid: true,
		UserId:  claims["sub"].(string),
		Claims:  stringClaims,
	}, nil
}
