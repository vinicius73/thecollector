package support

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

func WithKillSignal(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		kill := make(chan os.Signal, 1)
		signal.Notify(kill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		sig := <-kill

		logger := zerolog.Ctx(ctx)
		logger.Warn().Msgf("OS Signal (%s)", sig.String())

		cancel()

		// Kill timeout
		<-time.After(time.Second * 20)
		logger.Error().Msg("Stop timeout...")
		os.Exit(1)
	}()

	return ctx, cancel
}
