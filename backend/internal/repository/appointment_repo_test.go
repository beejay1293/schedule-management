package repository

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	database "github.com/beejay1293/schedule-management/backend/internal/db"
	"github.com/beejay1293/schedule-management/backend/internal/models"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAppointmentRepo_AllMethods(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewAppointmentRepo(db)

	now := time.Now()

	test1 := uuid.New()
	test2 := uuid.New()

	// --- Test Create ---
	appt := &models.Appointment{
		ID:        test1,
		Title:     "Test Appointment",
		StartTime: now,
		EndTime:   now.Add(30 * time.Minute),
	}

	if err := repo.Create(appt); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// --- Test conflict (overlapping time range) ---
	conflictAppt := &models.Appointment{
		ID:        test2,
		Title:     "Conflict Appointment",
		StartTime: now, // overlaps with test-1
		EndTime:   now.Add(30 * time.Minute),
	}

	err := repo.Create(conflictAppt)
	if err == nil || status.Code(err) != codes.AlreadyExists {
		t.Fatalf("expected conflict error, got: %v", err)
	}

	// --- Test ExistsInRange ---
	exists, err := repo.ExistsInRange(now.Add(20*time.Minute), now.Add(50*time.Minute))
	if err != nil {
		t.Fatalf("ExistsInRange failed: %v", err)
	}
	if !exists {
		t.Fatalf("ExistsInRange should return true for overlapping range")
	}

	// Non-overlapping range
	nonOverlapStart := now.Add(4 * time.Hour)
	nonOverlapEnd := now.Add(5 * time.Hour)
	exists, err = repo.ExistsInRange(nonOverlapStart, nonOverlapEnd)
	if err != nil {
		t.Fatalf("ExistsInRange failed: %v", err)
	}
	if exists {
		t.Fatalf("ExistsInRange should return false for non-overlapping range")
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
	got, err := repo.GetByID(test1.String())
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got == nil || got.ID != test1 {
		t.Fatalf("GetByID returned wrong appointment")
	}

	// Non-existing ID should return nil
	got, _ = repo.GetByID("non-existent")
	if got != nil {
		t.Fatalf("GetByID should return nil for non-existent ID")
	}

	// --- Test Delete ---
	err = repo.Delete(test1.String())
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Deleting non-existent ID should not error
	err = repo.Delete(uuid.NewString())
	if err != nil {
		t.Fatalf("Delete failed for non-existent ID: %v", err)
	}

	// --- Test Search ---
	// Add multiple appointments for search
	repo.Create(&models.Appointment{
		ID:        uuid.New(),
		Title:     "Doctor Visit",
		StartTime: now.Add(2 * time.Hour),
		EndTime:   now.Add(3 * time.Hour),
	})
	repo.Create(&models.Appointment{
		ID:        uuid.New(),
		Title:     "Dentist",
		StartTime: now.Add(3 * time.Hour),
		EndTime:   now.Add(4 * time.Hour),
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
func TestAppointmentRepo_Created_ConcurrentConflict(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewAppointmentRepo(db)

	start := time.Now().Add(8 * time.Hour)
	end := start.Add(30 * time.Minute)

	numGoroutines := 2
	errs := make(chan error, numGoroutines)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 1; i <= numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			appt := &models.Appointment{
				ID:        uuid.New(),
				Title:     fmt.Sprintf("Concurrent %d", i),
				StartTime: start,
				EndTime:   end,
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
			t.Fatalf("unexpected error: %v, success_count, %v", err, successCount)
		}
	}

	if successCount != 1 || conflictCount != 1 {
		t.Fatalf("expected 1 success and 1 conflict, got %d success and %d conflict", successCount, conflictCount)
	}
}
