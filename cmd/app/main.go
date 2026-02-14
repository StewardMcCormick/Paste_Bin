package main

import (
	"context"
	"fmt"
	"github.com/StewardMcCormick/Paste_Bin/config"
	"github.com/StewardMcCormick/Paste_Bin/internal/controller/http"
	"github.com/StewardMcCormick/Paste_Bin/pkg/httpserver"
	"github.com/StewardMcCormick/Paste_Bin/pkg/logging"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	AppRun(context.Background(), cfg)
}

func AppRun(ctx context.Context, cfg *config.Config) {
	logger, err := logging.NewLogger(cfg.Logger, cfg.App.Env, cfg.App.Name, cfg.App.Version)
	if err != nil {
		panic(err)
	}

	router := http.Router(logger)

	server := httpserver.New(router, &cfg.Server)

	logger.Info(fmt.Sprintf("Server starts on %s:%s", cfg.Server.Host, cfg.Server.Port))
	err = server.Run()
	if err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig
	server.Close()
}
