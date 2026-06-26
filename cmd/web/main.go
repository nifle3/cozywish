package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/nifle3/cozywishlist/internal/config"
)

//go:embed static/*
var static embed.FS

//go:embed templates/*
var templates embed.FS

func main() {
	exitCode, err := run()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(exitCode)
	}

	slog.Warn("Application shutdown success")
}

func run() (int, error) {
	slog.Info("App start with default log. Start reading cfg")
	cfg, err := config.FromEnv()
	if err != nil {
		return 78, err
	}
	slog.Info("Read config complete")

	loggerSetup(cfg.LoggerLevel, cfg.AppEnv)

	slog.Info("logger setup with", slog.String("LEVEL", cfg.LoggerLevel))

	// TODO: db setup
	// TODO: routes setup

	fsys, err := fs.Sub(static, "static")
	if err != nil {
		return 2, err
	}

	fsServer := http.FileServerFS(fsys)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFS(templates, "templates/layout.html", "templates/index.html"))
		type PageData struct {
			Title string
		}

		data := PageData{Title: "Home"}
		err := t.ExecuteTemplate(w, "layout", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.Handle("GET /static/", http.StripPrefix("/static/", fsServer))

	err = http.ListenAndServe(":8080", nil)

	return 0, err
}

func loggerSetup(level string, development string) {
	var handler slog.Handler
	var slogLevel slog.Level

	switch strings.ToUpper(level) {
	case "INFO":
		slogLevel = slog.LevelInfo
	case "WARN":
		slogLevel = slog.LevelWarn
	case "ERROR":
		slogLevel = slog.LevelError
	case "DEBUG":
		slogLevel = slog.LevelDebug
	}

	switch strings.ToUpper(development) {
	case "PRODUCTION":
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slogLevel,
		})
	case "DEVELOPMENT":
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slogLevel,
		})
	}

	slog.SetDefault(slog.New(handler))
}
