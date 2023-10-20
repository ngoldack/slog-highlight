package main

import (
	"context"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/highlight/highlight/sdk/highlight-go"
	"github.com/joho/godotenv"
	sloghighlight "github.com/ngoldack/slog-highlight"
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
		highlight.WithServiceName("slog-highlight-example-minimal"),
		highlight.WithServiceVersion(v),
	)
	defer highlight.Stop()

	// this logger only uses the highlight handler
	// you will not see any output in the console, only in highlight
	logger := slog.New(
		sloghighlight.NewHighlightHandler(),
	)
	slog.SetDefault(logger)

	ctx := context.Background()

	type userKey string
	uk := userKey("user")
	uv := struct {
		user string
	}{
		user: "test",
	}

	ctx = context.WithValue(ctx, uk, uv)

	slog.Debug("Hello Debug!")
	slog.DebugContext(ctx, "Hello Debug Context!")
	slog.Info("Hello World!")
	slog.InfoContext(ctx, "Hello World Context!")
	slog.Warn("Hello Warn!")
	slog.WarnContext(ctx, "Hello Warn Context!")
	slog.Error("Hello Error!")
	slog.ErrorContext(ctx, "Hello Error Context!")
}
