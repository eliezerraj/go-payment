package service

import (
	"context"
	"strconv"
	"errors"
	"github.com/rs/zerolog/log"
	"encoding/json"


	"github.com/mitchellh/mapstructure"
	"github.com/go-payment/internal/core"
	"github.com/go-payment/internal/erro"
	"github.com/go-payment/internal/adapter/restapi"
	"github.com/go-payment/internal/repository/postgre"
	"github.com/go-payment/internal/adapter/grpc"
	
	"go.opentelemetry.io/otel"
)

var childLogger = log.With().Str("service", "service").Logger()
type WorkerService struct {
	workerRepository 		*postgre.WorkerRepository
	restapi					*restapi.RestApiSConfig
	grpcClient 				*grpc.GrpcClient
}

func NewWorkerService(	workerRepository 	*postgre.WorkerRepository,
						restapi				*restapi.RestApiSConfig,
						grpcClient 			*grpc.GrpcClient) *WorkerService{
	childLogger.Debug().Msg("NewWorkerService")

	return &WorkerService{
		workerRepository:	workerRepository,
		restapi:			restapi,
		grpcClient: 		grpcClient,
	}
}

func (s WorkerService) SetSessionVariable(ctx context.Context, userCredential string) (bool, error){
	childLogger.Debug().Msg("SetSessionVariable")

	res, err := s.workerRepository.SetSessionVariable(ctx, userCredential)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (s WorkerService) Get(ctx context.Context, payment core.Payment) (*core.Payment, error){
	childLogger.Debug().Msg("Get")
	
	ctx, svcspan := otel.Tracer("go-payment").Start(ctx,"svc.Get")
	defer svcspan.End()

	res, err := s.workerRepository.Get(ctx, payment)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s WorkerService) Pay(ctx context.Context, payment core.Payment) (*core.Payment, error){
	childLogger.Debug().Msg("Pay")
	
	ctx, svcspan := otel.Tracer("go-payment").Start(ctx,"svc.Pay")
	defer svcspan.End()

	tx, err := s.workerRepository.StartTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if (payment.CardType != "CREDIT") && (payment.CardType != "DEBIT") {
		return nil, erro.ErrCardTypeInvalid
	}
	// Read Card
	card := core.Card{}
	card.CardNumber = payment.CardNumber
	res_interface_card, err := s.workerRepository.GetCard(ctx, card)
	if err != nil {
		return nil, err
	}

	var card_parsed core.Card
	err = mapstructure.Decode(res_interface_card, &card_parsed)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		return nil, errors.New(err.Error())
    }

	// Read Terminal
	terminal := core.Terminal{}
	terminal.Name = payment.TerminalName
	res_interface_term, err := s.workerRepository.GetTerminal(ctx,terminal)
	if err != nil {
		return nil, err
	}
	var terminal_parsed core.Terminal
	err = mapstructure.Decode(res_interface_term, &terminal_parsed)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		return nil, errors.New(err.Error())
    }

	// Get Account for Just for Check
	res_interface_acc, err := s.restapi.GetData(ctx, 
												s.restapi.ServerUrlDomain,
												s.restapi.ServerHost,
												s.restapi.XApigwId,
												"/getId", 
												strconv.Itoa(card_parsed.FkAccountID))
	if err != nil {
		return nil, err
	}
	var account_parsed core.Account
	jsonString, err := json.Marshal(res_interface_acc)
	if err != nil {
		childLogger.Error().Err(err).Msg("Error Marshal")
		return nil, err
	}
	json.Unmarshal(jsonString, &account_parsed)
	svcspan.AddEvent("Begin Transaction - lock")
	
	payment.FkCardID = card_parsed.ID
	payment.FkTerminalId = terminal_parsed.ID
	payment.Status = "PENDING"
	res, err := s.workerRepository.Add(ctx, tx ,payment)
	if err != nil {
		return nil, err
	}

	// Get Fund
	res_interface_data, err := s.restapi.GetData(ctx, 
												s.restapi.ServerUrlDomain,
												s.restapi.ServerHost, 
												s.restapi.XApigwId,
												"/fundBalanceAccount", 
												account_parsed.AccountID)
	if err != nil {
		return nil, err
	}
	var account_balance_parsed core.AccountBalance
	err = mapstructure.Decode(res_interface_data, &account_balance_parsed)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		return nil, errors.New(err.Error())
    }

	// Update the status payment
	if (account_balance_parsed.Amount < payment.Amount) {
		res.Status = "DECLINED:NO-FUND"
	} else {
		res.Status = "APPROVED"
	}
	res_update, err := s.workerRepository.Update(ctx, tx ,*res)
	if err != nil {
		return nil, err
	}
	if res_update == 0 {
		err = erro.ErrUpdate
		return nil, err
	}

	svcspan.AddEvent("Release Transaction - unlock")

	return res, nil
}

