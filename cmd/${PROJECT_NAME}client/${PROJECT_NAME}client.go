// kick:render
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal"
	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/pbhelloworld"
	"github.com/docopt/docopt-go"
	"google.golang.org/grpc"
)

const (
	defaultName = "world"
)

type Options struct {
	Connect string `docopt:"-c"`
}

func GetUsage(argv []string, version string) *Options {
	usage := `${PROJECT_NAME}client

Usage:
  ${PROJECT_NAME}client [-c <address>]

Options:
  -h --help     Show this screen.
  --version     Show version.
  -c <address>  Connect to address [default: 127.0.0.1:9090]
`

	opts, err := docopt.ParseArgs(usage, argv, version)
	if err != nil {
		panic(err)
	}
	config := &Options{}
	err = opts.Bind(config)
	if err != nil {
		panic(err)
	}
	return config
}

func main() {
	opts := GetUsage(os.Args[1:], internal.Version)
	// Set up a connection to the server.
	conn, err := grpc.Dial(opts.Connect, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pbhelloworld.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pbhelloworld.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Printf("Greeting: %s", r.GetMessage())
}
