package main

import (
	"fmt"
	"log"
	"os"

	pb "github.com/gregory-vc/user-service/proto/user"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/nats"
	k8s "github.com/micro/kubernetes/go/micro"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"

	bkr "github.com/micro/go-plugins/broker/grpc"
	cli "github.com/micro/go-plugins/client/grpc"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/selector/static"
	srv "github.com/micro/go-plugins/server/grpc"
	"github.com/micro/micro/api"
)

const topic = "user.created"

func main() {

	// disable namespace
	api.Namespace = ""

	// set values for registry/selector
	os.Setenv("MICRO_REGISTRY", "kubernetes")
	os.Setenv("MICRO_SELECTOR", "static")

	// setup broker/client/server
	broker.DefaultBroker = bkr.NewBroker()
	client.DefaultClient = cli.NewClient()
	server.DefaultServer = srv.NewServer()

	// Creates a database connection and handles
	// closing it again before exit.
	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	// Automatically migrates the user struct
	// into database columns/types etc. This will
	// check for changes and migrate them each time
	// this service is restarted.
	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}

	tokenService := &TokenService{repo}

	// Create a new service. Optionally include some options here.
	srv := k8s.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("user"),
		micro.Version("latest"),
	)

	// Init will parse the command line flags.
	srv.Init()

	pubsub := srv.Server().Options().Broker

	// Register handler
	pb.RegisterUserServiceHandler(srv.Server(), &service{repo, tokenService, pubsub})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
