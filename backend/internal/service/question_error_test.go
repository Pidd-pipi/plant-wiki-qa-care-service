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

func newQuestionTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.Question{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func questionTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func questionTestSvc(t *testing.T) *QuestionService {
	db := newQuestionTestDB(t)
	return NewQuestionService(repository.NewQuestionRepository(db), nil, nil, questionTestLogger())
}

func TestQuestionGet404(t *testing.T) {
	_, err := questionTestSvc(t).Get(999)
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404 AppError, got %v", err)
	}
}

func TestQuestionCreateError(t *testing.T) {
	db := newQuestionTestDB(t)
	svc := NewQuestionService(repository.NewQuestionRepository(db), nil, nil, questionTestLogger())
	sqlDB, _ := db.DB()
	sqlDB.Close()
	_, err := svc.Create(1, &model.Question{Title: "t", Content: "c"})
	if err == nil {
		t.Fatal("expected error when create fails")
	}
}

func TestQuestionListError(t *testing.T) {
	db := newQuestionTestDB(t)
	svc := NewQuestionService(repository.NewQuestionRepository(db), nil, nil, questionTestLogger())
	sqlDB, _ := db.DB()
	sqlDB.Close()
	_, _, err := svc.List(1, 10)
	if err == nil {
		t.Fatal("expected error when list fails")
	}
}
