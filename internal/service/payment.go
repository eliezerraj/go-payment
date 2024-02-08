package service

import (
	//"os"
	//"time"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	//"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/go-payment/internal/core"
	"github.com/go-payment/internal/erro"
	"github.com/go-payment/internal/adapter/restapi"
	"github.com/go-payment/internal/repository/postgre"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/go-payment/internal/adapter/grpc"

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
	
	_, root := xray.BeginSubsegment(ctx, "Service.Get")
	defer root.Close(nil)

	res, err := s.workerRepository.Get(ctx, payment)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s WorkerService) Pay(ctx context.Context, payment core.Payment) (*core.Payment, error){
	childLogger.Debug().Msg("Pay")
	childLogger.Debug().Interface("payment:",payment).Msg("")
	
	_, root := xray.BeginSubsegment(ctx, "Service.Pay")

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
		root.Close(nil)
	}()

	if (payment.CardType != "CREDIT") && (payment.CardType != "DEBIT") {
		return nil, erro.ErrCardTypeInvalid
	}
	// Get Account
	rest_interface_data, err := s.restapi.GetData(ctx, s.restapi.ServerUrlDomain, s.restapi.XApigwId,"/get", payment.AccountID)
	if err != nil {
		return nil, err
	}
	var account_parsed core.Account
	err = mapstructure.Decode(rest_interface_data, &account_parsed)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		return nil, errors.New(err.Error())
    }

	payment.FkAccountID = account_parsed.ID
	payment.Status = "PENDING"
	res, err := s.workerRepository.Add(ctx, tx ,payment)
	if err != nil {
		return nil, err
	}

	// Get Fund
	rest_interface_data, err = s.restapi.GetData(ctx, s.restapi.ServerUrlDomain, s.restapi.XApigwId,"/fundBalanceAccount", payment.AccountID)
	if err != nil {
		return nil, err
	}
	var account_balance_parsed core.AccountBalance
	err = mapstructure.Decode(rest_interface_data, &account_balance_parsed)
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
		return nil, erro.ErrUpdate
	}

	return res, nil
}