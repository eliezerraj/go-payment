package service

import(
	"context"
	"net/http"
	"encoding/json"
	"errors"

	"github.com/go-payment/internal/core/model"
	"github.com/go-payment/internal/core/erro"
	go_core_observ "github.com/eliezerraj/go-core/observability"
	go_core_api "github.com/eliezerraj/go-core/api"
)

var tracerProvider go_core_observ.TracerProvider
var apiService go_core_api.ApiService

func errorStatusCode(statusCode int) error{
	var err error
	switch statusCode {
	case http.StatusUnauthorized:
		err = erro.ErrUnauthorized
	case http.StatusForbidden:
		err = erro.ErrHTTPForbiden
	case http.StatusNotFound:
		err = erro.ErrNotFound
	default:
		err = erro.ErrServer
	}
	return err
}

func (s WorkerService) AddPayment(ctx context.Context, payment *model.Payment) (*model.Payment, error){
	childLogger.Debug().Msg("AddPayment")
	childLogger.Debug().Interface("payment: ",payment).Msg("")

	// Trace
	span := tracerProvider.Span(ctx, "service.AddPayment")
	
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

	// Businness rule
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

	res, err := s.workerRepository.AddPayment(ctx, tx, payment)
	if err != nil {
		return nil, err
	}

	// get fund balance
	res_payload, statusCode, err := apiService.CallApi(ctx,
														s.apiService[1].Url + "/" + res_card.AccountID,
														s.apiService[1].Method,
														&s.apiService[1].Header_x_apigw_api_id,
														nil, 
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

	if (movimentAccount.AccountBalance.Amount < payment.Amount) {
		res.Status = "DECLINED:NO-FUND"
	} else {
		res.Status = "APPROVED"
	}
	// update status payment
	res_update, err := s.workerRepository.UpdatePayment(ctx, tx, res)
	if err != nil {
		return nil, err
	}
	if res_update == 0 {
		err = erro.ErrUpdate
		return nil, err
	}

	return res, nil
}

func (s WorkerService) GetPayment(ctx context.Context, payment *model.Payment) (*model.Payment, error){
	childLogger.Debug().Msg("GetPayment")
	childLogger.Debug().Interface("payment: ",payment).Msg("")

	span := tracerProvider.Span(ctx, "service.GetPayment")
	defer span.End()
	
	res, err := s.workerRepository.GetPayment(ctx, payment)
	if err != nil {
		return nil, err
	}
	return res, nil
}
