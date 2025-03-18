package client

import (
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"

	proto "github.com/go-payment/internal/adapter/grpc/proto"
)

var childLogger = log.With().Str("adapter.grpc", "client").Logger()

type GrpcClient struct {
	ServiceClient 	proto.FraudServiceClient
	GrcpClient		*grpc.ClientConn
}

// About start a grpc client
func StartGrpcClient(host string ) (GrpcClient, error){
	childLogger.Debug().Msg("StartGrpcClient")

	// Prepare
	var opts []grpc.DialOption
	opts = append(opts, grpc.FailOnNonTempDialError(true)) // Wait for ready
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)) // 

	opts = append(opts, grpc.WithInsecure())
	//opts = append(opts, grpc.WithBlock()) // Wait for ready
	
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
func (s GrpcClient) GetConnection() (proto.FraudServiceClient) {
	childLogger.Debug().Msg("GetConnection")
	return s.ServiceClient
}

// About close connection
func (s GrpcClient) CloseConnection() () {
	childLogger.Debug().Msg("CloseConnection")

	if err := s.GrcpClient.Close(); err != nil {
		childLogger.Error().Err(err).Msg("Failed to close gPRC connection")
	}
}