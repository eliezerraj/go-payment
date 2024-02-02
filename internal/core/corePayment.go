package core

import (
	"time"

)

type Payment struct {
	ID				int			`json:"id,omitempty"`
	AccountID		string		`json:"account_id,omitempty"`
	FkAccountID		int			`json:"fk_account_id,omitempty"`
	CardNumber		string		`json:"card_number,omitempty"`
	CardType		string  	`json:"card_type,omitempty"`
	Status			string  	`json:"status,omitempty"`
	Currency		string  	`json:"currency,omitempty"`
	Amount			float64 	`json:"amount,omitempty"`
	MCC				string  	`json:"mcc,omitempty"`
	CreateAt		time.Time 	`json:"create_at,omitempty"`
	UpdateAt		*time.Time 	`json:"update_at,omitempty"`
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