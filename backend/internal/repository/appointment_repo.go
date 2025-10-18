package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/beejay1293/schedule-management/backend/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type appointmentRepo struct {
	db *gorm.DB
}

func NewAppointmentRepo(db *gorm.DB) AppointmentRepository {
	return &appointmentRepo{db: db}
}

// Create inserts a new appointment with proper concurrency handling
func (r *appointmentRepo) Create(a *models.Appointment) error {
	const maxRetries = 2

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		txErr := r.db.Transaction(func(tx *gorm.DB) error {
			return tx.Create(a).Error
		})

		if txErr == nil {
			return nil
		}

		errMsg := txErr.Error()

		// Conflict / constraint violations
		if strings.Contains(errMsg, "duplicate key value violates unique constraint") ||
			strings.Contains(errMsg, "no_overlapping_appointments") {
			return status.Error(codes.AlreadyExists, "conflict: appointment already exists at this date/time")
		}

		// Retry on deadlock only
		if strings.Contains(errMsg, "deadlock detected") {
			if attempt < maxRetries {
				time.Sleep(50 * time.Millisecond) // small backoff
				continue
			}
			return status.Errorf(codes.Aborted, "retry failed after %d attempts: %v", maxRetries, txErr)
		}

		// Any other error: fail
		lastErr = txErr
		break
	}

	return status.Errorf(codes.Internal, "failed to create appointment: %v", lastErr)
}

// List returns all appointments
func (r *appointmentRepo) List() ([]models.Appointment, error) {
	var appts []models.Appointment
	if err := r.db.Order("start_time ASC").Find(&appts).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list appointments: %v", err)
	}
	return appts, nil
}

// Delete removes an appointment by ID
func (r *appointmentRepo) Delete(id string) error {
	if err := r.db.Delete(&models.Appointment{}, "id = ?", id).Error; err != nil {
		return status.Errorf(codes.Internal, "failed to delete appointment: %v", err)
	}
	return nil
}

// ExistsAt checks if an appointment exists at a given date/time
func (r *appointmentRepo) ExistsInRange(start, end time.Time) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Appointment{}).
		Where("start_time < ? AND end_time > ?", end, start).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Search appointments by title and/or date
// Search appointments by title and/or date (based on start_time)
func (r *appointmentRepo) Search(title, date string) ([]models.Appointment, error) {
	var appts []models.Appointment
	tx := r.db.Model(&models.Appointment{})

	dialect := r.db.Dialector.Name()

	// Filter by title (case-insensitive where supported)
	if title != "" {
		if dialect == "sqlite" {
			tx = tx.Where("title LIKE ?", "%"+title+"%")
		} else { // postgres
			tx = tx.Where("title ILIKE ?", "%"+title+"%")
		}
	}

	// Filter by date (matching the day portion of start_time)
	if date != "" {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid date format")
		}

		startOfDay := parsedDate.Format("2006-01-02 00:00:00")
		endOfDay := parsedDate.Add(24 * time.Hour).Format("2006-01-02 00:00:00")

		if dialect == "sqlite" {
			tx = tx.Where("start_time >= ? AND start_time < ?", startOfDay, endOfDay)
		} else { // postgres
			tx = tx.Where("start_time >= ? AND start_time < ?", startOfDay, endOfDay)
		}
	}

	if err := tx.Order("start_time ASC").Find(&appts).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search appointments: %v", err)
	}

	return appts, nil
}

// GetByID retrieves an appointment by ID
func (r *appointmentRepo) GetByID(id string) (*models.Appointment, error) {
	var appt models.Appointment
	err := r.db.First(&appt, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "appointment not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve appointment: %v", err)
	}
	return &appt, nil
}
