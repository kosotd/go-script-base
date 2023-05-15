package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Handler interface {
	Handle(ctx context.Context) error
}

func init() {
	pflag.Duration("sleep_interval", 0, "sleep interval")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
}

func Run(handler Handler) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	ticker := time.NewTicker(viper.GetDuration("sleep_interval"))
	defer ticker.Stop()

	if err := handler.Handle(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := handler.Handle(ctx); err != nil {
				return err
			}
		}
	}
}
