package service

import (
	"context"
	"fmt"
	"time"

	"github.com/beejay1293/schedule-management/backend/internal/models"
	"github.com/beejay1293/schedule-management/backend/internal/repository"

	pb "github.com/beejay1293/schedule-management/backend/internal/pb"
	"github.com/google/uuid"
)

type AppointmentService struct {
	repo repository.AppointmentRepository
	pb.UnimplementedAppointmentServiceServer
}

func NewAppointmentService(repo repository.AppointmentRepository) *AppointmentService {
	return &AppointmentService{repo: repo}
}

func (s *AppointmentService) CreateAppointment(ctx context.Context, req *pb.CreateAppointmentRequest) (*pb.CreateAppointmentResponse, error) {
	date, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", req.Date, req.Time))
	if err != nil {
		return nil, fmt.Errorf("invalid date/time format")
	}

	conflict, err := s.repo.ExistsAt(date)
	if err != nil {
		return nil, err
	}

	if conflict {
		return nil, fmt.Errorf("conflict: another appointment exists at that time")
	}

	a := &models.Appointment{
		ID:    uuid.New().String(),
		Title: req.Title,
		Date:  date,
	}

	if err := s.repo.Create(a); err != nil {
		return nil, err
	}

	return &pb.CreateAppointmentResponse{
		Appointment: &pb.Appointment{
			Id:    a.ID,
			Title: a.Title,
			Date:  a.Date.Format("2006-01-02"),
			Time:  a.Date.Format("15:04"),
		},
	}, nil
}

func (s *AppointmentService) ListAppointments(ctx context.Context, req *pb.ListAppointmentsRequest) (*pb.ListAppointmentsResponse, error) {
	appointments, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	var pbAppointments []*pb.Appointment
	for _, a := range appointments {
		pbAppointments = append(pbAppointments, &pb.Appointment{
			Id:    a.ID,
			Title: a.Title,
			Date:  a.Date.Format("2006-01-02"),
			Time:  a.Date.Format("15:04"),
		})
	}

	return &pb.ListAppointmentsResponse{
		Appointments: pbAppointments,
	}, nil
}

func (s *AppointmentService) DeleteAppointment(ctx context.Context, req *pb.DeleteAppointmentRequest) (*pb.DeleteAppointmentResponse, error) {
	err := s.repo.Delete(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteAppointmentResponse{}, nil
}
