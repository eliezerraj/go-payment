package postgre

import (
	"context"
	"time"
	"errors"

	_ "github.com/lib/pq"
	"database/sql"

	"github.com/go-payment/internal/core"
	"github.com/aws/aws-xray-sdk-go/xray"

)

func (w WorkerRepository) Get(ctx context.Context, payment core.Payment) (*core.Payment, error){
	childLogger.Debug().Msg("Get")

	_, root := xray.BeginSubsegment(ctx, "Repository.Get")
	defer func() {
		root.Close(nil)
	}()

	client:= w.databaseHelper.GetConnection()
	
	result_query := core.Payment{}
	rows, err := client.QueryContext(ctx, `SELECT 	id, 
													fk_account_id, 
													card_number, 
													card_type,
													status,
													currency, 
													amount,
													mcc,
													create_at,
													update_at,
													tenant_id
											FROM payment 
											WHERE id =$1 `, payment.ID)
	if err != nil {
		childLogger.Error().Err(err).Msg("SELECT statement")
		return nil, errors.New(err.Error())
	}

	for rows.Next() {
		err := rows.Scan( 	&result_query.ID, 
							&result_query.FkAccountID, 
							&result_query.CardNumber, 
							&result_query.CardType, 
							&result_query.Status, 
							&result_query.Currency,
							&result_query.Amount,
							&result_query.MCC,
							&result_query.CreateAt,
							&result_query.UpdateAt,
							&result_query.TenantID,
						)
		if err != nil {
			childLogger.Error().Err(err).Msg("Scan statement")
			return nil, errors.New(err.Error())
        }
	}

	defer rows.Close()
	return &result_query , nil
}

func (w WorkerRepository) Add(ctx context.Context, tx *sql.Tx, payment core.Payment) (*core.Payment, error){
	childLogger.Debug().Msg("Pay")

	_, root := xray.BeginSubsegment(ctx, "Repository.Add")
	defer func() {
		root.Close(nil)
	}()
	
	stmt, err := tx.Prepare(`INSERT INTO payment ( 	fk_account_id, 
													card_number,
													card_type,
													status, 
													currency,
													amount,
													mcc,
													create_at,
													tenant_id) 
									VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id `)
	if err != nil {
		childLogger.Error().Err(err).Msg("INSERT statement")
		return nil, errors.New(err.Error())
	}

	var id int
	var_createAt := time.Now()
	err = stmt.QueryRowContext(ctx, 
								payment.FkAccountID,
								payment.CardNumber,
								payment.CardType,
								payment.Status,
								payment.Currency,
								payment.Amount,
								payment.MCC,
								var_createAt,
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

	_, root := xray.BeginSubsegment(ctx, "Repository.Update")
	defer func() {
		root.Close(nil)
	}()
	
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