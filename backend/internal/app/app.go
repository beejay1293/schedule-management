package app

import (
	"log"
	"net"

	"github.com/beejay1293/schedule-management/backend/internal/repository"
	"github.com/beejay1293/schedule-management/backend/internal/service"

	pb "github.com/beejay1293/schedule-management/backend/internal/pb"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// App holds the gRPC server and listener
type App struct {
	Server   *grpc.Server
	Listener net.Listener
}

// New creates the App, sets up repositories, services, and registers gRPC services
func New(database *gorm.DB) (*App, error) {
	// Create TCP listener
	lis, err := net.Listen("tcp", ":50051") // default port
	if err != nil {
		return nil, err
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// --- Repositories ---
	appointmentRepo := repository.NewAppointmentRepo(database)

	// --- Services ---
	appointmentSvc := service.NewAppointmentService(appointmentRepo)

	// --- Register gRPC services ---
	pb.RegisterAppointmentServiceServer(grpcServer, appointmentSvc)

	log.Println("All gRPC services registered")

	return &App{
		Server:   grpcServer,
		Listener: lis,
	}, nil
}

// Run starts the gRPC server
func (a *App) Run() {
	log.Printf("gRPC server running on %s", a.Listener.Addr())
	if err := a.Server.Serve(a.Listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
