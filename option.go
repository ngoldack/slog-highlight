package sloghighlight

import "log/slog"

type HighlightHandlerOption func(*HighlightHandler) error

func WithLevel(level slog.Leveler) HighlightHandlerOption {
	return func(h *HighlightHandler) error {
		h.level = level
		return nil
	}
}
