package database

import (
	"context"
	"time"
	"errors"
	
	"github.com/go-payment/internal/core/model"
	"github.com/go-payment/internal/core/erro"

	go_core_observ "github.com/eliezerraj/go-core/observability"
	go_core_pg "github.com/eliezerraj/go-core/database/pg"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

var tracerProvider go_core_observ.TracerProvider
var childLogger = log.With().Str("adapter", "database").Logger()

type WorkerRepository struct {
	DatabasePGServer *go_core_pg.DatabasePGServer
}

func NewWorkerRepository(databasePGServer *go_core_pg.DatabasePGServer) *WorkerRepository{
	childLogger.Debug().Msg("NewWorkerRepository")

	return &WorkerRepository{
		DatabasePGServer: databasePGServer,
	}
}

// About add payment
func (w WorkerRepository) AddPayment(ctx context.Context, tx pgx.Tx, payment *model.Payment) (*model.Payment, error){
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Msg("AddPayment")

	// Trace
	span := tracerProvider.Span(ctx, "database.AddPayment")
	defer span.End()

	// Prepare
	payment.CreateAt = time.Now()
	if payment.PaymentAt.IsZero(){
		payment.PaymentAt = payment.CreateAt
	}

	// Query and execute
	query := `INSERT INTO payment (fk_card_id, 
									card_number, 
									fk_terminal_id, 
									terminal_name, 
									card_type, 
									card_model, 
									payment_at, 
									mcc, 
									status, 
									currency, 
									amount, 
									create_at,
									fraud, 
									tenant_id)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`

	row := tx.QueryRow(ctx, query, payment.FkCardID,
									payment.CardNumber,
									payment.FkTerminalId,
									payment.TerminalName,
									payment.CardType,
									payment.CardMode,
									payment.PaymentAt,
									payment.MCC,
									payment.Status,
									payment.Currency,
									payment.Amount,
									payment.CreateAt ,
									payment.Fraud,
									payment.TenantID)

	var id int
	if err := row.Scan(&id); err != nil {
		childLogger.Error().Err(err).Msg("QueryRow INSERT")
		return nil, errors.New(err.Error())
	}

	// set PK
	payment.ID = id
	return payment , nil
}

// About update payment
func (w WorkerRepository) UpdatePayment(ctx context.Context, tx pgx.Tx, payment *model.Payment) (int64, error){
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Msg("UpdatePayment")

	// Trace
	span := tracerProvider.Span(ctx, "database.UpdatePayment")
	defer span.End()

	// Query and execute
	query := `update payment
				set status = $2,
					update_at = $3
				where id = $1`

	row, err := tx.Exec(ctx, query,	payment.ID,
									payment.Status,
									time.Now())
	if err != nil {
		return 0, errors.New(err.Error())
	}
	return row.RowsAffected(), nil
}

// About get payment
func (w WorkerRepository) GetPayment(ctx context.Context, payment *model.Payment) (*model.Payment, error){
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Msg("GetPayment")
	
	// Trace
	span := tracerProvider.Span(ctx, "database.GetPayment")
	defer span.End()

	// Get connection
	conn, err := w.DatabasePGServer.Acquire(ctx)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer w.DatabasePGServer.Release(conn)

	// Prepare
	res_payment := model.Payment{}

	// query and execute
	query := `SELECT id, 
						fk_card_id, 
						card_number, 
						fk_terminal_id, 
						card_type, 
						card_model, 
						payment_at, 
						mcc, 
						status, 
						currency, 
						amount, 
						create_at, 
						update_at,
						fraud, 
						tenant_id
				FROM payment
				WHERE id =$1`

	rows, err := conn.Query(ctx, query, payment.ID)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan( 	&res_payment.ID, 
							&res_payment.FkCardID, 
							&res_payment.CardNumber, 
							&res_payment.FkTerminalId, 
							&res_payment.CardType, 
							&res_payment.CardMode,
							&res_payment.PaymentAt,
							&res_payment.MCC,
							&res_payment.Status,							
							&res_payment.Currency,
							&res_payment.Amount,
							&res_payment.CreateAt,
							&res_payment.UpdateAt,
							&res_payment.Fraud,
							&res_payment.TenantID,
						)
		if err != nil {
			return nil, errors.New(err.Error())
        }
		return &res_payment, nil
	}
	
	return nil, erro.ErrNotFound
}

// About add card
func (w WorkerRepository) GetCard(ctx context.Context, card *model.Card) (*model.Card, error){
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Msg("GetCard")
	
	// Trace
	span := tracerProvider.Span(ctx, "database.GetCard")
	defer span.End()

	// Get connection
	conn, err := w.DatabasePGServer.Acquire(ctx)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer w.DatabasePGServer.Release(conn)

	// prepare
	res_card := model.Card{}

	// query and execute
	query :=  `SELECT 	c.id, 
						c.fk_account_id,
						a.account_id, 
						c.card_number, 
						c.card_type, 
						c.card_model, 
						c.card_pin, 
						c.status, 
						c.expire_at, 
						c.create_at, 
						c.update_at, 
						c.tenant_id
				FROM card c,
					account a 
				WHERE c.card_number = $1
				and a.id = c.fk_account_id`

	rows, err := conn.Query(ctx, query, card.CardNumber)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan( 	&res_card.ID, 
							&res_card.FkAccountID, 
							&res_card.AccountID,
							&res_card.CardNumber, 
							&res_card.Type, 
							&res_card.Model,
							&res_card.Pin,
							&res_card.Status,
							&res_card.ExpireAt,
							&res_card.CreateAt,
							&res_card.UpdateAt,
							&res_card.TenantID,
		)
		if err != nil {
			return nil, errors.New(err.Error())
        }
		return &res_card, nil
	}
	
	return nil, erro.ErrNotFound
}

// About get terminal
func (w WorkerRepository) GetTerminal(ctx context.Context, terminal *model.Terminal) (*model.Terminal, error){
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Msg("GetTerminal")
	
	// Trace
	span := tracerProvider.Span(ctx, "database.GetTerminal")
	defer span.End()

	// Get connection
	conn, err := w.DatabasePGServer.Acquire(ctx)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer w.DatabasePGServer.Release(conn)

	// prepare
	res_terminal := model.Terminal{}

	// query and execute
	query :=  `SELECT 	id, 
						terminal_name, 
						coord_x, 
						coord_y, 
						status, 
						create_at, 
						update_at
				FROM terminal
				WHERE terminal_name =$1`

	rows, err := conn.Query(ctx, query, terminal.Name)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan( 	&res_terminal.ID, 
							&res_terminal.Name, 
							&res_terminal.CoordX, 
							&res_terminal.CoordY, 
							&res_terminal.Status,
							&res_terminal.CreateAt,
							&res_terminal.UpdateAt,
		)
		if err != nil {
			return nil, errors.New(err.Error())
        }
		return &res_terminal, nil
	}
	
	return nil, erro.ErrNotFound
}

// About get payment fraud features
func (w WorkerRepository) GetPaymentFraudFeature(ctx context.Context, payment *model.Payment) (*model.PaymentFraud, error){
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Msg("GetPaymentFraudFeature")
	
	// Trace
	span := tracerProvider.Span(ctx, "database.GetPaymentFraudFeature")
	defer span.End()

	// Get connection
	conn, err := w.DatabasePGServer.Acquire(ctx)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer w.DatabasePGServer.Release(conn)

	// prepare
	result_paymentFraud := model.PaymentFraud{}

	// query and execute
	query :=  `select 	p.payment_at,
						p.card_model, 
						p.card_type,
						t.coord_x,
						t.coord_y,
						(select count(*) as tx_1d
							from payment p1
							where p1.card_number = p.card_number
							and p1.payment_at::date = p.payment_at::date
							group by p1.card_number, p1.payment_at::date),
						(select avg(p1.amount):: float as avg_1d
							from payment p1
							where p1.card_number = p.card_number
							and p1.payment_at::date = p.payment_at::date
							group by p1.card_number, p1.payment_at::date),
						(select count(*) as tx_7d
							from payment p1
							where p1.card_number = p.card_number
							and p1.payment_at::date between (p.payment_at::date - interval '6 days') and p.payment_at::date
							group by p1.card_number	),
						(select avg(p1.amount):: float as avg_7d
							from payment p1
							where p1.card_number = p.card_number
							and p1.payment_at::date between (p.payment_at::date - interval '6 days') and p.payment_at::date
							group by p1.card_number	),
						(select count(*) as tx_30d
							from payment p1
							where p1.card_number = p.card_number
							and p1.payment_at::date between (p.payment_at::date - interval '31 days') and p.payment_at::date
							group by p1.card_number	),
						(select avg(p1.amount):: float as avg_30d
							from payment p1
							where p1.card_number = p.card_number
							and p1.payment_at::date between (p.payment_at::date - interval '31 days') and p.payment_at::date
							group by p1.card_number	),
						coalesce ( (select extract(epoch from p.payment_at - p1.payment_at )
													from payment p1 
													where p1.card_number = p.card_number
													and p1.payment_at < p.payment_at
													order by p1.payment_at desc
													limit 1),0):: int as time_btw_cc_tx 
				from payment p,
					terminal t
				where p.fk_terminal_id = t.id
				and p.card_number = $1
				order by p.payment_at desc
				limit 1 `

	rows, err := conn.Query(ctx, query, payment.CardNumber)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&result_paymentFraud.PaymentAt, 
						&result_paymentFraud.CardModel, 
						&result_paymentFraud.CardType, 
						&result_paymentFraud.CoordX, 
						&result_paymentFraud.CoordY, 
						&result_paymentFraud.Tx1Day,
						&result_paymentFraud.Avg1Day,
						&result_paymentFraud.Tx7Day,
						&result_paymentFraud.Avg7Day,							
						&result_paymentFraud.Tx30Day,
						&result_paymentFraud.Avg30Day,
						&result_paymentFraud.TimeBtwTx)
		if err != nil {
			return nil, errors.New(err.Error())
        }
		return &result_paymentFraud, nil
	}
	
	return nil, erro.ErrNotFound
}