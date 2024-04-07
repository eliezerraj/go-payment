package postgre

import (
	"context"
	"time"
	"errors"

	_ "github.com/lib/pq"
	"database/sql"

	"github.com/go-payment/internal/erro"
	"github.com/go-payment/internal/core"

	"go.opentelemetry.io/otel"
)

func (w WorkerRepository) Get(ctx context.Context, payment core.Payment) (*core.Payment, error){
	childLogger.Debug().Msg("Get")

	ctx, repospan := otel.Tracer("go-payment").Start(ctx,"repo.Get")
	defer repospan.End()

	client:= w.databaseHelper.GetConnection()
	
	result_query := core.Payment{}
	rows, err := client.QueryContext(ctx, `	SELECT id, 
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
											WHERE id =$1 `, payment.ID)
	if err != nil {
		childLogger.Error().Err(err).Msg("SELECT statement")
		return nil, errors.New(err.Error())
	}

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
			return nil, errors.New(err.Error())
        }
		return &result_query , nil
	}

	defer rows.Close()
	return nil, erro.ErrNotFound
}

func (w WorkerRepository) Add(ctx context.Context, tx *sql.Tx, payment core.Payment) (*core.Payment, error){
	childLogger.Debug().Msg("Add")

	childLogger.Debug().Interface("payment: ",payment).Msg("*****")

	ctx, repospan := otel.Tracer("go-payment").Start(ctx,"repo.Add")
	defer repospan.End()

	var_createAt := time.Now()
	if payment.PaymentAt.IsZero(){
		payment.PaymentAt = time.Now()
	}
	stmt, err := tx.Prepare(`INSERT INTO payment ( 	fk_card_id, 
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
									VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id `)
	if err != nil {
		childLogger.Error().Err(err).Msg("INSERT statement")
		return nil, errors.New(err.Error())
	}

	var id int
	err = stmt.QueryRowContext(ctx, 
								payment.FkCardID,
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
								var_createAt,
								payment.Fraud,
								payment.TenantID).Scan(&id)
	if err != nil {
		childLogger.Error().Err(err).Msg("Exec statement")
		return nil, errors.New(err.Error())
	}
	defer stmt.Close()

	payment.ID = id
	payment.CreateAt = var_createAt

	return &payment , nil
}

func (w WorkerRepository) Update(ctx context.Context, tx *sql.Tx, payment core.Payment) (int64, error){
	childLogger.Debug().Msg("Update")

	ctx, repospan := otel.Tracer("go-payment").Start(ctx,"repo.Update")
	defer repospan.End()
	
	stmt, err := tx.Prepare(`update payment
							set status = $2,
								update_at = $3
							where id = $1 `)
	if err != nil {
		childLogger.Error().Err(err).Msg("UPDATE statement")
		return 0, errors.New(err.Error())
	}

	result, err := stmt.ExecContext(ctx,	
									payment.ID,
									payment.Status,
									time.Now(),
								)
	if err != nil {
		childLogger.Error().Err(err).Msg("Exec statement")
		return 0, errors.New(err.Error())
	}

	rowsAffected, _ := result.RowsAffected()
	childLogger.Debug().Int("rowsAffected : ",int(rowsAffected)).Msg("")

	defer stmt.Close()
	return rowsAffected , nil
}

func (w WorkerRepository) GetPaymentFraudFeature(ctx context.Context, payment core.Payment) (*core.PaymentFraud, error){
	childLogger.Debug().Msg("GetPaymentFraudFeature")
	childLogger.Debug().Interface("===>payment :", payment).Msg("")

	ctx, repospan := otel.Tracer("go-payment").Start(ctx,"repo.GetPaymentFraudFeature")
	defer repospan.End()

	client:= w.databaseHelper.GetConnection()
	
	result_query := core.PaymentFraud{}
	rows := client.QueryRowContext(ctx,  `select 	p.payment_at,
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
												limit 1 `, payment.CardNumber)

	err := rows.Scan(&result_query.PaymentAt, 
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
		return nil, erro.ErrNotFound
	}

	return &result_query , nil
}

func (w WorkerRepository) GetPaymentFraudFeature2(ctx context.Context, payment core.Payment) (*core.PaymentFraud, error){
	childLogger.Debug().Msg("GetPaymentFraudFeature")
	childLogger.Debug().Interface("===>payment :", payment).Msg("")

	ctx, repospan := otel.Tracer("go-payment").Start(ctx,"repo.GetPaymentFraudFeature")
	defer repospan.End()

	client:= w.databaseHelper.GetConnection()
	
	result_query := core.PaymentFraud{}
	rows, err := client.QueryContext(ctx,  `select 	p.payment_at,
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
												limit 1 `, payment.CardNumber)
	if err != nil {
		childLogger.Error().Err(err).Msg("SELECT statement")
		return nil, errors.New(err.Error())
	}

	childLogger.Debug().Msg("GetPaymentFraudFeature - step 1")
	childLogger.Debug().Interface("GetPaymentFraudFeature rows : ", rows).Msg("")

	childLogger.Debug().Msg("GetPaymentFraudFeature - step 2")
	
	for rows.Next() {
		childLogger.Debug().Msg("GetPaymentFraudFeature - step 3")

		err := rows.Scan( 	&result_query.PaymentAt, 
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
							&result_query.TimeBtwTx,
						)
	
		if err != nil {
			childLogger.Error().Err(err).Msg("Scan statement")
			return nil, errors.New(err.Error())
        }
		return &result_query , nil
	}

	defer rows.Close()
	return nil, erro.ErrNotFound
}