package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/beejay1293/schedule-management/backend/internal/models"
	"gorm.io/gorm"
)

type appointmentRepo struct {
	db *gorm.DB
}

func NewAppointmentRepo(db *gorm.DB) AppointmentRepository {
	return &appointmentRepo{db: db}
}

func (r *appointmentRepo) Create(a *models.Appointment) error {
	// Try inserting the new appointment
	if err := r.db.Create(a).Error; err != nil {
		// Check if the error is a unique constraint violation
		if strings.Contains(err.Error(), "unique_appointment_time") {
			return fmt.Errorf("conflict: appointment already exists at this date/time")
		}
		return err
	}
	return nil
}

func (r *appointmentRepo) List() ([]models.Appointment, error) {
	var appts []models.Appointment
	err := r.db.Order("date ASC").Find(&appts).Error
	return appts, err
}

func (r *appointmentRepo) Delete(id string) error {
	return r.db.Delete(&models.Appointment{}, "id = ?", id).Error
}

func (r *appointmentRepo) ExistsAt(date time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&models.Appointment{}).Where("date = ?", date).Count(&count).Error
	return count > 0, err
}

func (r *appointmentRepo) Search(title, date string) ([]models.Appointment, error) {
	var appts []models.Appointment
	tx := r.db.Model(&models.Appointment{})

	if title != "" {
		tx = tx.Where("title ILIKE ?", "%"+title+"%") // case-insensitive search
	}

	if date != "" {
		// parse date string to match the date column
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format")
		}
		tx = tx.Where("date::date = ?", parsedDate.Format("2006-01-02"))
	}

	err := tx.Order("date ASC").Find(&appts).Error
	return appts, err
}

func (r *appointmentRepo) GetByID(id string) (*models.Appointment, error) {
	var appt models.Appointment
	err := r.db.First(&appt, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &appt, nil
}
