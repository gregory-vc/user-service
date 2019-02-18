package main

import (
	"fmt"
	"log"
	"os"

	pb "github.com/gregory-vc/user-service/proto/user"
	"github.com/micro/cli"
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

		micro.Flags(cli.BoolFlag{
			Name:  "migrate",
			Usage: "Launch the migration",
		}),
	)

	srv.Init(
		micro.Action(func(c *cli.Context) {
			if c.Bool("migrate") {
				Migrate()
				os.Exit(0)
			}
		}),
	)

	pubsub := srv.Server().Options().Broker
	pb.RegisterUserServiceHandler(srv.Server(), &service{repo, tokenService, pubsub})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
