package grpc

import (
	"encoding/base64"
	"crypto/x509"
	"crypto/tls"
	"fmt"

	"github.com/go-payment/internal/core"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc"

	proto "github.com/go-payment/internal/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

var childLogger = log.With().Str("adapter/grpc", "fraud").Logger()

type GrpcClient struct {
	Client 		proto.FraudServiceClient
	GrcpClient	*grpc.ClientConn
}

func loadClientCertsTLS(cert *core.Cert) (credentials.TransportCredentials, error) {
	childLogger.Debug().Msg("loadClientCertsTLS")

	//log.Debug().Interface("cert.CaFraudPEM :",cert.CaFraudPEM).Msg("")

	var clientTLSConf *tls.Config

	caPEM_Raw, err := base64.StdEncoding.DecodeString(string(cert.CaFraudPEM))
	if err != nil {
		childLogger.Error().Err(err).Msg("Erro caPEM_Raw !!!")
		return nil, err
	}

	childLogger.Info().Msg("------------------------------------------------")
	fmt.Println(string(caPEM_Raw))
	childLogger.Info().Msg("------------------------------------------------")

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caPEM_Raw); !ok {
		childLogger.Error().Err(err).Msg("Erro AppendCertsFromPEM !!!")
		return nil, err
	}

	clientTLSConf = &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(clientTLSConf), nil
}

func StartGrpcClient(host string,
					cert *core.Cert) (GrpcClient, error){
	childLogger.Debug().Msg("StartGrpcClient")

	var opts []grpc.DialOption
	opts = append(opts, grpc.FailOnNonTempDialError(true)) // Wait for ready
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)) // 
	opts = append(opts, grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor())) // Otel

	// -------------- Load Certs -------------------------	
	if string(cert.CaFraudPEM) != "" {
		tlsCredentials, err := loadClientCertsTLS(cert)
		if err != nil {
			childLogger.Error().Err(err).Msg("Erro loadClientCertsTLS")
			return GrpcClient{}, err
		}
		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials)) // with TLS
	}else {
		opts = append(opts, grpc.WithInsecure()) // no TLS
	}
	// -------------- Load Certs -------------------------

	opts = append(opts, grpc.WithBlock()) // Wait for ready
	
	conn, err := grpc.Dial(host, opts...)
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