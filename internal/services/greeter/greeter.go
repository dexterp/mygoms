// kick:render
package greater

import (
	"context"
	"fmt"

	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/errs"
	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/logger"
	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/pbhelloworld"
)

// Greeter GRPC server
type Greeter struct {
	err errs.LogIface
	log logger.OutputIface
	pbhelloworld.UnimplementedGreeterServer
}

// Options contructor options
type Options struct {
	Err errs.LogIface
	Log logger.OutputIface
}

// New construct microservice
func New(opts Options) *Greeter {
	return &Greeter{
		err: opts.Err,
		log: opts.Log,
	}
}

// SayHello add service
func (s *Greeter) SayHello(ctx context.Context, in *pbhelloworld.HelloRequest) (*pbhelloworld.HelloReply, error) {
	name := in.Name
	return &pbhelloworld.HelloReply{Message: fmt.Sprintf("Hello %s\n", name)}, nil
}
