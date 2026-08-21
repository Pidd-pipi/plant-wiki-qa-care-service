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

func newPlantTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil { t.Fatalf("open sqlite: %v", err) }
	if err := db.AutoMigrate(&model.PlantSpecies{}); err != nil { t.Fatalf("migrate: %v", err) }
	return db
}

func plantTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func plantTestSvc(t *testing.T) *PlantSpeciesService {
	db := newPlantTestDB(t)
	return NewPlantSpeciesService(repository.NewPlantSpeciesRepository(db), plantTestLogger())
}

func assertPlant404(t *testing.T, err error) {
	t.Helper()
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404 AppError, got %v", err)
	}
}

func TestPlantGet404(t *testing.T) { _, err := plantTestSvc(t).Get(999); assertPlant404(t, err) }
func TestPlantDelete404(t *testing.T) { assertPlant404(t, plantTestSvc(t).Delete(999)) }
func TestPlantUpdate404(t *testing.T) {
	_, err := plantTestSvc(t).Update(999, &model.PlantSpecies{Name: "x"})
	assertPlant404(t, err)
}
