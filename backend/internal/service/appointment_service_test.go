package service

import (
	"context"
	"testing"
	"time"

	"github.com/beejay1293/schedule-management/backend/internal/pb"
	"github.com/beejay1293/schedule-management/backend/internal/repository"

	database "github.com/beejay1293/schedule-management/backend/internal/db"
)

func setupTestService(t *testing.T) *AppointmentService {
	db := database.SetupTestDB(t)
	repo := repository.NewAppointmentRepo(db)
	return NewAppointmentService(repo)
}

func TestAppointmentService_CreateAppointment(t *testing.T) {
	svc := setupTestService(t)

	req := &pb.CreateAppointmentRequest{
		Title: "Meeting",
		Date:  time.Now().Add(1 * time.Hour).Format("2006-01-02"),
		Time:  time.Now().Add(1 * time.Hour).Format("15:04"),
	}

	resp, err := svc.CreateAppointment(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Appointment.Title != req.Title {
		t.Fatalf("expected title %s, got %s", req.Title, resp.Appointment.Title)
	}
}

func TestAppointmentService_ListAppointments(t *testing.T) {
	svc := setupTestService(t)
	ctx := context.Background()

	// Seed data
	for i := range 2 {
		req := &pb.CreateAppointmentRequest{
			Title: "Checkup",
			Date:  time.Now().Add(time.Duration(i) * time.Hour).Format("2006-01-02"),
			Time:  time.Now().Add(time.Duration(i) * time.Hour).Format("15:04"),
		}
		_, _ = svc.CreateAppointment(ctx, req)
	}

	resp, err := svc.ListAppointments(ctx, &pb.ListAppointmentsRequest{})
	if err != nil {
		t.Fatalf("failed to list appointments: %v", err)
	}
	if len(resp.Appointments) < 2 {
		t.Fatalf("expected at least 2 appointments, got %d", len(resp.Appointments))
	}
}

func TestAppointmentService_SearchAppointments(t *testing.T) {
	svc := setupTestService(t)
	ctx := context.Background()

	// Seed data
	req := &pb.CreateAppointmentRequest{
		Title: "Doctor Visit",
		Date:  "2025-10-16",
		Time:  "14:00",
	}
	_, _ = svc.CreateAppointment(ctx, req)

	searchReq := &pb.SearchAppointmentsRequest{Title: "Doctor"}
	resp, err := svc.SearchAppointments(ctx, searchReq)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(resp.Appointments) != 1 || resp.Appointments[0].Title != "Doctor Visit" {
		t.Fatalf("unexpected search result: %+v", resp.Appointments)
	}
}

func TestAppointmentService_DeleteAppointment(t *testing.T) {
	svc := setupTestService(t)
	ctx := context.Background()

	// Create appointment
	req := &pb.CreateAppointmentRequest{
		Title: "To Delete",
		Date:  time.Now().Format("2006-01-02"),
		Time:  time.Now().Add(2 * time.Hour).Format("15:04"),
	}
	createResp, err := svc.CreateAppointment(ctx, req)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Delete appointment
	_, err = svc.DeleteAppointment(ctx, &pb.DeleteAppointmentRequest{Id: createResp.Appointment.Id})
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Verify deletion
	listResp, _ := svc.ListAppointments(ctx, &pb.ListAppointmentsRequest{})
	for _, a := range listResp.Appointments {
		if a.Id == createResp.Appointment.Id {
			t.Fatalf("appointment not deleted")
		}
	}
}
