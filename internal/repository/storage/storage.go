package storage

import (
	"context"
	"time"
	"errors"

	"github.com/go-payment/internal/repository/pg"
	"github.com/go-payment/internal/core"
	"github.com/go-payment/internal/lib"
	"github.com/go-payment/internal/erro"

	"github.com/rs/zerolog/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var childLogger = log.With().Str("repository.pg", "storage").Logger()

//-----------------------------------------------
type WorkerRepository struct {
	databasePG pg.DatabasePG
}

func NewWorkerRepository(databasePG pg.DatabasePG) WorkerRepository {
	childLogger.Debug().Msg("NewWorkerRepository")
	return WorkerRepository{
		databasePG: databasePG,
	}
}
//-----------------------------------------------
func (w WorkerRepository) StartTx(ctx context.Context) (pgx.Tx, *pgxpool.Conn, error) {
	childLogger.Debug().Msg("StartTx")

	span := lib.Span(ctx, "storage.StartTx")
	defer span.End()

	conn, err := w.databasePG.Acquire(ctx)
	if err != nil {
		childLogger.Error().Err(err).Msg("Erro Acquire")
		return nil, nil, errors.New(err.Error())
	}

	tenant := ctx.Value("tenant_id").(string)
	res_rls, err := w.SetSessionRLS(ctx, conn, tenant)
	if err != nil {
		childLogger.Error().Err(err).Msg("Erro SetSessionRLS")
		return nil, nil, errors.New(err.Error())
	}
	if res_rls != true{
		childLogger.Error().Err(err).Msg("using RLS error !!!")
	}

	tx, err := conn.Begin(ctx)
    if err != nil {
        return nil, nil ,errors.New(err.Error())
    }

	return tx, conn, nil
}

func (w WorkerRepository) ReleaseTx(connection *pgxpool.Conn) {
	childLogger.Debug().Msg("ReleaseTx")

	defer connection.Release()
}

func (w WorkerRepository) SetSessionRLS(ctx context.Context, conn *pgxpool.Conn, userCredential string) (bool, error) {
	childLogger.Debug().Msg("++++++++++++++++++++++++++++++++")
	childLogger.Debug().Msg("SetSessionRLS")
	
	rows, err := conn.Query(ctx, "SET app.current_tenant = '" + userCredential+ "'")
	if err != nil {
		childLogger.Error().Err(err).Msg("SET SESSION statement ERROR")
		return false, errors.New(err.Error())
	}
	for rows.Next() {
	}
	return true, nil
}
//---------------------------------------------------------------
func (w WorkerRepository) GetCard(ctx context.Context, card *core.Card) (*core.Card, error){
	childLogger.Debug().Msg("GetCard")
	//childLogger.Debug().Interface("card: ",card).Msg("*****")

	span := lib.Span(ctx, "storage.getCard")	
    defer span.End()

	conn, err := w.databasePG.Acquire(ctx)
	if err != nil {
		childLogger.Error().Err(err).Msg("Erro Acquire")
		return nil, errors.New(err.Error())
	}
	defer w.databasePG.Release(conn)

	result_query := core.Card{}
	query :=  `SELECT 	id, 
						fk_account_id, 
						card_number, 
						card_type, 
						card_model, 
						card_pin, 
						status, 
						expire_at, 
						create_at, 
						update_at, 
						tenant_id
				FROM card
				WHERE card_number =$1`

	rows, err := conn.Query(ctx, query, card.CardNumber)
	if err != nil {
		childLogger.Error().Err(err).Msg("Query statement")
		return nil, errors.New(err.Error())
	}

	for rows.Next() {
		err := rows.Scan( 	&result_query.ID, 
							&result_query.FkAccountID, 
							&result_query.CardNumber, 
							&result_query.Type, 
							&result_query.Model,
							&result_query.Pin,
							&result_query.Status,
							&result_query.ExpireAt,
							&result_query.CreateAt,
							&result_query.UpdateAt,
							&result_query.TenantID,
						)
		if err != nil {
			childLogger.Error().Err(err).Msg("Scan statement")
			span.RecordError(err)
			return nil, errors.New(err.Error())
        }
		return &result_query , nil
	}

	defer rows.Close()
	return nil, erro.ErrNotFound
}

func (w WorkerRepository) GetTerminal(ctx context.Context, terminal *core.Terminal) (*core.Terminal, error){
	childLogger.Debug().Msg("GetTerminal")

	span := lib.Span(ctx, "storage.getTerminal")	
    defer span.End()

	conn, err := w.databasePG.Acquire(ctx)
	if err != nil {
		childLogger.Error().Err(err).Msg("Erro Acquire")
		return nil, errors.New(err.Error())
	}
	defer w.databasePG.Release(conn)
	
	result_query := core.Terminal{}
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
		childLogger.Error().Err(err).Msg("Query statement")
		return nil, errors.New(err.Error())
	}

	for rows.Next() {
		err := rows.Scan( 	&result_query.ID, 
							&result_query.Name, 
							&result_query.CoordX, 
							&result_query.CoordY, 
							&result_query.Status,
							&result_query.CreateAt,
							&result_query.UpdateAt,
						)
		if err != nil {
			childLogger.Error().Err(err).Msg("Scan statement")
			span.RecordError(err)
			return nil, errors.New(err.Error())
        }
		return &result_query , nil
	}

	defer rows.Close()
	return nil, erro.ErrNotFound
}
//-------------------------------------------------------------------
func (w WorkerRepository) Get(ctx context.Context, payment *core.Payment) (*core.Payment, error){
	childLogger.Debug().Msg("Get")

	span := lib.Span(ctx, "storage.get")	
    defer span.End()

	conn, err := w.databasePG.Acquire(ctx)
	if err != nil {
		childLogger.Error().Err(err).Msg("Erro Acquire")
		return nil, errors.New(err.Error())
	}
	defer w.databasePG.Release(conn)

	tenant := ctx.Value("tenant_id").(string)
	res_rls, err := w.SetSessionRLS(ctx, conn, tenant)
	if err != nil {
		childLogger.Error().Err(err).Msg("Erro SetSessionRLS")
		return nil, errors.New(err.Error())
	}
	if res_rls != true{
		childLogger.Error().Err(err).Msg("using RLS error !!!")
	}
	
	result_query := core.Payment{}
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
		childLogger.Error().Err(err).Msg("SELECT statement")
		span.RecordError(err)
		return nil, errors.New(err.Error())
	}

	lib.Event(span, query)
	//span.AddEvent("Executing SQL query", trace.WithAttributes(attribute.String("db.statement", query)))

	for rows.Next() {
		err := rows.Scan( 	&result_query.ID, 
							&result_query.FkCardID, 
							&result_query.CardNumber, 
							&result_query.FkTerminalId, 
							&result_query.CardType, 
							&result_query.CardMode,
							&result_query.PaymentAt,
							&result_query.MCC,
							&result_query.Status,							
							&result_query.Currency,
							&result_query.Amount,
							&result_query.CreateAt,
							&result_query.UpdateAt,
							&result_query.Fraud,
							&result_query.TenantID,
						)
		if err != nil {
			childLogger.Error().Err(err).Msg("Scan statement")
			span.RecordError(err)
			return nil, errors.New(err.Error())
        }
		return &result_query , nil
	}

	defer rows.Close()
	return nil, erro.ErrNotFound
}

func (w WorkerRepository) Add(ctx context.Context, tx pgx.Tx, payment *core.Payment) (*core.Payment, error){
	childLogger.Debug().Msg("Add")
	childLogger.Debug().Interface("payment: ",payment).Msg("*****")

	span := lib.Span(ctx, "storage.add")	
    defer span.End()

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

	payment.CreateAt = time.Now()
	if payment.PaymentAt.IsZero(){
		payment.PaymentAt = payment.CreateAt
	}

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

	payment.ID = id
	return payment , nil
}

func (w WorkerRepository) Update(ctx context.Context, tx pgx.Tx, payment *core.Payment) (int64, error){
	childLogger.Debug().Msg("Update")

	span := lib.Span(ctx, "storage.update")	
    defer span.End()
	
	query := `update payment
				set status = $2,
					update_at = $3
				where id = $1`

	row, err := tx.Exec(ctx, query,	payment.ID,
									payment.Status,
									time.Now())
	if err != nil {
		childLogger.Error().Err(err).Msg("Exec statement")
		span.RecordError(err)
		return 0, errors.New(err.Error())
	}

	childLogger.Debug().Int("rowsAffected : ", int(row.RowsAffected())).Msg("")

	return row.RowsAffected(), nil
}

func (w WorkerRepository) GetPaymentFraudFeature(ctx context.Context, payment *core.Payment) (*core.PaymentFraud, error){
	childLogger.Debug().Msg("GetPaymentFraudFeature")
	childLogger.Debug().Interface("===>payment :", payment).Msg("")

	span := lib.Span(ctx, "storage.getPaymentFraudFeature")	
    defer span.End()

	conn, err := w.databasePG.Acquire(ctx)
	if err != nil {
		childLogger.Error().Err(err).Msg("Erro Acquire")
		return nil, errors.New(err.Error())
	}
	defer w.databasePG.Release(conn)

	result_query := core.PaymentFraud{}

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
						(select to_char(avg(p1.amount),'FM999999999.00') as avg_1d
							from payment p1
							where p1.card_number = p.card_number
							and p1.payment_at::date = p.payment_at::date
							group by p1.card_number, p1.payment_at::date),
						(select count(*) as tx_7d
							from payment p1
							where p1.card_number = p.card_number
							and p1.payment_at::date between (p.payment_at::date - interval '6 days') and p.payment_at::date
							group by p1.card_number	),
						(select to_char(avg(p1.amount),'FM999999999.00') as avg_7d
							from payment p1
							where p1.card_number = p.card_number
							and p1.payment_at::date between (p.payment_at::date - interval '6 days') and p.payment_at::date
							group by p1.card_number	),
						(select count(*) as tx_30d
							from payment p1
							where p1.card_number = p.card_number
							and p1.payment_at::date between (p.payment_at::date - interval '31 days') and p.payment_at::date
							group by p1.card_number	),
						(select to_char(avg(p1.amount),'FM999999999.00') as avg_30d
							from payment p1
							where p1.card_number = p.card_number
							and p1.payment_at::date between (p.payment_at::date - interval '31 days') and p.payment_at::date
							group by p1.card_number	),
							to_char(coalesce ( (select extract(epoch from p.payment_at - p1.payment_at )
													from payment p1 
													where p1.card_number = p.card_number
													and p1.payment_at < p.payment_at
													order by p1.payment_at desc
													limit 1),0),'FM999999999') as time_btw_cc_tx
												from payment p,
												terminal t
				where p.fk_terminal_id = t.id
				and p.card_number = $1
				order by p.payment_at desc
				limit 1 `

	rows, err := conn.Query(ctx, query, payment.CardNumber)
	if err != nil {
		childLogger.Error().Err(err).Msg("Query statement")
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	err = rows.Scan(&result_query.PaymentAt, 
					&result_query.CardModel, 
					&result_query.CardType, 
					&result_query.CoordX, 
					&result_query.CoordY, 
					&result_query.Tx1Day,
					&result_query.Avg1Day,
					&result_query.Tx7Day,
					&result_query.Avg7Day,							
					&result_query.Tx30Day,
					&result_query.Avg30Day,
					&result_query.TimeBtwTx)
	if err != nil {
		childLogger.Error().Err(err).Msg("Scan statement")
		span.RecordError(err)
		return nil, erro.ErrNotFound
	}

	return &result_query , nil
}