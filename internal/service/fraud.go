package service

import (
	"context"
	"encoding/json"

	"github.com/go-payment/internal/core"
	"google.golang.org/grpc/metadata"

	pb "github.com/golang/protobuf/proto"
	proto "github.com/go-payment/internal/proto"
	"github.com/golang/protobuf/jsonpb"

	"go.opentelemetry.io/otel"
)

func ProtoToJSON(msg pb.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}

	return marshaler.MarshalToString(msg)
}

func JSONToProto(data string, msg pb.Message) error {
	return jsonpb.UnmarshalString(data, msg)
}

func (s WorkerService) GetInfoPodGrpc(ctx context.Context) (interface{}, error){
	childLogger.Debug().Msg("GetInfoPodGrpc")

	ctx, span := otel.Tracer("appName").Start(ctx,"svc.GetInfoPodGrpc")
	defer span.End()

	header := metadata.New(map[string]string{"client-id": "client-001", "authorization": "Beared cookie"})
	ctx = metadata.NewOutgoingContext(ctx, header)

	data := &proto.PodInfoRequest {}
	client := s.grpcClient.GetConnection()

	response, err := client.GetPodInfo(ctx, data)
	if err != nil {
		childLogger.Error().Err(err).Msg("Error not GetPodInfo")
	  	return nil, err
	}
	response_str, err := ProtoToJSON(response)
	if err != nil {
		childLogger.Error().Err(err).Msg("Error not ProtoToJSON")
		return nil, err
  	}

	var result_final map[string]interface{}
	err = json.Unmarshal([]byte(response_str), &result_final)
	if err != nil {
		childLogger.Error().Err(err).Msg("Error Unmarshal")
		return nil, err
	}

	result_filtered := result_final["podInfo"].(map[string]interface{})
	var podInfo core.InfoPod

	childLogger.Debug().Interface("result_filtered :", result_filtered).Msg("")

	jsonString, err := json.Marshal(result_filtered)
	if err != nil {
		childLogger.Error().Err(err).Msg("Error Marshal")
		return nil, err
	}
	json.Unmarshal(jsonString, &podInfo)

	childLogger.Debug().Interface("podInfo :", podInfo).Msg("")

	return &podInfo, nil
}

func (s WorkerService) CheckPaymentFraud(ctx context.Context) (interface{}, error){
	childLogger.Debug().Msg("CheckPaymentFraud")


	return true, nil
}