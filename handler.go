package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	pb "github.com/gregory-vc/user-service/proto/user"
	"github.com/micro/go-micro/broker"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type service struct {
	repo         Repository
	tokenService Authable
	PubSub       broker.Broker
}

func (srv *service) GetUser(ctx context.Context, req *pb.ID, res *pb.User) error {
	user, err := srv.repo.Get(req.Id)
	if err != nil {
		log.Println(err)
	}
	res = user
	return nil
}

func (srv *service) UpdateUser(ctx context.Context, req *pb.User, res *pb.User) error {
	user, err := srv.repo.Get(req.Id)
	if err != nil {
		log.Println(err)
	}
	res = user
	return nil
}

func (srv *service) DeleteUser(ctx context.Context, req *pb.ID, res *pb.ID) error {
	_, err := srv.repo.Get(req.Id)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (srv *service) ListUsers(ctx context.Context, req *pb.ListUsersRequest, res *pb.ListUsersResponse) error {
	users, err := srv.repo.GetAll()
	if err != nil {
		log.Println(err)
	}
	res.Users = users
	res.Count = uint32(len(users))
	return nil
}

func (srv *service) ListUsersByID(ctx context.Context, req *pb.ListUsersByIDRequest, res *pb.ListUsersResponse) error {
	users, err := srv.repo.GetAll()
	if err != nil {
		log.Println(err)
	}
	res.Users = users
	res.Count = uint32(len(users))
	return nil
}

func (srv *service) AuthUser(ctx context.Context, req *pb.AuthUserRequest, res *pb.AuthUserResponse) error {
	log.Println("Logging in with:", req.Email, req.Password)
	user, err := srv.repo.GetByEmail(req.Email)
	if err != nil {
		log.Println(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Println(err)
	}

	token, err := srv.tokenService.Encode(user)
	if err != nil {
		log.Println(err)
	}

	res.Jwt = token
	return nil
}

func (srv *service) CreateUser(ctx context.Context, req *pb.User, res *pb.User) error {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	req.Password = string(hashedPass)
	if err := srv.repo.Create(req); err != nil {
		log.Println(err)
	}
	*res = *req

	fmt.Println(res)

	if err := srv.publishEvent(req); err != nil {
		log.Println(err)
	}
	return nil
}

func (srv *service) publishEvent(user *pb.User) error {
	// Marshal to JSON string
	body, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
	}

	// Create a broker message
	msg := &broker.Message{
		Header: map[string]string{
			"id": string(user.Id),
		},
		Body: body,
	}

	// Publish message to broker
	if err := srv.PubSub.Publish(topic, msg); err != nil {
		log.Printf("[pub] failed: %v", err)
	}

	return nil
}

func (srv *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {

	// Decode token
	claims, err := srv.tokenService.Decode(req.Token)

	if err != nil {
		log.Println(err)
	}

	if claims.User.Id == 0 {
		log.Println(errors.New("invalid user"))
	}

	res.Valid = true

	return nil
}
