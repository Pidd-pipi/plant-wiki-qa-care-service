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

func newArticleTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.CareArticle{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func articleTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func articleTestSvc(t *testing.T) *CareArticleService {
	db := newArticleTestDB(t)
	return NewCareArticleService(repository.NewCareArticleRepository(db), articleTestLogger())
}

func assert404(t *testing.T, err error) {
	t.Helper()
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404 AppError, got %v", err)
	}
}

func TestCareArticleGetMissing404(t *testing.T) {
	_, err := articleTestSvc(t).Get(999)
	assert404(t, err)
}

func TestCareArticleDeleteMissing404(t *testing.T) {
	err := articleTestSvc(t).Delete(999, 1)
	assert404(t, err)
}

func TestCareArticleUpdateMissing404(t *testing.T) {
	_, err := articleTestSvc(t).Update(999, 1, &model.CareArticle{Title: "x"})
	assert404(t, err)
}
