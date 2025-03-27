package model

import (
	"time"
	go_core_pg "github.com/eliezerraj/go-core/database/pg"
	go_core_observ "github.com/eliezerraj/go-core/observability" 
)

type AppServer struct {
	InfoPod 		*InfoPod 					`json:"info_pod"`
	Server     		*Server     				`json:"server"`
	ConfigOTEL		*go_core_observ.ConfigOTEL	`json:"otel_config"`
	DatabaseConfig	*go_core_pg.DatabaseConfig  `json:"database"`
	ApiService 		[]ApiService 				`json:"api_endpoints"`
}

type InfoPod struct {
	PodName				string 	`json:"pod_name"`
	ApiVersion			string 	`json:"version"`
	OSPID				string 	`json:"os_pid"`
	IPAddress			string 	`json:"ip_address"`
	AvailabilityZone 	string 	`json:"availabilityZone"`
	IsAZ				bool   	`json:"is_az"`
	Env					string `json:"enviroment,omitempty"`
	AccountID			string `json:"account_id,omitempty"`
}

type Server struct {
	Port 			int `json:"port"`
	ReadTimeout		int `json:"readTimeout"`
	WriteTimeout	int `json:"writeTimeout"`
	IdleTimeout		int `json:"idleTimeout"`
	CtxTimeout		int `json:"ctxTimeout"`
}

type ApiService struct {
	Name			string `json:"name_service"`
	Url				string `json:"url"`
	Method			string `json:"method"`
	Header_x_apigw_api_id	string `json:"x-apigw-api-id"`
}

type MessageRouter struct {
	Message			string `json:"message"`
}

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

type MovimentAccount struct {
	AccountBalance					*AccountBalance	`json:"account_balance,omitempty"`
	AccountBalanceStatementCredit	float64			`json:"account_balance_statement_credit,omitempty"`
	AccountBalanceStatementDebit	float64			`json:"account_balance_statement_debit,omitempty"`
	AccountBalanceStatementTotal	float64			`json:"account_balance_debit.debit_total,omitempty"`
}

type AccountBalance struct {
	ID				int			`json:"id,omitempty"`
	AccountID		string		`json:"account_id,omitempty"`
	FkAccountID		int			`json:"fk_account_id,omitempty"`
	Currency		string  	`json:"currency,omitempty"`
	Amount			float64 	`json:"amount"`
	TenantID		string  	`json:"tenant_id,omitempty"`
	CreateAt		time.Time 	`json:"create_at,omitempty"`
	UpdateAt		*time.Time 	`json:"update_at,omitempty"`
	UserLastUpdate	*string  	`json:"user_last_update,omitempty"`
	JwtId			*string  	`json:"jwt_id,omitempty"`
	RequestId		*string  	`json:"request_id,omitempty"`
	TransactionID	*string  	`json:"transaction_id,omitempty"`
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