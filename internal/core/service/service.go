package service

import(
	"github.com/go-payment/internal/adapter/grpc/client"
	"github.com/go-payment/internal/core/model"
	"github.com/go-payment/internal/adapter/database"
	"github.com/rs/zerolog/log"
)

var childLogger = log.With().Str("component","go-payment").Str("package","internal.core.service").Logger()

type WorkerService struct {
	workerRepository *database.WorkerRepository
	apiService		[]model.ApiService
	grpcClient		*client.GrpcClient
}

func NewWorkerService(	workerRepository *database.WorkerRepository,
						apiService		[]model.ApiService,
						grpcClient		*client.GrpcClient) *WorkerService{
	childLogger.Info().Str("func","NewWorkerService").Send()

	return &WorkerService{
		workerRepository: workerRepository,
		apiService: apiService,
		grpcClient: grpcClient,
	}
}