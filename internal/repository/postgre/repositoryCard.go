package postgre

import (
	"context"
	"errors"
	_ "github.com/lib/pq"
	"github.com/go-payment/internal/erro"
	"github.com/go-payment/internal/core"
	"github.com/go-payment/internal/lib"
)

func (w WorkerRepository) GetCard(ctx context.Context, card core.Card) (*core.Card, error){
	childLogger.Debug().Msg("GetCard")
	//childLogger.Debug().Interface("card: ",card).Msg("*****")

	span := lib.Span(ctx, "repo.getCard")	
    defer span.End()

	client:= w.databaseHelper.GetConnection()
	
	result_query := core.Card{}
	rows, err := client.QueryContext(ctx, `	SELECT id, 
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
											WHERE card_number =$1 `, card.CardNumber)
	if err != nil {
		childLogger.Error().Err(err).Msg("SELECT statement")
		span.RecordError(err)
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

func (w WorkerRepository) GetTerminal(ctx context.Context, terminal core.Terminal) (*core.Terminal, error){
	childLogger.Debug().Msg("GetTerminal")

	span := lib.Span(ctx, "repo.getTerminal")	
    defer span.End()

	client:= w.databaseHelper.GetConnection()
	
	result_query := core.Terminal{}
	rows, err := client.QueryContext(ctx, `SELECT id, 
													terminal_name, 
													coord_x, 
													coord_y, 
													status, 
													create_at, 
													update_at
											FROM terminal
											WHERE terminal_name =$1 `, terminal.Name)
	if err != nil {
		childLogger.Error().Err(err).Msg("SELECT statement")
		span.RecordError(err)
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
