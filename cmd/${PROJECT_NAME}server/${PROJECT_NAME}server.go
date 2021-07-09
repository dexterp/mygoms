// kick:render
package main

import (
	"os"

	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal"
	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/di"
	"github.com/docopt/docopt-go"
)

type Options struct {
	Listen string `docopt:"-l"`
	Xds    bool   `docopt:"--xds"`
}

func GetUsage(argv []string, version string) *Options {
	usage := `${PROJECT_NAME}server

Usage:
  ${PROJECT_NAME}server [-l <listen>] [--xds]

Options:
  -h --help     show this screen
  --version     show version
  -l <listen>   listen on host:port [default: 127.0.0.1:9090]
  --xds         whether the server should use xDS APIs to receive security configuration
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
	cli := GetUsage(os.Args[1:], internal.Version)
	inject := di.New(di.Options{
		Listen: cli.Listen,
	})

	switch {
	case cli.Xds:
		startServer := inject.CallXDSServer()
		startServer()
	default:
		startServer := inject.CallGRPCServer()
		startServer()
	}
}
