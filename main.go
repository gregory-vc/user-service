package main

import (
	"fmt"
	"log"

	pb "github.com/gregory-vc/user-service/proto/user"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	k8s "github.com/micro/kubernetes/go/micro"
)

const topic = "user.created"

func main() {

	bucket, err := CreateGlobalBucket()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	repo := &UserRepository{bucket}
	tokenService := &TokenService{repo}

	srv := k8s.NewService(
		micro.Name("user"),
		micro.Version("latest"),
	)

	srv.Init()
	pubsub := srv.Server().Options().Broker
	pb.RegisterUserServiceHandler(srv.Server(), &service{repo, tokenService, pubsub})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
