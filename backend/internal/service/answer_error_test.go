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

func newAnswerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil { t.Fatalf("open sqlite: %v", err) }
	if err := db.AutoMigrate(&model.Question{}, &model.Answer{}); err != nil { t.Fatalf("migrate: %v", err) }
	return db
}

func answerTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func answerTestSvc(t *testing.T) *AnswerService {
	db := newAnswerTestDB(t)
	return NewAnswerService(db, repository.NewAnswerRepository(db), repository.NewQuestionRepository(db), answerTestLogger())
}

func TestAnswerAdoptMissingAnswer404(t *testing.T) {
	db := newAnswerTestDB(t)
	qRepo := repository.NewQuestionRepository(db)
	q := &model.Question{UserID: 1, Title: "t", Content: "c"}
	if err := qRepo.Create(q); err != nil { t.Fatalf("seed question: %v", err) }
	svc := NewAnswerService(db, repository.NewAnswerRepository(db), qRepo, answerTestLogger())
	_, err := svc.Adopt(1, q.ID, 999)
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404 AppError, got %v", err)
	}
}

func TestAnswerAdoptMissingQuestion404(t *testing.T) {
	svc := answerTestSvc(t)
	_, err := svc.Adopt(1, 999, 1)
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404 AppError, got %v", err)
	}
}

func TestAnswerLikeMissing404(t *testing.T) {
	svc := answerTestSvc(t)
	_, err := svc.Like(999)
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404 AppError, got %v", err)
	}
}

func TestAnswerCreateMissingQuestion404(t *testing.T) {
	svc := answerTestSvc(t)
	_, err := svc.Create(1, 999, "x")
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404 AppError, got %v", err)
	}
}
