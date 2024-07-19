package core

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