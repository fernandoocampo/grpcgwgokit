package main

import (
	"os"

	"github.com/fernandoocampo/grpcgwgokit/internal/application"
	"go.uber.org/zap"
)

func main() {
	run()
}

func run() {
	c, err := application.New(os.Args)
	if err != nil {
		panic(err)
	}
	if err := c.Run(); err != nil {
		zap.L().Error("error", zap.Error(err))
		zap.L().Fatal("app run failed", zap.Error(err))
	}
}
