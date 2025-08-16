package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

type CliHandler struct {
	w     io.Writer
	level slog.Leveler
}

func (h *CliHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.level != nil {
		minLevel = h.level.Level()
	}
	return level >= minLevel
}

func NewCliHandler(w io.Writer, level slog.Leveler) slog.Handler {
	h := &CliHandler{w: w, level: level}
	return h
}

func (h *CliHandler) Handle(_ context.Context, r slog.Record) error {

	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "error" {
			attrs += fmt.Sprintf("%v ", a.Value)
			return true
		}
		return false
	})

	if attrs != "" {
		fmt.Fprintf(h.w, "%s: %s\n", r.Message, attrs[:len(attrs)-1])
	} else {
		fmt.Fprintln(h.w, r.Message)
	}
	return nil
}

func (h *CliHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *CliHandler) WithGroup(name string) slog.Handler       { return h }
