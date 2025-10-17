package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/beejay1293/schedule-management/backend/internal/models"
	"github.com/beejay1293/schedule-management/backend/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/beejay1293/schedule-management/backend/internal/pb"
)

type AppointmentService struct {
	repo repository.AppointmentRepository
	pb.UnimplementedAppointmentServiceServer

	// subscribers for streaming events
	mu          sync.Mutex
	subscribers map[int64]chan *pb.AppointmentEvent
	nextID      int64
}

func NewAppointmentService(repo repository.AppointmentRepository) *AppointmentService {
	return &AppointmentService{
		repo:        repo,
		subscribers: make(map[int64]chan *pb.AppointmentEvent),
	}
}

// Create Appointment
// No in-memory mutex here because database-level uniqueness ensures concurrent safety across multiple service instances
func (s *AppointmentService) CreateAppointment(ctx context.Context, req *pb.CreateAppointmentRequest) (*pb.CreateAppointmentResponse, error) {
	date, err := time.Parse("2006-01-02 15:04", req.Date+" "+req.Time)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date/time format")
	}

	conflict, err := s.repo.ExistsAt(date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check existing appointments: %v", err)
	}
	if conflict {
		return nil, status.Error(codes.AlreadyExists, "conflict: another appointment exists at that time")
	}

	a := &models.Appointment{
		ID:    uuid.New().String(),
		Title: req.Title,
		Date:  date,
	}

	if err := s.repo.Create(a); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create appointment: %v", err)
	}

	pbAppt := &pb.Appointment{
		Id:    a.ID,
		Title: a.Title,
		Date:  a.Date.UTC().Format("2006-01-02"),
		Time:  a.Date.UTC().Format("15:04"),
	}

	s.broadcastEvent(&pb.AppointmentEvent{
		Type:        pb.EventType_CREATED,
		Appointment: pbAppt,
	})

	return &pb.CreateAppointmentResponse{Appointment: pbAppt}, nil
}

// Delete Appointment
func (s *AppointmentService) DeleteAppointment(ctx context.Context, req *pb.DeleteAppointmentRequest) (*pb.DeleteAppointmentResponse, error) {
	a, err := s.repo.GetByID(req.Id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Error(codes.NotFound, "appointment not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve appointment: %v", err)
	}

	if err := s.repo.Delete(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete appointment: %v", err)
	}

	pbAppt := &pb.Appointment{
		Id:    a.ID,
		Title: a.Title,
		Date:  a.Date.UTC().Format("2006-01-02"),
		Time:  a.Date.UTC().Format("15:04"),
	}

	s.broadcastEvent(&pb.AppointmentEvent{
		Type:        pb.EventType_DELETED,
		Appointment: pbAppt,
	})

	return &pb.DeleteAppointmentResponse{}, nil
}

// List Appointments
func (s *AppointmentService) ListAppointments(ctx context.Context, req *pb.ListAppointmentsRequest) (*pb.ListAppointmentsResponse, error) {
	appointments, err := s.repo.List()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list appointments: %v", err)
	}

	var pbAppointments []*pb.Appointment
	for _, a := range appointments {
		utcDate := a.Date.UTC()
		pbAppointments = append(pbAppointments, &pb.Appointment{
			Id:    a.ID,
			Title: a.Title,
			Date:  utcDate.Format("2006-01-02"),
			Time:  utcDate.Format("15:04"),
		})
	}

	return &pb.ListAppointmentsResponse{Appointments: pbAppointments}, nil
}

// Search Appointments
func (s *AppointmentService) SearchAppointments(ctx context.Context, req *pb.SearchAppointmentsRequest) (*pb.ListAppointmentsResponse, error) {
	appointments, err := s.repo.Search(req.Title, req.Date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search appointments: %v", err)
	}

	var pbAppointments []*pb.Appointment
	for _, a := range appointments {
		utcDate := a.Date.UTC()
		pbAppointments = append(pbAppointments, &pb.Appointment{
			Id:    a.ID,
			Title: a.Title,
			Date:  utcDate.Format("2006-01-02"),
			Time:  utcDate.Format("15:04"),
		})
	}

	return &pb.ListAppointmentsResponse{Appointments: pbAppointments}, nil
}

// Stream Appointments
func (s *AppointmentService) StreamAppointments(req *pb.StreamAppointmentsRequest, stream pb.AppointmentService_StreamAppointmentsServer) error {
	ch := make(chan *pb.AppointmentEvent, 10)

	s.mu.Lock()
	id := s.nextID
	s.nextID++
	s.subscribers[id] = ch
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.subscribers, id)
		close(ch)
		s.mu.Unlock()
	}()

	for event := range ch {
		if err := stream.Send(event); err != nil {
			log.Printf("stream send error: %v", err)
			if status.Code(err) == codes.Canceled {
				log.Printf("Client closed the stream (context canceled)")
				break
			}
			log.Printf("Error sending event: %v", err)
		}
	}

	return nil
}

// broadcastEvent broadcasts to all subscribers
func (s *AppointmentService) broadcastEvent(event *pb.AppointmentEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ch := range s.subscribers {
		select {
		case ch <- event:
		default:
		}
	}
}
