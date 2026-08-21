package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/gbplantwiki/gbplantwiki/internal/config"
	"github.com/gbplantwiki/gbplantwiki/internal/router"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

func main() {
	logger := util.NewLogger()
	cfg := config.Load()

	var dialector gorm.Dialector
	if cfg.DBDriver == "sqlite" {
		dialector = sqlite.Open(cfg.DBName)
	} else {
		dialector = mysql.Open(cfg.DSN())
	}
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	if err := migrate(db); err != nil {
		logger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}
	if err := seed(db); err != nil {
		logger.Error("failed to seed database", "error", err)
		os.Exit(1)
	}

	r := router.Setup(cfg, db, logger)
	r.Static("/uploads", cfg.UploadDir)

	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("gbplantwiki backend listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", "error", err)
	}
}
