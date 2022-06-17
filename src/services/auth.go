package services

import (
	"crypto/rsa"
	"dh-backend-auth-sv/src/auth"
	"dh-backend-auth-sv/src/helpers"
	"dh-backend-auth-sv/src/models"
	"dh-backend-auth-sv/src/ports"
	rabbitMQ2 "dh-backend-auth-sv/src/rabbitMQ"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"strings"
	"time"
)

type Server struct {
	DB         ports.DB
	RedisCache ports.RedisCache
	auth.UnimplementedAuthServiceServer
	jwtKey *rsa.PrivateKey
}

func (s *Server) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	email := strings.TrimSpace(request.GetEmail())
	if !helpers.IsEmailValid(email) {
		helpers.LogEvent("ERROR", fmt.Sprintf("%s: %s", helpers.ErrInvalidEmail, email))
		return nil, status.Errorf(codes.InvalidArgument, "Email is not valid")
	}
	password := strings.TrimSpace(request.GetPassword())
	sevenOrMore, number, upper, special := helpers.VerifyPassword(password)
	if !sevenOrMore {
		helpers.LogEvent("ERROR", fmt.Sprintf("%s: %s", helpers.ErrPassword, password))
		return nil, status.Errorf(codes.InvalidArgument, helpers.ErrPassword)
	}
	if !number {
		helpers.LogEvent("ERROR", "Password must contain at least one number")
		return nil, status.Errorf(codes.InvalidArgument, "Password must contain at least one number")
	}
	if !upper {
		helpers.LogEvent("ERROR", "Password must contain at least one uppercase letter")
		return nil, status.Errorf(codes.InvalidArgument, "Password must contain at least one uppercase letter")
	}
	if !special {
		helpers.LogEvent("ERROR", "Password must contain at least one special character")
		return nil, status.Errorf(codes.InvalidArgument, "Password must contain at least one special character")
	}

	if len(password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "Password must be at least 8 characters long")
	}
	hashedPassword, err := helpers.GenerateHashPassword(password)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("%s: %s", helpers.ErrGenerateHashPassword, err.Error()))
	}

	err = rabbitMQ2.PublishToLoginQueue(hashedPassword, email)
	if err != nil {
		return nil, err
	}

	user := s.RedisCache.GetSubChannel(email)
	if user.Email == "" {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	if !helpers.CheckPasswordHash(password, []byte(user.HashedPassword)) {
		return nil, status.Error(codes.NotFound, "user password incorrect")
	}

	now := time.Now()
	exp := now.Add(time.Hour * 24)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"aud": "auth-service",
		"iss": "auth-service",
		"exp": exp.Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
	})
	tokenStr, err := token.SignedString([]byte(email))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	err = rabbitMQ2.PublishToRoleQueue(user.ID)
	if err != nil {
		return nil, err
	}

	var roleID []*auth.Role
	userRoles := s.RedisCache.GetRoleChannels("roles")

	for _, role := range userRoles {
		roles := &auth.Role{
			UserId: role.UserID,
			RoleId: role.RoleID,
		}
		roleID = append(roleID, roles)
	}

	//, _ := metadata.FromIncomingContext(ctx)
	activities := &models.Activities{
		ID:     uuid.New().String(),
		UserID: user.ID,
		Token:  tokenStr,
		Time:   time.Now(),
		Device: string(rune(os.Getpid())),
	}

	err = s.DB.SaveActivities(activities)
	if err != nil {
		log.Printf("err %s", err)
	}

	response := &auth.LoginResponse{
		Token: tokenStr,
		Roles: roleID,
	}
	return response, nil
}
