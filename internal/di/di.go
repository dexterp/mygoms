// kick:render
package di

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/errs"
	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/exit"
	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/logger"
	srvc "${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/services/greeter"
	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/pbhelloworld"
	opentracing "github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"
	"google.golang.org/grpc/xds"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/uber/jaeger-client-go"
)

// DI dependency injection struct
type DI struct {
	exitMode int
	listen   string
	logLvl   logger.Level
	stderr   io.Writer
	stdout   io.Writer
	trace    bool
}

// Options constructor options
type Options struct {
	ExitMode int          // Exit mode one of exit.MNone or exit.Panic
	Listen   string       // Listen host:port
	LogLevel logger.Level // Logger level
	Stderr   io.Writer    // Inject Stderr
	Stdout   io.Writer    // Inject Stdout
	Trace    bool         // Enable microservice tracing
}

// New create IoC (DI) container
func New(opts Options) *DI {
	di := &DI{
		exitMode: opts.ExitMode,
		listen:   opts.Listen,
		logLvl:   opts.LogLevel,
		stderr:   opts.Stderr,
		stdout:   opts.Stdout,
		trace:    opts.Trace,
	}
	if di.stderr == nil {
		di.stderr = os.Stderr

	}
	if di.stdout == nil {
		di.stdout = os.Stdin
	}
	return di
}

//
// Microservices
//

// CallGRPCServer start GRPC server
func (i *DI) CallGRPCServer() func() {
	logger := i.MakeLoggerOutput("", "")
	chk := i.MakeErrorHandler()
	var (
		opts []grpc.ServerOption
	)
	if i.trace {
		tracer := i.MakeTracer()
		opts = append(opts, grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
		opts = append(opts, grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(tracer)))
	}
	fn := func() {
		_, closer := i.MakeTracerCloser()
		defer func() {
			closer.Close()
		}()
		lis, err := net.Listen("tcp", i.listen)
		chk.Fatalf(`failed to listen "%s": %v`, i.listen, err)
		srv := i.MakeGreeter()
		s := grpc.NewServer()
		pbhelloworld.RegisterGreeterServer(s, srv)
		logger.Printf("server listening at %s", lis.Addr())
		err = s.Serve(lis)
		chk.Fatalf(`failed to serve "%s": %v`, lis.Addr(), err)
	}
	return fn
}

// CallXDSServer start XDS transport server
func (i *DI) CallXDSServer() func() {
	logger := i.MakeLoggerOutput("", "")
	chk := i.MakeErrorHandler()
	var (
		opts []grpc.ServerOption
	)
	if i.trace {
		tracer := i.MakeTracer()
		opts = append(opts, grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
		opts = append(opts, grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(tracer)))
	}
	creds, err := xdscreds.NewServerCredentials(xdscreds.ServerOptions{FallbackCreds: insecure.NewCredentials()})
	chk.Fatalf(`failed to create server-side xDS credentials: %v`, err)
	opts = append(opts, grpc.Creds(creds))

	fn := func() {
		lis, err := net.Listen("tcp", i.listen)
		chk.Fatalf(`failed to listen "%s": %v`, lis, err)
		s := xds.NewGRPCServer(opts...)
		srv := i.MakeGreeter()
		pbhelloworld.RegisterGreeterServer(s, srv)
		logger.Printf("server listening at %s", lis.Addr())
		err = s.Serve(lis)
		chk.Fatalf(`failed to serve "%s": %v`, lis.Addr(), err)
	}
	return fn
}

var cacheServer *srvc.Greeter

// MakeGreeter create greater microservice
func (i *DI) MakeGreeter() *srvc.Greeter {
	if cacheServer != nil {
		return cacheServer
	}
	cacheServer = srvc.New(srvc.Options{
		Log: i.MakeLoggerOutput("", ""),
	})

	return cacheServer
}

var cacheTracer *opentracing.Tracer
var cacheTraceCloser io.Closer

// MakeTracer create a tracer
func (i *DI) MakeTracer() opentracing.Tracer {
	if cacheTracer != nil {
		return *cacheTracer
	}
	// cacheTracer global is set in MakeTracerCloser
	_, _ = i.MakeTracerCloser()
	return *cacheTracer
}

// MakeTracerCloser make a tracer and a way to close
func (i *DI) MakeTracerCloser() (opentracing.Tracer, io.Closer) {
	if cacheTracer != nil {
		return *cacheTracer, cacheTraceCloser
	}
	cfg := i.MakeTraceConfig()
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		panic(err)
	}
	cacheTracer := &tracer
	cacheTraceCloser := closer
	return *cacheTracer, cacheTraceCloser
}

var traceConfig *jaegercfg.Configuration

func (i *DI) MakeTraceConfig() *jaegercfg.Configuration {
	if traceConfig != nil {
		return traceConfig
	}
	traceConfig := &jaegercfg.Configuration{
		ServiceName: "your_service_name",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	return traceConfig
}

//
// Other
//

var cacheErrHandler *errs.Handler

// MakeErrorHandler dependency injector
func (i *DI) MakeErrorHandler() *errs.Handler {
	if cacheErrHandler != nil {
		return cacheErrHandler
	}
	handler := errs.New(errs.Options{
		ExitHandler: i.MakeExitHandler(),
		Logger:      i.MakeLoggerOutput("", ""),
	})
	cacheErrHandler = handler
	return handler
}

var cacheExitHandler *exit.Handler

// MakeExitHandler dependency injector
func (i *DI) MakeExitHandler() *exit.Handler {
	if cacheExitHandler != nil {
		return cacheExitHandler
	}
	handler := exit.New(exit.Options{
		Mode: i.exitMode,
	})
	cacheExitHandler = handler
	return cacheExitHandler
}

var cacheLogFile *os.File

// MakeLogFile create a logfile and return the interface
func (i *DI) MakeLogFile(logfile string) *os.File {
	chk := i.MakeErrorHandler()
	if cacheLogFile != nil {
		return cacheLogFile
	}
	var (
		f   *os.File
		err error
	)

	fInfo, err := os.Stat(logfile)
	if err != nil && !os.IsNotExist(err) {
		// Simple output because logging is not available
		fmt.Printf(`can not open log file %s: %v`, logfile, err)
	} else if err == nil && fInfo.Size() > 1024*1024*2 {
		// Remove files greater than 2M
		os.Remove(logfile)
	}
	f, err = os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		chk.Fatalf(`can not open logfile %s: %v`, logfile, err)
	}
	cacheLogFile = f
	return cacheLogFile
}

// MakeLoggerOutput inject logger.OutputIface.
func (i *DI) MakeLoggerOutput(prefix string, logfile string) *logger.Router {
	toStderr := logger.New(i.stderr, prefix, log.Lmsgprefix, i.logLvl, i.MakeExitHandler())
	if logfile != "" {
		toFile := logger.New(i.MakeLogFile(logfile), prefix, log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix, i.logLvl, i.MakeExitHandler())
		return logger.NewRouter(toFile, toStderr)
	}
	return logger.NewRouter(toStderr)
}
