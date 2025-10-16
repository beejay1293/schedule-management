package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/beejay1293/schedule-management/backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type appointmentRepo struct {
	db *gorm.DB
}

func NewAppointmentRepo(db *gorm.DB) AppointmentRepository {
	return &appointmentRepo{db: db}
}

func (r *appointmentRepo) Create(a *models.Appointment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Lock rows that could conflict (same date/time)
		var existing models.Appointment
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("date = ?", a.Date).
			Take(&existing).Error

		if err == nil {
			return fmt.Errorf("conflict: appointment already exists at this date/time")
		}

		// If no record found, create a new one
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(a).Error; err != nil {
				return err
			}
			return nil
		}

		return err
	})
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
