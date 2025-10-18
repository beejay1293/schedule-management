package repository

import (
	"time"

	"github.com/beejay1293/schedule-management/backend/internal/models"
)

type AppointmentRepository interface {
	Create(a *models.Appointment) error
	List() ([]models.Appointment, error)
	Delete(id string) error
	GetByID(id string) (*models.Appointment, error)
	Search(title, date string) ([]models.Appointment, error)
	ExistsInRange(start, end time.Time) (bool, error)
}
