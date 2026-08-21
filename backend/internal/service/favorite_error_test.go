package service

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gbplantwiki/gbplantwiki/internal/model"
	"github.com/gbplantwiki/gbplantwiki/internal/repository"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

func newFavoriteTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.Favorite{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func favoriteTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func favoriteTestSvc(t *testing.T) *FavoriteService {
	db := newFavoriteTestDB(t)
	return NewFavoriteService(repository.NewFavoriteRepository(db), favoriteTestLogger())
}

func TestFavoriteRemoveNotFound404(t *testing.T) {
	err := favoriteTestSvc(t).Remove(1, "plant", 999)
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404 AppError, got %v", err)
	}
}

func TestFavoriteAddDuplicate409(t *testing.T) {
	svc := favoriteTestSvc(t)
	if _, err := svc.Add(1, "plant", 7); err != nil {
		t.Fatalf("first add: %v", err)
	}
	_, err := svc.Add(1, "plant", 7)
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 409 {
		t.Fatalf("expected 409 AppError, got %v", err)
	}
}
