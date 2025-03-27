package slogtestlog

import (
	"log/slog"
	"testing"
)

func TestNew(t *testing.T) {
	l := New()
	slog.SetDefault(slog.New(l))
}
