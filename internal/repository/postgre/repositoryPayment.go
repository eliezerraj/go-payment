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
	childLogger.Debug().Msg("Pay")

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