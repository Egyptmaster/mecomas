package slog

import (
	"context"
	"fmt"
	"log/slog"
)

const (
	ErrorPropertyKey = "error"
)

func ErrorContext(ctx context.Context, msg string, err error, props ...slog.Attr) error {
	args := make([]any, 0, len(props)+1)
	for _, prop := range props {
		args = append(args, prop)
	}
	args = append(args, slog.String(ErrorPropertyKey, err.Error()))
	slog.ErrorContext(ctx, msg, args)
	return err
}

func InfoContextf(ctx context.Context, msg string, args ...any) {
	if len(args) == 0 {
		slog.InfoContext(ctx, msg)
	} else {
		slog.InfoContext(ctx, fmt.Sprintf(msg, args...))
	}
}
