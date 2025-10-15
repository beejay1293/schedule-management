package main

import (
	"context"
	"log"
	"time"

	pb "github.com/beejay1293/schedule-management/backend/internal/pb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAppointmentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Create appointment
	res, err := client.CreateAppointment(ctx, &pb.CreateAppointmentRequest{
		Title: "Doctor Appointment",
		Date:  "2025-10-14",
		Time:  "14:00",
	})
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Created: %+v", res.Appointment)

	// List appointments
	list, _ := client.ListAppointments(ctx, &pb.ListAppointmentsRequest{})
	log.Printf("All appointments: %+v", list.Appointments)
}
