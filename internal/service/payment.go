package service

import (
	"context"
	"strconv"
	"errors"
	"github.com/rs/zerolog/log"
	"encoding/json"
	"github.com/go-payment/internal/lib"

	"github.com/mitchellh/mapstructure"
	"github.com/go-payment/internal/core"
	"github.com/go-payment/internal/erro"
	"github.com/go-payment/internal/adapter/restapi"
	"github.com/go-payment/internal/repository/storage"
	"github.com/go-payment/internal/adapter/grpc"
)

var childLogger = log.With().Str("service", "service").Logger()

type WorkerService struct {
	workerRepo		 	*storage.WorkerRepository
	appServer			*core.AppServer
	restApiService		*restapi.RestApiService
	grpcClient 			*grpc.GrpcClient
}

func NewWorkerService(workerRepo	*storage.WorkerRepository,
						appServer	*core.AppServer,
						restApiService	*restapi.RestApiService,
						grpcClient 	*grpc.GrpcClient) *WorkerService{
	childLogger.Debug().Msg("NewWorkerService")

	return &WorkerService{
		workerRepo:		 	workerRepo,
		appServer:			appServer,
		restApiService:		restApiService,
		grpcClient: 		grpcClient,
	}
}

func (s WorkerService) Auth(ctx context.Context, authUser core.AuthUser) (*core.AuthUser, error){
	childLogger.Debug().Msg("Auth")

	span := lib.Span(ctx, "service.auth")	
	defer span.End()

	path := s.appServer.RestEndpoint.AuthUrlDomain + "/login"
	res_interface, err := s.restApiService.CallRestApi(ctx,	"POST",	path, nil ,authUser)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	var auth_user_parsed core.AuthUser
	err = mapstructure.Decode(res_interface, &auth_user_parsed)
	if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		span.RecordError(err)
		return nil, errors.New(err.Error())
	}

	return &auth_user_parsed, nil
}

func (s WorkerService) Get(ctx context.Context, payment *core.Payment) (*core.Payment, error){
	childLogger.Debug().Msg("Get")
	
	span := lib.Span(ctx, "service.get")	
    defer span.End()

	res, err := s.workerRepo.Get(ctx, payment)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return res, nil
}

func (s WorkerService) Pay(ctx context.Context, payment *core.Payment) (*core.Payment, error){
	childLogger.Debug().Msg("Pay")
	
	span := lib.Span(ctx, "service.pay")	

	tx, conn, err := s.workerRepo.StartTx(ctx)
	if err != nil {
		return nil, err
	}
	
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
		s.workerRepo.ReleaseTx(conn)
		span.End()
	}()

	if (payment.CardType != "CREDIT") && (payment.CardType != "DEBIT") {
		span.RecordError(erro.ErrCardTypeInvalid)
		return nil, erro.ErrCardTypeInvalid
	}
	// Read Card
	card := core.Card{}
	card.CardNumber = payment.CardNumber
	res_interface_card, err := s.workerRepo.GetCard(ctx, &card)
	if err != nil {
		childLogger.Error().Err(err).Msg("error workerRepository.GetCard")
		span.RecordError(err)
		return nil, err
	}

	var card_parsed core.Card
	err = mapstructure.Decode(res_interface_card, &card_parsed)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		span.RecordError(err)
		return nil, errors.New(err.Error())
    }

	// Read Terminal
	terminal := core.Terminal{}
	terminal.Name = payment.TerminalName
	res_interface_term, err := s.workerRepo.GetTerminal(ctx, &terminal)
	if err != nil {
		childLogger.Error().Err(err).Msg("error workerRepository.GetTerminal")
		span.RecordError(err)
		return nil, err
	}
	var terminal_parsed core.Terminal
	err = mapstructure.Decode(res_interface_term, &terminal_parsed)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		span.RecordError(err)
		return nil, errors.New(err.Error())
    }

	// Get Account for Just for Check
	path := s.appServer.RestEndpoint.ServiceUrlDomain + "/getId/" + strconv.Itoa(card_parsed.FkAccountID)
	res_interface_acc, err := s.restApiService.CallRestApi(ctx,	"GET", path, &s.appServer.RestEndpoint.XApigwId, nil)

	if err != nil {
		childLogger.Error().Err(err).Msg("error restApiService.GetData")
		span.RecordError(err)
		return nil, err
	}
	var account_parsed core.Account
	jsonString, err := json.Marshal(res_interface_acc)
	if err != nil {
		childLogger.Error().Err(err).Msg("Error Marshal")
		span.RecordError(err)
		return nil, err
	}

	json.Unmarshal(jsonString, &account_parsed)
	span.AddEvent("Begin Transaction - lock")
	
	payment.FkCardID = card_parsed.ID
	payment.FkTerminalId = terminal_parsed.ID
	payment.Status = "PENDING"
	res, err := s.workerRepo.Add(ctx, tx ,payment)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Get Fund
	path = s.appServer.RestEndpoint.ServiceUrlDomain + "/fundBalanceAccount/" + account_parsed.AccountID
	res_interface_data, err := s.restApiService.CallRestApi(ctx,"GET", path, &s.appServer.RestEndpoint.XApigwId, nil)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	var account_balance_parsed core.AccountBalance
	err = mapstructure.Decode(res_interface_data, &account_balance_parsed)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		span.RecordError(err)
		return nil, errors.New(err.Error())
    }

	// Update the status payment
	if (account_balance_parsed.Amount < payment.Amount) {
		res.Status = "DECLINED:NO-FUND"
	} else {
		res.Status = "APPROVED"
	}
	
	res_update, err := s.workerRepo.Update(ctx, tx, res)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	if res_update == 0 {
		span.RecordError(erro.ErrUpdate)
		err = erro.ErrUpdate
		return nil, err
	}

	span.AddEvent("Release Transaction - unlock")

	return res, nil
}

func (s WorkerService) PayWithCheckFraud(ctx context.Context, payment *core.Payment) (*core.Payment, error){
	childLogger.Debug().Msg("PayWithCheckFraud")
	
	span := lib.Span(ctx, "service.payWithCheckFraud")	

	tx, conn, err := s.workerRepo.StartTx(ctx)
	if err != nil {
		return nil, err
	}
	
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
		s.workerRepo.ReleaseTx(conn)
		span.End()
	}()

	if (payment.CardType != "CREDIT") && (payment.CardType != "DEBIT") {
		span.RecordError(erro.ErrCardTypeInvalid)
		return nil, erro.ErrCardTypeInvalid
	}
	// Read Card
	card := core.Card{}
	card.CardNumber = payment.CardNumber
	res_interface_card, err := s.workerRepo.GetCard(ctx, &card)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	var card_parsed core.Card
	err = mapstructure.Decode(res_interface_card, &card_parsed)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		span.RecordError(err)
		return nil, errors.New(err.Error())
    }

	// Read Terminal
	terminal := core.Terminal{}
	terminal.Name = payment.TerminalName
	res_interface_term, err := s.workerRepo.GetTerminal(ctx, &terminal)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	var terminal_parsed core.Terminal
	err = mapstructure.Decode(res_interface_term, &terminal_parsed)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		span.RecordError(err)
		return nil, errors.New(err.Error())
    }

	// Get Account for Just for Check
	path := s.appServer.RestEndpoint.ServiceUrlDomain + "/getId/" + strconv.Itoa(card_parsed.FkAccountID)
	res_interface_acc, err := s.restApiService.CallRestApi(ctx,	"GET", path, &s.appServer.RestEndpoint.XApigwId, nil)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	var account_parsed core.Account
	jsonString, err := json.Marshal(res_interface_acc)
	if err != nil {
		childLogger.Error().Err(err).Msg("Error Marshal")
		span.RecordError(err)
		return nil, err
	}
	json.Unmarshal(jsonString, &account_parsed)
	span.AddEvent("Begin Transaction - lock")
	
	payment.FkCardID = card_parsed.ID
	payment.FkTerminalId = terminal_parsed.ID
	payment.Status = "PENDING"
	res, err := s.workerRepo.Add(ctx, tx ,payment)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Get Fund
	path = s.appServer.RestEndpoint.ServiceUrlDomain + "/fundBalanceAccount/" + account_parsed.AccountID
	res_interface_data, err := s.restApiService.CallRestApi(ctx,"GET", path, &s.appServer.RestEndpoint.XApigwId, nil)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	var account_balance_parsed core.AccountBalance
	err = mapstructure.Decode(res_interface_data, &account_balance_parsed)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		span.RecordError(err)
		return nil, errors.New(err.Error())
    }

	// Get Payment Feature for ML Fraud xgboost Grpc
	payment_fraud := core.PaymentFraud{}
	
	res_pay_fraud, err := s.workerRepo.GetPaymentFraudFeature(ctx, payment)
	if err != nil {
		switch err {
			case erro.ErrNotFound:
				payment_fraud.CardNumber = payment.CardNumber
				payment_fraud.TerminalName = payment.TerminalName
				payment_fraud.MCC = payment.MCC
				payment_fraud.CoordX = int32(terminal_parsed.CoordX)
				payment_fraud.CoordY = int32(terminal_parsed.CoordY)
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

	childLogger.Debug().Interface("===> res_pay_fraud :", res_pay_fraud).Msg("")

	res_svc_fraud, err := s.CheckPaymentFraudGrpc(ctx, &payment_fraud)
    if err != nil {
		childLogger.Error().Err(err).Msg("error CheckPaymentFraudGrpc")
		span.RecordError(err)
		return nil, errors.New(err.Error())
    }

	var parse_paymentFraud core.PaymentFraud
	jsonString, err = json.Marshal(res_svc_fraud)
	if err != nil {
		childLogger.Error().Err(err).Msg("Error Marshal")
		span.RecordError(err)
		return nil, err
	}
	json.Unmarshal(jsonString, &parse_paymentFraud)

	childLogger.Debug().Interface("*#########> parse_paymentFraud :", parse_paymentFraud).Msg("")
	res.Fraud = parse_paymentFraud.Fraud

	// Get Payment ML Anomaly
	path = s.appServer.RestEndpoint.GatewayMlHost + "/payment/anomaly"
	res_interface_anomaly, err := s.restApiService.CallRestApi(ctx,	"POST",	path, &s.appServer.RestEndpoint.XApigwIdMl,payment_fraud)
	/*res_interface_anomaly, err := s.restApiService.PostData(ctx, 
													s.appServer.RestEndpoint.GatewayMlHost,
													s.appServer.RestEndpoint.ServerHost, 
													s.appServer.RestEndpoint.XApigwIdMl,
													"/payment/anomaly", 
													payment_fraud)*/
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	childLogger.Debug().Interface("*#########> res_interface_anomaly :", res_interface_anomaly).Msg("")
	
	jsonString, err = json.Marshal(res_interface_anomaly)
	if err != nil {
		childLogger.Error().Err(err).Msg("Error Marshal")
		span.RecordError(err)
		return nil, err
	}

	var result map[string]interface{}
	json.Unmarshal(jsonString, &result)
	res.Anomaly = result["score"].(float64)

	// Update the status payment
	if (account_balance_parsed.Amount < payment.Amount) {
		res.Status = "DECLINED:NO-FUND"
	} else {
		res.Status = "APPROVED"
	}
	res_update, err := s.workerRepo.Update(ctx, tx, res)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	if res_update == 0 {
		span.RecordError(erro.ErrUpdate)
		err = erro.ErrUpdate
		return nil, err
	}

	span.AddEvent("Release Transaction - unlock")

	return res, nil
}