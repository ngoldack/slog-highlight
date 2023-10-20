package sloghighlight

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"slices"

	"github.com/highlight/highlight/sdk/highlight-go"
	hlog "github.com/highlight/highlight/sdk/highlight-go/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type HighlightHandler struct {
	level slog.Leveler

	attrs  []slog.Attr
	groups []string
}

var _ slog.Handler = (*HighlightHandler)(nil)

// NewHighlightHandler creates a new HighlightHandler.
// Default level is LevelInfo.
func NewHighlightHandler(options ...HighlightHandlerOption) slog.Handler {
	h := &HighlightHandler{
		level: slog.LevelInfo,
	}
	for _, option := range options {
		option(h)
	}
	return h
}

// Enabled returns true if the level is enabled.
func (h *HighlightHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.level.Level() <= level.Level()
}

// Handle handles the record.
func (h *HighlightHandler) Handle(ctx context.Context, r slog.Record) error {
	span, _ := highlight.StartTrace(ctx, "highlight-go/log")
	defer highlight.EndTrace(span)

	// add default attributes
	attrs := []attribute.KeyValue{
		hlog.LogSeverityKey.String(r.Level.String()),
		hlog.LogMessageKey.String(r.Message),
	}

	// if PC is not nil, get caller function info
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()

		source := attribute.Key("source").String(fmt.Sprintf("%s:%d (%s)", f.File, f.Line, f.Function))
		attrs = append(attrs, source)
	}

	r.Attrs(func(attr slog.Attr) bool {
		key := attribute.Key(attr.Key)

		// TODO: correctly type infer the values
		attrs = append(attrs, key.String(fmt.Sprintf("%v", attr.Value)))
		return true
	})

	attrs = append(attrs, attribute.Key("groups").StringSlice(h.groups))

	span.AddEvent(highlight.LogEvent, trace.WithAttributes(attrs...), trace.WithTimestamp(r.Time))

	if r.Level <= slog.LevelError {
		span.SetStatus(codes.Error, r.Message)
	}

	return nil
}

// WithLevel returns a new HighlightHandler with the attributes added.
func (h *HighlightHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := append(slices.Clone[[]slog.Attr](attrs), h.attrs...)

	return &HighlightHandler{
		level:  h.level,
		attrs:  newAttrs,
		groups: slices.Clone[[]string](h.groups),
	}
}

// WithGroup returns a new HighlightHandler with the group added.
func (h *HighlightHandler) WithGroup(name string) slog.Handler {
	newGroups := append(slices.Clone[[]string](h.groups), name)

	return &HighlightHandler{
		level:  h.level,
		attrs:  slices.Clone[[]slog.Attr](h.attrs),
		groups: newGroups,
	}
}
