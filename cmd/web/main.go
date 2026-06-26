package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/nifle3/cozywishlist/internal/config"
)

//go:embed static/*
var static embed.FS

func run() (int, error) {
	config.FromEnv()
	// logger setup
	// gracegul shutdown setup
	// db setup
	// routes setup

	fsys, err := fs.Sub(static, "static")
	if err != nil {
		return 2, err
	}

	fsServer := http.FileServerFS(fsys)
	http.Handle("/static/", http.StripPrefix("/static/", fsServer))

	err = http.ListenAndServe(":8080", nil)

	return 0, err
}

func main() {
	exitCode, err := run()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(exitCode)
	}

	slog.Warn("Application shutdown success")
}
