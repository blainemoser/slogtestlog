package slogtestlog

// Sometimes one uses `slog` for their application logging.
// This is something I can use for Unit Tests so that I can
// capture and review logs without publishing them elsewhere.
// Makes the unit testing a little less annoying.

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"
)

type (
	TestLog struct {
		mu       sync.Mutex
		messages []string
	}
)

func New() *TestLog {
	return &TestLog{
		mu:       sync.Mutex{},
		messages: make([]string, 0),
	}
}

func (h *TestLog) Enabled(_ context.Context, level slog.Level) bool {
	switch level {
	case slog.LevelDebug:
		fallthrough
	case slog.LevelInfo:
		fallthrough
	case slog.LevelWarn:
		fallthrough
	case slog.LevelError:
		return true
	default:
		panic("unreachable")
	}
}

func (h *TestLog) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.messages = make([]string, 0)
}

func (h *TestLog) Read() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return slices.Clone(h.messages)
}

func (h *TestLog) Handle(context context.Context, record slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	message := []string{
		fmt.Sprintf(
			"%s [%s] %s",
			record.Time.Format(time.RFC3339),
			record.Level,
			record.Message,
		),
	}
	record.Attrs(func(a slog.Attr) bool {
		message = append(message, fmt.Sprintf("%s: %s", a.Key, a.Value))
		return true
	})
	switch record.Level {
	case slog.LevelDebug:
		fallthrough
	case slog.LevelInfo:
		fallthrough
	case slog.LevelWarn:
		fallthrough
	case slog.LevelError:
		h.messages = append(h.messages, strings.Join(message, ", "))
	default:
		break
	}
	return nil
}
