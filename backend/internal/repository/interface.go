package repository

import (
	"time"

	"github.com/beejay1293/schedule-management/backend/internal/models"
)

type AppointmentRepository interface {
	Create(a *models.Appointment) error
	List() ([]models.Appointment, error)
	Delete(id string) error
	ExistsAt(date time.Time) (bool, error)
}
