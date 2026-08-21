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

func newPestTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil { t.Fatalf("open sqlite: %v", err) }
	if err := db.AutoMigrate(&model.DiseasePest{}); err != nil { t.Fatalf("migrate: %v", err) }
	return db
}

func pestTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func pestTestSvc(t *testing.T) *DiseasePestService {
	db := newPestTestDB(t)
	return NewDiseasePestService(repository.NewDiseasePestRepository(db), pestTestLogger())
}

func assertPest404(t *testing.T, err error) {
	t.Helper()
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404 AppError, got %v", err)
	}
}

func TestPestGet404(t *testing.T) { _, err := pestTestSvc(t).Get(999); assertPest404(t, err) }
func TestPestDelete404(t *testing.T) { assertPest404(t, pestTestSvc(t).Delete(999)) }
func TestPestUpdate404(t *testing.T) {
	_, err := pestTestSvc(t).Update(999, &model.DiseasePest{Name: "x"})
	assertPest404(t, err)
}
