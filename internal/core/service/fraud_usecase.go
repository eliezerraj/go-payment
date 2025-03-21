package service

import (
	"fmt"
	"errors"
	"context"
	"encoding/json"

	"github.com/go-payment/internal/core/model"
	"github.com/go-payment/internal/core/erro"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/grpc/metadata"
	"github.com/golang/protobuf/jsonpb"
	pb "github.com/golang/protobuf/proto"
	proto "github.com/go-payment/internal/adapter/grpc/proto"
	proto_pod "github.com/go-payment/internal/core/proto/pod"
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

// About get gprc server information pod 
func (s WorkerService) GetInfoPodGrpc(ctx context.Context) (*model.InfoPod, error){
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Msg("GetInfoPodGrpc")

	// Trace
	span := tracerProvider.Span(ctx, "service.GetInfoPodGrpc")
	defer span.End()
	
	// Prepare to receive proto data
	data := &proto.PodInfoRequest {}
	data_pod := &proto_pod.PodInfoRequest {}
	_ = data_pod
	// Prepare the client
	client := s.grpcClient.GetConnection()

	// Set header for authorization
	header := metadata.New(map[string]string{"client-id": "client-001", "authorization": "Beared cookie"})
	ctx = metadata.NewOutgoingContext(ctx, header)

	// request the data from grpc
	response, err := client.GetPodInfo(ctx, data)
	if err != nil {
		childLogger.Error().Err(err).Msg("Error not GetPodInfo")
	  	return nil, err
	}

	// conver proto to json
	response_str, err := ProtoToJSON(response)
	if err != nil {
		return nil, err
  	}

	// convert json to struct
	var result_final map[string]interface{}
	err = json.Unmarshal([]byte(response_str), &result_final)
	if err != nil {
		return nil, err
	}

	result_filtered := result_final["podInfo"].(map[string]interface{})

	var infoPod model.InfoPod
	jsonString, err := json.Marshal(result_filtered)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(jsonString, &infoPod)

	return &infoPod, nil
}

// About check the fraud score from featrures 
func (s WorkerService) CheckFeaturePaymentFraudGrpc(ctx context.Context, paymentFraud *model.PaymentFraud) (*model.PaymentFraud, error){
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Msg("CheckFeaturePaymentFraudGrpc")

	// Trace
	span := tracerProvider.Span(ctx, "service.CheckFeaturePaymentFraudGrpc")
	defer span.End()
	
	// set header for authorization
	header := metadata.New(map[string]string{"client-id": "client-001", "authorization": "Beared cookie"})
	ctx = metadata.NewOutgoingContext(ctx, header)

	// set time
	ts_proto_paymentAt := timestamppb.New(paymentFraud.PaymentAt)

	// proto data
	payment_proto := proto.Payment{	AccountId: paymentFraud.AccountID,
								CardNumber: paymentFraud.CardNumber,
								TerminalName: paymentFraud.TerminalName,
								CoordX:  paymentFraud.CoordX,
								CoordY:  paymentFraud.CoordY,
								CardType: paymentFraud.CardType,
								CardModel: paymentFraud.CardModel,
								Currency: paymentFraud.Currency,
								Mcc: paymentFraud.MCC,
								Amount: paymentFraud.Amount,
								Status: paymentFraud.Status,
								PaymentAt: ts_proto_paymentAt,
								Tx_1D:  paymentFraud.Tx1Day,
								Avg_1D: paymentFraud.Avg1Day,
								Tx_7D: paymentFraud.Tx7Day,
								Avg_7D: paymentFraud.Avg7Day,
								Tx_30D: paymentFraud.Tx30Day,
								Avg_30D: paymentFraud.Avg30Day,
								TimeBtwCcTx: paymentFraud.TimeBtwTx,
	}

	// Prepare
	data := &proto.PaymentRequest {Payment: &payment_proto}
	client := s.grpcClient.GetConnection()

	// request data from grpc server
	response, err := client.CheckPaymentFraud(ctx, data)
	if err != nil {
	  	return nil, err
	}
	
	// convert proto to json
	response_str, err := ProtoToJSON(response)
	if err != nil {
		return nil, err
  	}

	var result_final map[string]interface{}
	err = json.Unmarshal([]byte(response_str), &result_final)
	if err != nil {
		return nil, err
	}
	
	result_filtered := result_final["payment"].(map[string]interface{})
	var parse_paymentFraud model.PaymentFraud

	jsonString, err := json.Marshal(result_filtered)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(jsonString, &parse_paymentFraud)
	
	return &parse_paymentFraud, nil
}

// About add a payment with fraud score 
func (s WorkerService) AddPaymentWithCheckFraud(ctx context.Context, payment *model.Payment) (*model.Payment, error){
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Msg("AddPaymentWithCheckFraud")

	// Trace
	span := tracerProvider.Span(ctx, "service.AddPaymentWithCheckFraud")
	trace_id := fmt.Sprintf("%v",ctx.Value("trace-request-id"))

	// get connection
	tx, conn, err := s.workerRepository.DatabasePGServer.StartTx(ctx)
	if err != nil {
		return nil, err
	}
	
	// handle tx
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
		s.workerRepository.DatabasePGServer.ReleaseTx(conn)
		span.End()
	}()

	//Businness rule
	if (payment.CardType != "CREDIT") && (payment.CardType != "DEBIT") {
		span.RecordError(erro.ErrCardTypeInvalid)
		return nil, erro.ErrCardTypeInvalid
	}

	// Get Card data
	card := model.Card{CardNumber: payment.CardNumber}
	res_card, err := s.workerRepository.GetCard(ctx, &card)
	if err != nil {
		return nil, err
	}

	// Get terminal
	terminal := model.Terminal{Name: payment.TerminalName}
	res_terminal, err := s.workerRepository.GetTerminal(ctx, &terminal)
	if err != nil {
		return nil, err
	}

	// add payment
	payment.FkCardID = res_card.ID
	payment.FkTerminalId = res_terminal.ID
	payment.Status = "PENDING"

	res_payment, err := s.workerRepository.AddPayment(ctx, tx, payment)
	if err != nil {
		return nil, err
	}

	// set the pk
	payment.ID = res_payment.ID

	// get fund balance
	res_payload, statusCode, err := apiService.CallApi(ctx,
														s.apiService[1].Url + "/" + res_card.AccountID,
														s.apiService[1].Method,
														&s.apiService[1].Header_x_apigw_api_id,
														nil,
														&trace_id,
														nil)
	if err != nil {
		return nil, errorStatusCode(statusCode)
	}

	jsonString, err  := json.Marshal(res_payload)
	if err != nil {
		return nil, errors.New(err.Error())
    }
	var movimentAccount model.MovimentAccount
	json.Unmarshal(jsonString, &movimentAccount)

	// apply business rule
	if (movimentAccount.AccountBalance.Amount < payment.Amount) {
		payment.Status = "DECLINED:NO-FUND"
	} else {
		payment.Status = "APPROVED"
	}

	// Get Payment Feature for ML Fraud xgboost Grpc
	payment_fraud := model.PaymentFraud{}
	res_pay_fraud, err := s.workerRepository.GetPaymentFraudFeature(ctx, payment)
	if err != nil {
		switch err {
			case erro.ErrNotFound:
				payment_fraud.CardNumber = payment.CardNumber
				payment_fraud.TerminalName = payment.TerminalName
				payment_fraud.MCC = payment.MCC
				payment_fraud.CoordX = int32(terminal.CoordX)
				payment_fraud.CoordY = int32(terminal.CoordY)
				payment_fraud.CardType = payment.CardType
				payment_fraud.CardModel = payment.CardMode
				payment_fraud.Currency = payment.Currency
				payment_fraud.Amount = payment.Amount
				payment_fraud.Tx1Day = 0
				payment_fraud.Avg1Day = 0
				payment_fraud.Tx7Day = 0
				payment_fraud.Avg7Day = 0
				payment_fraud.Tx30Day = 0
				payment_fraud.Avg30Day = 0
				payment_fraud.TimeBtwTx = 0
			default:
				return nil, err
		}
	}else {
		payment_fraud = *res_pay_fraud
	}
	
	// call the method GetPaymentFraudGrpc (above)
	res_svc_fraud, err := s.CheckFeaturePaymentFraudGrpc(ctx, &payment_fraud)
    if err != nil {
		return nil, errors.New(err.Error())
    }

	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Interface("======>>>>res_svc_fraud :", res_svc_fraud).Msg("")
	// set the fraud score
	payment.Fraud = res_svc_fraud.Fraud

	// update status payment and ml features
	res_update, err := s.workerRepository.UpdatePayment(ctx, tx, payment)
	if err != nil {
		return nil, err
	}
	if res_update == 0 {
		err = erro.ErrUpdate
		return nil, err
	}

	return payment, nil
}