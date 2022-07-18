package services

import (
	"context"
	"crypto/rsa"
	"dh-backend-auth-sv/internal/helpers"
	"dh-backend-auth-sv/internal/models"
	"dh-backend-auth-sv/internal/ports"
	"dh-backend-auth-sv/internal/proto"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
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
	proto.UnimplementedAuthServiceServer
	jwtKey      *rsa.PrivateKey
	UserService proto.UserServiceClient
}

func (s *Server) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
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

	fmt.Println(hashedPassword)

	userRequest := proto.GetUserDetailsByEmailRequest{
		Email: email,
	}
	res, err := s.UserService.GetUserDetailsByEmail(context.Background(), &userRequest)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("user with this email does not exist!"))
		return nil, status.Errorf(codes.NotFound, "user with this email does not exist!")
	}

	user := &models.User{}

	err = json.Unmarshal(res.GetResponse(), user)
	if err != nil {
		fmt.Println(err)
		helpers.LogEvent("ERROR", fmt.Sprintf("cannot unmarshal user"))
		return nil, status.Errorf(codes.Internal, "cannot process user info")
	}

	fmt.Println(user.Role.Title)
	//userRoleRequest := proto.GetUserRolesRequest{}
	//userRoles := s.UserService.GetUserRoles(ctx)

	//err = rabbitMQ2.PublishToLoginQueue(hashedPassword, email)
	//if err != nil {
	//	return nil, err
	//}

	//user := s.RedisCache.GetSubChannel(email)
	//if user.Email == "" {
	//	return nil, status.Error(codes.NotFound, "User not found")
	//}

	if !helpers.CheckPasswordHash(password, []byte(user.HashedPassword)) {
		return nil, status.Error(codes.NotFound, "user password incorrect")
	}

	now := time.Now()
	exp := now.Add(time.Hour * 24)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": map[string]string{
			"userId": user.ID,
			"roleId": user.RoleID,
		},
		"aud": "proto-service",
		"iss": "proto-service",
		"exp": exp.Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
	})
	tokenStr, err := token.SignedString([]byte(email))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	//
	//err = rabbitMQ2.PublishToRoleQueue(user.ID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var roleID []*proto.Role
	//userRoles := s.RedisCache.GetRoleChannels("roles")
	//
	//for _, role := range userRoles {
	//	roles := &proto.Role{
	//		UserId: role.UserID,
	//		RoleId: role.RoleID,
	//	}
	//	roleID = append(roleID, roles)
	//}

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

	//var roleID []*proto.Roles
	//roles := &proto.Roles{
	//	UserId: user.ID,
	//	RoleId: user,
	//}
	//
	//var userRole proto.Role
	//
	//fmt.Println(user.Roles)
	//
	//for _, role := range user.Roles {
	//	fmt.Println(role.ID, role.Title)
	//	userRoles = append(userRoles, &proto.Role{
	//		RoleId: role.ID,
	//		Title:  role.Title,
	//	})
	//}

	response := &proto.LoginResponse{
		Token: tokenStr,
	}
	return response, nil
}
