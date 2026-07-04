package main

import (
	"context"
	"embed"
	"log"

	"filesweep/internal/bootstrap"
	wailsapi "filesweep/internal/transport/wails"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	container, err := bootstrap.NewContainer()
	if err != nil {
		log.Fatal(err)
	}
	api := wailsapi.NewAppAPI(container)
	err = wails.Run(&options.App{
		Title:  "FileSweep",
		Width:  1280,
		Height: 820,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: api.Startup,
		OnShutdown: func(_ context.Context) {
			_ = container.Close()
		},
		Bind: []interface{}{api},
	})
	if err != nil {
		log.Fatal(err)
	}
}
