package watcher

import (
	"fmt"
	"time"

	"github.com/aottr/nox/internal/cache"
	"github.com/aottr/nox/internal/config"
	"github.com/aottr/nox/internal/logging"
	"github.com/aottr/nox/internal/processor"
)

func Start(cfg *config.Config) {
	log := logging.Get()
	logging.SetLevel("debug")
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	ctx, err := config.BuildRuntimeCtxFromConfig(cfg)
	if err != nil {
		log.Error("error building runtime context", "error", err.Error())
		return
	}
	log.Info(fmt.Sprintf("Starting watcher (interval: %s)\n", cfg.Interval))

	for {
		if err := cache.GlobalCache.RefreshCache(); err != nil {
			log.Error("error pre-fetching secrets", "error", err.Error())
		}
		if err := processor.SyncApps(ctx); err != nil {
			log.Error("error syncing secrets", "error", err.Error())
		}
		<-ticker.C
	}
}
