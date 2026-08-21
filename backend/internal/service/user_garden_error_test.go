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

func newGardenTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.UserGarden{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func gardenTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func gardenTestSvc(t *testing.T) *UserGardenService {
	db := newGardenTestDB(t)
	return NewUserGardenService(repository.NewUserGardenRepository(db), gardenTestLogger())
}

func TestUserGardenRemoveMissingNoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Remove panicked on missing item: %v", r)
		}
	}()
	err := gardenTestSvc(t).Remove(1, 999)
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404 AppError, got %v", err)
	}
}

func TestUserGardenBindMissingNoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BindReminder panicked on missing item: %v", r)
		}
	}()
	_, err := gardenTestSvc(t).BindReminder(1, 999, 1)
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404 AppError, got %v", err)
	}
}
