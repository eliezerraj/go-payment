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

type HttpAppServer struct {
	InfoPod 	*InfoPod 		`json:"info_pod"`
	Server     	*Server     	`json:"server"`
	Cert		*Cert			`json:"cert"`
}

type InfoPod struct {
	PodName				string `json:"pod_name,omitempty"`
	ApiVersion			string `json:"version,omitempty"`
	OSPID				string `json:"os_pid,omitempty"`
	IPAddress			string `json:"ip_address,omitempty"`
	AvailabilityZone 	string `json:"availabilityZone,omitempty"`
	Database			*DatabaseRDS `json:"database,omitempty"`
	GrpcHost			string `json:"grpc_host,omitempty"`
	GatewayMlHost		string `json:"gateway_ml_host,omitempty"`
	isTLS				string `json:"is_tls,omitempty"`
	OtelExportEndpoint	string `json:"otel_export_endpoint,omitempty"`
}

type Server struct {
	Port 			int `json:"port"`
	ReadTimeout		int `json:"readTimeout"`
	WriteTimeout	int `json:"writeTimeout"`
	IdleTimeout		int `json:"idleTimeout"`
	CtxTimeout		int `json:"ctxTimeout"`
}
