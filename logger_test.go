package sloghighlight_test

import (
	"context"
	"log/slog"
	"math/rand"
	"runtime"
	"testing"
	"time"

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

	pcs := createPCs(t)
	records := make([]slog.Record, 0, len(pcs))
	for _, pc := range pcs {
		r := slog.NewRecord(time.Now(), slog.LevelDebug, "Hello Debug!", pc)
		r.AddAttrs(slog.Attr{
			Key:   "test",
			Value: slog.StringValue("test"),
		})

		switch rand.Intn(4) % 4 {
		case 0:
			r.Level = slog.LevelDebug
			records = append(records, r)
			continue
		case 1:
			r.Level = slog.LevelInfo
			records = append(records, r)
			continue
		case 2:
			r.Level = slog.LevelWarn
			records = append(records, r)
			continue
		case 3:
			r.Level = slog.LevelError
			records = append(records, r)
			continue
		}
	}

	ctx = context.WithValue(ctx, uk, uv)
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

func createPCs(t *testing.T) []uintptr {
	pc := make([]uintptr, 100)
	n := runtime.Callers(0, pc)
	if n == 0 {
		t.Fatal("could not get PC")
	}
	return pc[:n]
}
