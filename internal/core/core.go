package core

import(
	"time"
)

type DatabaseRDS struct {
    Host 				string `json:"host"`
    Port  				string `json:"port"`
	Schema				string `json:"schema"`
	DatabaseName		string `json:"databaseName"`
	User				string `json:"user"`
	Password			string `json:"password"`
	Db_timeout			int	`json:"db_timeout"`
	Postgres_Driver		string `json:"postgres_driver"`
}

type AppServer struct {
	InfoPod 		*InfoPod 		`json:"info_pod"`
	Server     		*Server     	`json:"server"`
	Database		*DatabaseRDS	`json:"database"`
	RestEndpoint	*RestEndpoint	`json:"rest_endpoint"`
	AuthUser		*AuthUser		`json:"auth_user"`
	ConfigOTEL		*ConfigOTEL		`json:"otel_config"`
	AwsServiceConfig 	*AwsServiceConfig	`json:"aws_service_config"`
	RestApiCallData 	*RestApiCallData `json:"rest_api_call_dsa_data"`
}

type InfoPod struct {
	PodName				string `json:"pod_name,omitempty"`
	ApiVersion			string `json:"version,omitempty"`
	OSPID				string `json:"os_pid,omitempty"`
	IPAddress			string `json:"ip_address,omitempty"`
	AvailabilityZone 	string `json:"availabilityZone,omitempty"`
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

type RestEndpoint struct {
	ServiceUrlDomain 	string `json:"service_url_domain"`
	XApigwId			string `json:"xApigwId"`
	CaCert				*Cert `json:"ca_cert"`
	GatewayMlHost		string `json:"gateway_ml_host,omitempty"`
	XApigwIdMl			string `json:"xApigwIdMl"`
	GrpcHost			string `json:"grpc_host,omitempty"`
	ServerHost			string `json:"server_host_localhost,omitempty"`
	AuthUrlDomain		string `json:"auth_url_domain,omitempty"`
	Token				string `json:"token,omitempty"`	
}

type Cert struct {
	CaAccountPEM 		[]byte
	CaFraudPEM 			[]byte  	 		
}

type ConfigOTEL struct {
	OtelExportEndpoint		string
	TimeInterval            int64    `mapstructure:"TimeInterval"`
	TimeAliveIncrementer    int64    `mapstructure:"RandomTimeAliveIncrementer"`
	TotalHeapSizeUpperBound int64    `mapstructure:"RandomTotalHeapSizeUpperBound"`
	ThreadsActiveUpperBound int64    `mapstructure:"RandomThreadsActiveUpperBound"`
	CpuUsageUpperBound      int64    `mapstructure:"RandomCpuUsageUpperBound"`
	SampleAppPorts          []string `mapstructure:"SampleAppPorts"`
}

type AuthUser struct {
	User 		string `json:"user,omitempty"`
	Password 	string `json:"password,omitempty"`
	Token		string `json:"token,omitempty"`	
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

type AwsServiceConfig struct {
	AwsRegion				string	`json:"aws_region"`
	ServiceUrlJwtSA 		string	`json:"service_url_jwt_sa"`
	SecretJwtSACredential 	string	`json:"secret_jwt_credential"`
	UsernameJwtDA			string	`json:"username_jwt_sa"`
	PasswordJwtDA			string	`json:"password_jwt_sa"`
}

type TokenSA struct {
	Token string `json:"token,omitempty"`
	Err   error
}

type RestApiCallData struct {
	Url				string `json:"url"`
	Method			string `json:"method"`
	X_Api_Id		*string `json:"x-apigw-api-id"`
	UsernameAuth	string `json:"user"`
	PasswordAuth 	string `json:"password"`
}