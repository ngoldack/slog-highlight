package sloghighlight_test

import (
	"testing"

	sloghighlight "github.com/ngoldack/slog-highlight"
	"github.com/stretchr/testify/require"
)

func TestNewHighlightHandler(t *testing.T) {
	handler := sloghighlight.NewHighlightHandler()
	require.NotNil(t, handler, "handler should not be nil")
}
