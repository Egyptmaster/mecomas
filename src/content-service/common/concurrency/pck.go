package concurrency

import (
	"context"
	"os"
	"os/signal"
)

func OnError(f func() error) <-chan error {
	ch := make(chan error, 1)
	go func() {
		if err := f(); err != nil {
			ch <- err
		}
	}()
	return ch
}

func WhenDoneOrError(ctx context.Context, f func(ctx context.Context) error) error {
	onError := OnError(func() error { return f(ctx) })
	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-onError:
			return e
		}
	}
}

func CancelOnSignal(signals ...os.Signal) context.Context {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, signals...)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer close(sigs)
		for {
			select {
			case <-sigs:
				cancel()
				return
			}
		}
	}()
	return ctx
}
