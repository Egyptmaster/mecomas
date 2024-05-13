package retry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type Settings struct {
	Retry uint8
	Delay time.Duration
}

func Execute(ctx context.Context, stx Settings, f func() error) error {
	return execute(ctx, stx.Retry, stx.Delay, f)
}

func ExecuteAndReturn[T any](ctx context.Context, stx Settings, f func() (T, error)) (res T, err error) {
	cb := func() error {
		res, err = f()
		return err
	}
	return res, execute(ctx, stx.Retry, stx.Delay, cb)
}

func execute(ctx context.Context, retry uint8, delay time.Duration, f func() error) error {
	var (
		err error
		try uint8 = 0
	)
	for {
		if err = invoke(f); err == nil {
			return nil
		}
		if try == retry {
			return err
		}
		try++
		slog.WarnContext(ctx, fmt.Sprintf("%d. try failed. Will wait %s and try again", try, delay), slog.String("err", err.Error()))
		if err = wait(ctx, delay); err != nil {
			return err
		}
	}
}

func invoke(f func() error) (err error) {
	defer func() {
		if re := recover(); re != nil {
			err = errors.New(re.(string))
		}
	}()
	return f()
}

func wait(ctx context.Context, delay time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			return nil
		}
	}
}
