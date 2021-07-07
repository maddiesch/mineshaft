package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func installSignalHandlers(_ context.Context) <-chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	return sigChan
}
