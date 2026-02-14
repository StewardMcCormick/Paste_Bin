package main

import (
	"context"
	"github.com/StewardMcCormick/Paste_Bin/config"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	RunApp(context.Background(), cfg)
}

func RunApp(ctx context.Context, cfg *config.Config) {

}
