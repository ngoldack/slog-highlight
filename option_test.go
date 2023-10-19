package sloghighlight_test

import (
	"context"
	"testing"

	"log/slog"

	sloghighlight "github.com/ngoldack/slog-highlight"
	"github.com/stretchr/testify/require"
)

func TestWithLevel(t *testing.T) {
	ctx := context.Background()
	handler := sloghighlight.NewHighlightHandler(sloghighlight.WithLevel(slog.LevelDebug))

	require.NotNil(t, handler, "handler should not be nil")
	require.True(t, handler.Enabled(ctx, slog.LevelDebug), "handler should be enabled")
}
