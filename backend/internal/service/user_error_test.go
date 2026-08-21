package service

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gbplantwiki/gbplantwiki/internal/config"
	"github.com/gbplantwiki/gbplantwiki/internal/model"
	"github.com/gbplantwiki/gbplantwiki/internal/repository"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

func newUserTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func userTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func userTestSvc(t *testing.T) *UserService {
	db := newUserTestDB(t)
	return NewUserService(repository.NewUserRepository(db), userTestLogger(), config.Load())
}

func TestUserUpdateMissingNoCrash(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("UpdateProfile panicked: %v", r)
		}
	}()
	_, err := userTestSvc(t).UpdateProfile(999, "n", "", "")
	if err == nil {
		t.Fatal("expected error for missing user")
	}
}

func TestUserGetMissingNoCrash(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("GetByID panicked: %v", r)
		}
	}()
	u, err := userTestSvc(t).GetByID(999)
	if err == nil || u != nil {
		t.Fatalf("expected error and nil user, got u=%v err=%v", u, err)
	}
}

func TestUserLoginMissing401(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Login panicked: %v", r)
		}
	}()
	_, _, err := userTestSvc(t).Login("nobody", "x")
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 401 {
		t.Fatalf("expected 401 AppError, got %v", err)
	}
}
