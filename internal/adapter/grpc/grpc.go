package grpc

import (

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	proto "github.com/go-payment/internal/proto"

)

var childLogger = log.With().Str("adapter/grpc", "fraud").Logger()

type GrpcClient struct {
	Client 		proto.FraudServiceClient
	GrcpClient	*grpc.ClientConn
}

func StartGrpcClient(HOST string) (GrpcClient, error){
	childLogger.Debug().Msg("StartGrpcClient")

	var opts []grpc.DialOption
	opts = append(opts, grpc.FailOnNonTempDialError(true)) // Wait for ready
	opts = append(opts, grpc.WithBlock()) // Wait for ready
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)) // 
	opts = append(opts, grpc.WithInsecure()) // no TLS

	conn, err := grpc.Dial(HOST, opts...)
	if err != nil {
	  childLogger.Error().Err(err).Msg("erro connect to grpc server")
	  return GrpcClient{}, err
	}

	client := proto.NewFraudServiceClient(conn)
	childLogger.Info().Interface("Grpc Client running : ", client).Msg("")

	return GrpcClient{
		Client: client,
		GrcpClient : conn,
	}, nil
}

func (s GrpcClient) GetConnection() (proto.FraudServiceClient) {
	childLogger.Debug().Msg("GetConnection")
	return s.Client
}

func (s GrpcClient) CloseConnection() () {
	childLogger.Debug().Msg("CloseConnection")

	if err := s.GrcpClient.Close(); err != nil {
		childLogger.Error().Err(err).Msg("Failed to close gPRC connection")
	}
}