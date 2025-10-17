package repository

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	database "github.com/beejay1293/schedule-management/backend/internal/db"
	"github.com/beejay1293/schedule-management/backend/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAppointmentRepo_AllMethods(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewAppointmentRepo(db)

	// --- Test Create ---
	appt := &models.Appointment{
		ID:    "test-1",
		Title: "Test Appointment",
		Date:  time.Now().Add(1 * time.Hour),
	}

	if err := repo.Create(appt); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Conflict scenario: creating another appointment at the same time
	conflictAppt := &models.Appointment{
		ID:    "test-2",
		Title: "Conflict Appointment",
		Date:  appt.Date,
	}
	err := repo.Create(conflictAppt)
	if err == nil || status.Code(err) != codes.AlreadyExists {
		t.Fatalf("expected conflict error, got: %v", err)
	}

	// --- Test ExistsAt ---
	exists, err := repo.ExistsAt(appt.Date)
	if err != nil {
		t.Fatalf("ExistsAt failed: %v", err)
	}
	if !exists {
		t.Fatalf("ExistsAt should return true")
	}

	nonExistDate := time.Now().Add(24 * time.Hour)
	exists, err = repo.ExistsAt(nonExistDate)
	if err != nil {
		t.Fatalf("ExistsAt failed: %v", err)
	}
	if exists {
		t.Fatalf("ExistsAt should return false for future date")
	}

	// --- Test List ---
	list, err := repo.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("List should return 1 appointment, got %d", len(list))
	}

	// --- Test GetByID ---
	got, err := repo.GetByID("test-1")
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got == nil || got.ID != "test-1" {
		t.Fatalf("GetByID returned wrong appointment")
	}

	// Non-existing ID should return nil
	got, _ = repo.GetByID("non-existent")
	if got != nil {
		t.Fatalf("GetByID should return nil for non-existent ID")
	}

	// --- Test Delete ---
	err = repo.Delete("test-1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Deleting non-existent ID should not error
	err = repo.Delete("non-existent")
	if err != nil {
		t.Fatalf("Delete failed for non-existent ID: %v", err)
	}

	// --- Test Search ---
	// Add multiple appointments for search
	repo.Create(&models.Appointment{
		ID:    "test-3",
		Title: "Doctor Visit",
		Date:  time.Now().Add(2 * time.Hour),
	})
	repo.Create(&models.Appointment{
		ID:    "test-4",
		Title: "Dentist",
		Date:  time.Now().Add(3 * time.Hour),
	})

	// Search by title
	results, err := repo.Search("Doctor", "")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 1 || results[0].Title != "Doctor Visit" {
		t.Fatalf("Search by title failed, got %+v", results)
	}

	// Invalid date format
	_, err = repo.Search("", "invalid-date")
	if err == nil || !strings.Contains(err.Error(), "invalid date") {
		t.Fatalf("Search should fail for invalid date, got: %v", err)
	}
}

// --- Test Concurrency ---
// Two users trying to create an appointment at the same time
func TestAppointmentRepo_Create_ConcurrentConflict(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewAppointmentRepo(db)

	date := time.Now().Add(2 * time.Hour)
	numGoroutines := 2
	errs := make(chan error, numGoroutines)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 1; i <= numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			appt := &models.Appointment{
				ID:    fmt.Sprintf("concurrent-%d", i),
				Title: fmt.Sprintf("Concurrent %d", i),
				Date:  date,
			}
			errs <- repo.Create(appt)
		}(i)
	}

	wg.Wait()
	close(errs)

	var successCount, conflictCount int
	for err := range errs {
		if err == nil {
			successCount++
		} else if status.Code(err) == codes.AlreadyExists {
			conflictCount++
		} else {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if successCount != 1 || conflictCount != 1 {
		t.Fatalf("expected 1 success and 1 conflict, got %d success and %d conflict", successCount, conflictCount)
	}
}
