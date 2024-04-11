package main   

import(
	"time"
	"context"
	"os"
	"strconv"
	"io/ioutil"
	"net"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

    "github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-payment/internal/core"
	"github.com/go-payment/internal/repository/postgre"
	"github.com/go-payment/internal/service"
	"github.com/go-payment/internal/adapter/restapi"
	"github.com/go-payment/internal/handler"
	"github.com/go-payment/internal/adapter/grpc"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

var(
	logLevel 	= 	zerolog.DebugLevel
	noAZ		=	true // set only if you get to split the xray trace per AZ
	serverUrlDomain, serverHost, xApigwId 		string
	infoPod					core.InfoPod
	envDB	 				core.DatabaseRDS
	httpAppServerConfig 	core.HttpAppServer
	server					core.Server
	dataBaseHelper 			postgre.DatabaseHelper
	repoDB					postgre.WorkerRepository
	configOTEL				core.ConfigOTEL
	cert					core.Cert
	caPEM					[]byte
	isTLS		= 	false
)

func getEnv() {
	log.Debug().Msg("getEnv")

	if os.Getenv("API_VERSION") !=  "" {
		infoPod.ApiVersion = os.Getenv("API_VERSION")
	}
	if os.Getenv("POD_NAME") !=  "" {
		infoPod.PodName = os.Getenv("POD_NAME")
	}
	if os.Getenv("PORT") !=  "" {
		intVar, _ := strconv.Atoi(os.Getenv("PORT"))
		server.Port = intVar
	}
	if os.Getenv("GRPC_HOST") !=  "" {	
		infoPod.GrpcHost = os.Getenv("GRPC_HOST")
	}
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") !=  "" {	
		infoPod.OtelExportEndpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	}

	if os.Getenv("DB_HOST") !=  "" {
		envDB.Host = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") !=  "" {
		envDB.Port = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_NAME") !=  "" {	
		envDB.DatabaseName = os.Getenv("DB_NAME")
	}
	if os.Getenv("DB_SCHEMA") !=  "" {	
		envDB.Schema = os.Getenv("DB_SCHEMA")
	}
	if os.Getenv("DB_DRIVER") !=  "" {	
		envDB.Postgres_Driver = os.Getenv("DB_DRIVER")
	}
	
	if os.Getenv("SERVER_URL_DOMAIN") !=  "" {	
		serverUrlDomain = os.Getenv("SERVER_URL_DOMAIN")
	}
	if os.Getenv("X_APIGW_API_ID") !=  "" {	
		xApigwId = os.Getenv("X_APIGW_API_ID")
	}
	if os.Getenv("SERVER_HOST") !=  "" {	
		serverHost = os.Getenv("SERVER_HOST")
	}

	if (os.Getenv("TLS") != "") {
		if (os.Getenv("TLS") != "false") {	
			isTLS = true
		}
	}

	if os.Getenv("NO_AZ") == "false" {	
		noAZ = false
	} else {
		noAZ = true
	}
}

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)
	
	err := godotenv.Load(".env")
	if err != nil {
		log.Info().Err(err).Msg("No .ENV File !!!!")
	}

	getEnv()

	server.ReadTimeout = 60
	server.WriteTimeout = 60
	server.IdleTimeout = 60
	server.CtxTimeout = 60

	configOTEL.TimeInterval = 1
	configOTEL.TimeAliveIncrementer = 1
	configOTEL.TotalHeapSizeUpperBound = 100
	configOTEL.ThreadsActiveUpperBound = 10
	configOTEL.CpuUsageUpperBound = 100
	configOTEL.SampleAppPorts = []string{}

	// Get Database Secrets
	file_user, err := ioutil.ReadFile("/var/pod/secret/username")
	if err != nil {
		log.Error().Err(err).Msg("ERRO FATAL recuperacao secret-user")
		os.Exit(3)
	}
	file_pass, err := ioutil.ReadFile("/var/pod/secret/password")
	if err != nil {
		log.Error().Err(err).Msg("ERRO FATAL recuperacao secret-pass")
		os.Exit(3)
	}
	envDB.User = string(file_user)
	envDB.Password = string(file_pass)
	envDB.Password="pass123"  //For local testes
	envDB.Db_timeout = 90

	// Load info pod
	// Get IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Error().Err(err).Msg("Error to get the POD IP address !!!")
		os.Exit(3)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				infoPod.IPAddress = ipnet.IP.String()
			}
		}
	}
	infoPod.OSPID = strconv.Itoa(os.Getpid())

	// Get AZ only if localtest is true
	if (noAZ != true) {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Error().Err(err).Msg("ERRO FATAL get Context !!!")
			os.Exit(3)
		}
		client := imds.NewFromConfig(cfg)
		response, err := client.GetInstanceIdentityDocument(context.TODO(), &imds.GetInstanceIdentityDocumentInput{})
		if err != nil {
			log.Error().Err(err).Msg("Unable to retrieve the region from the EC2 instance !!!")
			os.Exit(3)
		}
		infoPod.AvailabilityZone = response.AvailabilityZone	
	} else {
		infoPod.AvailabilityZone = "-"
	}

	// Load Certs
	if (isTLS) {
		caPEM, err = ioutil.ReadFile("/var/pod/cert/caB64.crt")
		if err != nil {
			log.Info().Err(err).Msg("Cert caPEM nao encontrado")
		} else {
			cert.CertPEM = caPEM
		}
	}

	// Load info pod
	infoPod.Database = &envDB
}

func main(){
	log.Debug().Msg("----------------------------------------------------")
	log.Debug().Msg("main")
	log.Debug().Interface("",envDB).Msg("")
	log.Debug().Msg("----------------------------------------------------")
	log.Debug().Interface("",server).Msg("")
	log.Debug().Msg("----------------------------------------------------")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration( server.ReadTimeout ) * time.Second)
	defer cancel()

	// Open Database
	count := 1
	var err error
	for {
		dataBaseHelper, err = postgre.NewDatabaseHelper(ctx, envDB)
		if err != nil {
			if count < 3 {
				log.Error().Err(err).Msg("Erro na abertura do Database")
			} else {
				log.Error().Err(err).Msg("ERRO FATAL na abertura do Database aborting")
				panic(err)
			}
			time.Sleep(3 * time.Second)
			count = count + 1
			continue
		}
		break
	}

	grpcClient, err  := grpc.StartGrpcClient(infoPod.GrpcHost,
											&cert)
	if err != nil {
		log.Error().Err(err).Msg("Erro connect to grpc server")
	}

	restapi	:= restapi.NewRestApi(serverUrlDomain, serverHost ,xApigwId)

	httpAppServerConfig.Server = server
	repoDB = postgre.NewWorkerRepository(dataBaseHelper)

	workerService := service.NewWorkerService(&repoDB, restapi, &grpcClient)
	httpWorkerAdapter := handler.NewHttpWorkerAdapter(workerService)

	httpAppServerConfig.InfoPod = &infoPod
	httpServer := handler.NewHttpAppServer(httpAppServerConfig)

	httpServer.StartHttpAppServer(ctx, httpWorkerAdapter)
}