package presence

import (
	pb "parkir-pintar/services/presence/gen/presence/v1"
	"parkir-pintar/services/presence/internal/presence/handler"
	"parkir-pintar/services/presence/internal/presence/repository"
	"parkir-pintar/services/presence/internal/presence/usecase"
	"parkir-pintar/services/presence/pkg/billingclient"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Service struct {
	uc usecase.Presence
}

func New(db *pgxpool.Pool, bc billingclient.BillingService) *Service {
	repo := repository.NewPresence(db)
	uc := usecase.NewPresence(repo, bc)
	return &Service{uc: uc}
}

func (s *Service) RegisterGRPC(grpcServer *grpc.Server) {
	srv := handler.NewPresenceServer(s.uc)
	pb.RegisterPresenceServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)
}
