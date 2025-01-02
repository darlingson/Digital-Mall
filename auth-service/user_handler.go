package main

import (
	"context"
	"errors"
	"log"

	pb "digital-mall/pkg/proto/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	db     *mongo.Collection
	config Config
}

func (s *AuthService) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.AuthResponse, error) {
	// Check if user already exists
	filter := bson.M{"username": req.Username}
	var existingUser bson.M
	err := s.db.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Insert new user into MongoDB
	_, err = s.db.InsertOne(ctx, bson.M{
		"username": req.Username,
		"password": string(hashedPassword),
	})
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate JWT token
	token, err := GenerateToken(req.Username, s.config.SecretKey)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	log.Printf("New user created: %s", req.Username)
	return &pb.AuthResponse{Token: token}, nil
}
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	// Find user in MongoDB
	filter := bson.M{"username": req.Username}
	var user bson.M
	err := s.db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user["password"].(string)), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := GenerateToken(req.Username, s.config.SecretKey)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	log.Printf("User logged in: %s", req.Username)
	return &pb.AuthResponse{Token: token}, nil
}
