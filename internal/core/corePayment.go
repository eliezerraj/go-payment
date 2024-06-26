package core

import (
	"time"

)

type Payment struct {
	ID				int			`json:"id,omitempty"`
	FkCardID		int			`json:"fk_card_id,omitempty"`
	CardNumber		string		`json:"card_number,omitempty"`
	FkTerminalId	int			`json:"fk_terminal_id,omitempty"`
	TerminalName	string		`json:"terminal_name,omitempty"`
	CardType		string  	`json:"card_type,omitempty"`
	CardMode		string  	`json:"card_model,omitempty"`
	PaymentAt		time.Time	`json:"payment_at,omitempty"`
	MCC				string  	`json:"mcc,omitempty"`
	Status			string  	`json:"status,omitempty"`
	Currency		string  	`json:"currency,omitempty"`
	Amount			float64 	`json:"amount,omitempty"`
	CreateAt		time.Time 	`json:"create_at,omitempty"`
	UpdateAt		*time.Time 	`json:"update_at,omitempty"`
	Fraud			float64	  	`json:"fraud,omitempty"`
	Anomaly			float64	  	`json:"anomaly,omitempty"`
	TenantID		string  	`json:"tenant_id,omitempty"`
}

type Account struct {
	ID				int			`json:"id,omitempty"`
	AccountID		string		`json:"account_id,omitempty"`
	PersonID		string  	`json:"person_id,omitempty"`
	CreateAt		time.Time 	`json:"create_at,omitempty"`
	UpdateAt		*time.Time 	`json:"update_at,omitempty"`
	TenantID		string  	`json:"tenant_id,omitempty"`
	UserLastUpdate	*string  	`json:"user_last_update,omitempty"`
}

type AccountBalance struct {
	ID				int			`json:"id,omitempty"`
	AccountID		string		`json:"account_id,omitempty"`
	FkAccountID		int			`json:"fk_account_id,omitempty"`
	Currency		string  	`json:"currency,omitempty"`
	Amount			float64 	`json:"amount,omitempty"`
	TenantID		string  	`json:"tenant_id,omitempty"`
	CreateAt		time.Time 	`json:"create_at,omitempty"`
	UpdateAt		*time.Time 	`json:"update_at,omitempty"`
	UserLastUpdate	*string  	`json:"user_last_update,omitempty"`
}

type Card struct {
	ID				int			`json:"id,omitempty"`
	FkAccountID		int			`json:"fk_account_id,omitempty"`
	AccountID		string		`json:"account_id,omitempty"`
	CardNumber		string  	`json:"card_number,omitempty"`
	Type			string  	`json:"card_type,omitempty"`
	Model			string  	`json:"card_model,omitempty"`
	Pin				string  	`json:"card_pin,omitempty"`
	Status			string  	`json:"status,omitempty"`
	ExpireAt		time.Time 	`json:"expire_at,omitempty"`
	CreateAt		time.Time 	`json:"create_at,omitempty"`
	UpdateAt		*time.Time 	`json:"update_at,omitempty"`
	TenantID		string  	`json:"tenant_id,omitempty"`
}

type Terminal struct {
	ID				int			`json:"id,omitempty"`
	Name			string		`json:"terminal_name,omitempty"`
	CoordX			float64  	`json:"coord_x,omitempty"`
	CoordY			float64  	`json:"coord_y,omitempty"`
	Status			string  	`json:"status,omitempty"`
	CreateAt		time.Time 	`json:"create_at,omitempty"`
	UpdateAt		*time.Time 	`json:"update_at,omitempty"`
}

type PaymentFraud struct {
	AccountID		string		`json:"account_id,omitempty"`
	CardNumber		string		`json:"card_number,omitempty"`
	TerminalName	string		`json:"terminal_name,omitempty"`
	CoordX			int32		`json:"coord_x,omitempty"`
	CoordY			int32		`json:"coord_y,omitempty"`
	CardType		string  	`json:"card_type,omitempty"`
	CardModel		string  	`json:"card_model,omitempty"`
	PaymentAt		time.Time	`json:"payment_at,omitempty"`
	MCC				string  	`json:"mcc,omitempty"`
	Status			string  	`json:"status,omitempty"`
	Currency		string  	`json:"currency,omitempty"`
	Amount			float64 	`json:"amount"`
	Tx1Day			float64 	`json:"tx_1d"`
	Avg1Day			float64 	`json:"avg_1d"`
	Tx7Day			float64 	`json:"tx_7d"`
	Avg7Day			float64 	`json:"avg_7d"`
	Tx30Day			float64 	`json:"tx_30d"`
	Avg30Day		float64 	`json:"avg_30d"`
	TimeBtwTx		int32 		`json:"time_btw_cc_tx"`
	Fraud			float64	  	`json:"fraud,omitempty"`
}