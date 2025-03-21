package client

import (
	"time"
	"context"
	"github.com/rs/zerolog/log"

	"github.com/go-payment/internal/core/erro"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	proto "github.com/go-payment/internal/adapter/grpc/proto"
)

var childLogger = log.With().Str("adapter.grpc", "client").Logger()

type GrpcClient struct {
	ServiceClient 	proto.FraudServiceClient
	GrcpClient		*grpc.ClientConn
}

// About start a grpc client
func StartGrpcClient(host string) (GrpcClient, error){
	childLogger.Info().Msg("StartGrpcClient")

	// Prepare options
	var opts []grpc.DialOption
	opts = append(opts, grpc.FailOnNonTempDialError(true)) // Wait for ready
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)) // 
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTimeout(5*time.Second))
	opts = append(opts, grpc.WithBlock()) // Wait for ready
	
	// Dail a server
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
	  childLogger.Error().Err(err).Msg("erro connect to grpc server")
	  return GrpcClient{}, err
	}

	// Create a client
	serviceClient := proto.NewFraudServiceClient(conn)
	childLogger.Info().Interface("Grpc Client running : ", serviceClient).Msg("")

	return GrpcClient{
		ServiceClient: serviceClient,
		GrcpClient : conn,
	}, nil
}

// About get connection
func (s GrpcClient) TestConnection(ctx context.Context) (error) {
	childLogger.Info().Msg("TestConnection")
	
	if (s.GrcpClient == nil){
		return erro.ErrGRPCConnection
	}
	client := grpc_health_v1.NewHealthClient(s.GrcpClient)
	res, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: ""})
	if err != nil {
		childLogger.Error().Err(err).Msg("health check failed:")
		return err
	}

	childLogger.Info().Interface("Grpc Client running : ", res).Msg("")

	return nil
}

// About get connection
func (s GrpcClient) GetConnection() (proto.FraudServiceClient) {
	childLogger.Info().Msg("GetConnection")
	return s.ServiceClient
}

// About close connection
func (s GrpcClient) CloseConnection() () {
	childLogger.Info().Msg("CloseConnection")

	if err := s.GrcpClient.Close(); err != nil {
		childLogger.Error().Err(err).Msg("Failed to close gPRC connection")
	}
}