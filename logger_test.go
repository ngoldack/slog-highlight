package sloghighlight_test

import (
	"context"
	"log/slog"
	"testing"

	sloghighlight "github.com/ngoldack/slog-highlight"
	"github.com/stretchr/testify/require"
)

func TestNewHighlightHandler(t *testing.T) {
	handler := sloghighlight.NewHighlightHandler()
	require.NotNil(t, handler, "handler should not be nil")
}

func TestEnabled(t *testing.T) {
	ctx := context.Background()
	handler := sloghighlight.NewHighlightHandler(
		sloghighlight.WithLevel(slog.LevelWarn),
	)
	require.NotNil(t, handler, "handler should not be nil")

	levels := map[slog.Leveler]bool{
		slog.LevelDebug: false,
		slog.LevelInfo:  false,
		slog.LevelWarn:  true,
		slog.LevelError: true,
	}

	for level, enabled := range levels {
		if !enabled {
			require.False(t, handler.Enabled(ctx, level.Level()), "handler should not be enabled")
			continue
		}
		require.Equal(t, enabled, handler.Enabled(ctx, level.Level()), "handler should be enabled")
	}
}

func TestHandle(t *testing.T) {
	ctx := context.Background()
	handler := sloghighlight.NewHighlightHandler()
	require.NotNil(t, handler, "handler should not be nil")

	type userKey string
	uk := userKey("user")
	uv := struct {
		user string
	}{
		user: "test",
	}

	ctx = context.WithValue(ctx, uk, uv)

	records := []slog.Record{
		{
			Level:   slog.LevelDebug,
			Message: "Hello Debug!",
			PC:      0,
		},
		{
			Level:   slog.LevelInfo,
			Message: "Hello World!",
			PC:      0,
		},
		{
			Level:   slog.LevelWarn,
			Message: "Hello Warn!",
			PC:      0,
		},
		{
			Level:   slog.LevelError,
			Message: "Hello Error!",
			PC:      0,
		},
	}

	for _, r := range records {
		err := handler.Handle(ctx, r)
		require.NoError(t, err, "handler should not return an error")
	}
}

func TestWithGroups(t *testing.T) {
	handler := sloghighlight.NewHighlightHandler()
	handler = handler.WithGroup("test")

	require.NotNil(t, handler, "handler should not be nil")
}

func TestWithAttrs(t *testing.T) {
	handler := sloghighlight.NewHighlightHandler()
	handler = handler.WithAttrs([]slog.Attr{
		{
			Key:   "test",
			Value: slog.StringValue("test"),
		},
	})

	require.NotNil(t, handler, "handler should not be nil")
}
