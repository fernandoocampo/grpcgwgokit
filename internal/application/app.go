package application

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcadapter "github.com/fernandoocampo/grpcgwgokit/internal/handler"
	"github.com/fernandoocampo/grpcgwgokit/internal/service"
	pb "github.com/fernandoocampo/grpcgwgokit/pkg/proto/grpcgwgokit/pb"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

/*
https://gist.github.com/pteich/c0bb58b0b7c8af7cc6a689dd0d3d26ef
https://github.com/urfave/cli/issues/945
*/

// App defines application.
type App struct {
	args   []string
	app    *cli.App
	logger *zap.Logger
}

// New creates a new App instance
func New(args []string) (*App, error) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	newApp := &App{
		args:   args,
		app:    cli.NewApp(),
		logger: zapLogger,
	}
	return newApp, nil
}

func (a *App) setInfo() {
	a.app.Name = "Thrive's GO Reset Guides Microservice"
	a.app.Usage = "A microservice for handling reset guides"
	a.app.Version = "latest"
}

func (a *App) setFlags() {
	a.app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "Enable logging with debug mode",
			EnvVars: []string{"DEBUG"},
		},
		&cli.StringFlag{
			Name:    "grpc_server_address",
			Usage:   "gRPC server address",
			EnvVars: []string{"GRPC_SERVER_ADDRESS"},
			Value:   ":50501",
		},
		&cli.StringFlag{
			Name:    "grpc_gateway_address",
			Usage:   "gRPC gateway address",
			EnvVars: []string{"GRPC_GATEWAY_ADDRESS"},
			Value:   ":8080",
		},
	}
}

func (a *App) setDefaultAction() {
	a.app.Action = a.createDefaultAction()
}

func (a *App) createDefaultAction() cli.ActionFunc {
	return func(c *cli.Context) error {
		if c.Bool("debug") {
			zapLogger, err := zap.NewDevelopment()
			if err != nil {
				return err
			}

			a.logger = zapLogger
		}

		if c.Context == nil {
			a.logger.Info("the c.Context is nil")
		}

		ctx, cancel := context.WithCancel(c.Context)
		g, ctx := errgroup.WithContext(ctx)

		// propagate new context to other actions
		c.Context = ctx

		g.Go(func() error {
			ch := make(chan os.Signal, 1)
			// signal.Notify(ch, []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL}...)
			signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)

			select {
			// wait on kill signal
			case <-ch:
				cancel()
				os.Exit(1)
			// wait on context cancel
			case <-ctx.Done():
				return ctx.Err()
			}

			return nil
		})

		g.Go(func() error {
			return a.StartServer(c)
		})

		return g.Wait()
	}
}

// Run runs the application.
func (a *App) Run() error {
	a.setInfo()
	a.setFlags()
	a.setDefaultAction()

	return a.app.Run(a.args)
}

// StartServer starts the app server
func (a *App) StartServer(c *cli.Context) error {
	serviceUser := service.NewBasicEcho()
	// finder endpoint
	endpoints := service.NewEndpoints(serviceUser)
	handler := grpcadapter.NewGRPCServer(endpoints)
	gRPCServer := grpc.NewServer()
	pb.RegisterYourServiceServer(gRPCServer, handler)
	a.logger.Info("starting grpc server", zap.String("grpc:", c.String("grpc_server_address")))
	listener, err := net.Listen("tcp", c.String("grpc_server_address"))
	if err != nil {
		log.Println("error", err)
		return err
	}
	return gRPCServer.Serve(listener)
}
