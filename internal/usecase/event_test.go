package usecase

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/binaryty/evbot/internal/config"
	domain "github.com/binaryty/evbot/internal/domain/entities"
	"github.com/binaryty/evbot/internal/repository/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite memory db: %v", err)
	}

	schema := `
	CREATE TABLE events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		date DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		archived INTEGER NOT NULL DEFAULT 0
	);`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func TestArchiveEventFlow(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() { db.Close() })

	repo := sqlite.NewEventRepository(db)
	cfg := &config.Config{AdminIDs: []int64{1}}

	uc := NewEventUseCase(repo, cfg)

	ctx := context.Background()
	event := domain.Event{
		Title:       "Test event",
		Description: "desc",
		Date:        time.Now().Add(time.Hour),
	}

	eventID, err := uc.CreateEvent(ctx, 1, event)
	if err != nil {
		t.Fatalf("create event failed: %v", err)
	}

	events, err := uc.ListEvents(ctx)
	if err != nil {
		t.Fatalf("list events failed: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if err := uc.ArchiveEvent(ctx, 1, eventID); err != nil {
		t.Fatalf("archive event failed: %v", err)
	}

	events, err = uc.ListEvents(ctx)
	if err != nil {
		t.Fatalf("list events failed after archive: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected 0 active events after archive, got %d", len(events))
	}

	archived, err := uc.ListArchivedEvents(ctx)
	if err != nil {
		t.Fatalf("list archived failed: %v", err)
	}
	if len(archived) != 1 || archived[0].ID != eventID {
		t.Fatalf("expected archived event with ID %d, got %+v", eventID, archived)
	}
}
