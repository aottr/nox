package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aottr/nox/internal/config"
	"github.com/aottr/nox/internal/processor"
)

func Run(ctx context.Context, rt *config.RuntimeContext, interval time.Duration) error {
	logger := rt.Logger
	logger.Println("🌀 Nox daemon starting...")

	// Set up signal handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Use a cancellable context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Shutdown listener
	go func() {
		sig := <-sigs
		logger.Printf("⚠️ Received signal: %v, shutting down...\n", sig)
		cancel()
	}()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Initial run
	if err := processor.SyncApps(rt); err != nil {
		logger.Printf("❌ Initial run failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Println("👋 Nox daemon exiting.")
			return nil
		case <-ticker.C:
			logger.Println("🔁 Processing apps...")
			if err := processor.SyncApps(rt); err != nil {
				logger.Printf("❌ Error: %v", err)
			}
		}
	}
}
