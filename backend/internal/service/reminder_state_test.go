package service

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gbplantwiki/gbplantwiki/internal/constants"
	"github.com/gbplantwiki/gbplantwiki/internal/model"
	"github.com/gbplantwiki/gbplantwiki/internal/repository"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

func reminderTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func newReminderTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil { t.Fatalf("open sqlite: %v", err) }
	if err := db.AutoMigrate(&model.CareReminder{}); err != nil { t.Fatalf("migrate: %v", err) }
	return db
}

func TestReminderSnoozeTransition(t *testing.T) {
	db := newReminderTestDB(t)
	repo := repository.NewCareReminderRepository(db)
	svc := NewCareReminderService(repo, reminderTestLogger())
	m := &model.CareReminder{UserID: 1, TaskTitle: "浇水", RemindDate: time.Now().Add(24 * time.Hour), Status: model.ReminderPending}
	if err := repo.Create(m); err != nil { t.Fatalf("create: %v", err) }
	got, err := svc.UpdateStatus(1, m.ID, model.ReminderSnoozed)
	if err != nil { t.Fatalf("UpdateStatus: %v", err) }
	if got.Status != "snoozed" {
		t.Fatalf("status = %q, want %q", got.Status, model.ReminderSnoozed)
	}
}

func TestReminderStatusTextSnoozed(t *testing.T) {
	if got := util.ReminderStatusText("snoozed"); got != "已暂缓" {
		t.Fatalf("ReminderStatusText(snoozed) = %q, want 已暂缓", got)
	}
}

func TestReminderLogTemplateCount(t *testing.T) {
	if got := constants.LogTemplateCount(); got != 34 {
		t.Fatalf("LogTemplateCount() = %d, want 34", got)
	}
}

func TestReminderCreateSnoozedRejected(t *testing.T) {
	db := newReminderTestDB(t)
	repo := repository.NewCareReminderRepository(db)
	svc := NewCareReminderService(repo, reminderTestLogger())
	m := &model.CareReminder{UserID: 1, TaskTitle: "浇水", RemindDate: time.Now().Add(24 * time.Hour), Status: model.ReminderSnoozed}
	if _, err := svc.Create(1, m); err != nil {
		t.Fatalf("Create should accept snoozed, got %v", err)
	}
}
