package main

import (
	"errors"
	"net/http"
	"os"
	"runtime/debug"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/highlight/highlight/sdk/highlight-go"
	highlightChi "github.com/highlight/highlight/sdk/highlight-go/middleware/chi"
	"github.com/joho/godotenv"
	sloghighlight "github.com/ngoldack/slog-highlight"
	slogmulti "github.com/samber/slog-multi"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	if _, ok := os.LookupEnv("HIGHLIGHT_PROJECT_ID"); !ok {
		panic("env HIGHLIGHT_PROJECT_ID not set")
	}

	info, ok := debug.ReadBuildInfo()
	v := "local-dev"
	if ok {
		v = info.Main.Sum
	}
	highlight.SetProjectID(os.Getenv("HIGHLIGHT_PROJECT_ID"))
	highlight.Start(
		highlight.WithServiceName("slog-highlight-example-api"),
		highlight.WithServiceVersion(v),
	)
	defer highlight.Stop()

	// to use multiple handlers, use slogmulti.Fanout
	logger := slog.New(slogmulti.Fanout(
		sloghighlight.NewHighlightHandler(),
		slog.NewTextHandler(os.Stderr, nil),
	))
	slog.SetDefault(logger)

	r := chi.NewRouter()
	r.Use(highlightChi.Middleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "Index called",
			"method", r.Method,
			"host", r.Host,
			"origin", r.RemoteAddr,
		)
		slog.ErrorContext(r.Context(), "Random error", "error", errors.New("random error"))
		w.Write([]byte("Hello World!"))
	})

	slog.Info("starting server", "addr", "http://localhost:3000")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
