package main

import(
	"time"
	"context"
	
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/go-payment/internal/infra/configuration"
	"github.com/go-payment/internal/core/model"
	"github.com/go-payment/internal/core/service"
	"github.com/go-payment/internal/infra/server"
	"github.com/go-payment/internal/adapter/api"
	"github.com/go-payment/internal/adapter/grpc/client"
	"github.com/go-payment/internal/adapter/database"
	go_core_pg "github.com/eliezerraj/go-core/database/pg"  
)

var(
	logLevel = 	zerolog.InfoLevel // zerolog.InfoLevel zerolog.DebugLevel
	appServer	model.AppServer
	databaseConfig go_core_pg.DatabaseConfig
	databasePGServer go_core_pg.DatabasePGServer
	childLogger = log.With().Str("component","go-payment").Str("package", "main").Logger()
)

// About initialize the enviroment var
func init(){
	childLogger.Info().Str("func","init").Send()

	zerolog.SetGlobalLevel(logLevel)

	infoPod, server := configuration.GetInfoPod()
	configOTEL 		:= configuration.GetOtelEnv()
	databaseConfig 	:= configuration.GetDatabaseEnv() 
	apiService 	:= configuration.GetEndpointEnv() 

	appServer.InfoPod = &infoPod
	appServer.Server = &server
	appServer.ConfigOTEL = &configOTEL
	appServer.DatabaseConfig = &databaseConfig
	appServer.ApiService = apiService
}

// About main
func main (){
	childLogger.Info().Str("func","main").Interface("appServer",appServer).Send()

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( appServer.Server.ReadTimeout ) * time.Second)
	defer cancel()

	// Open Database
	count := 1
	var err error
	for {
		databasePGServer, err = databasePGServer.NewDatabasePGServer(ctx, *appServer.DatabaseConfig)
		if err != nil {
			if count < 3 {
				childLogger.Error().Err(err).Msg("error open database... trying again !!")
			} else {
				childLogger.Error().Err(err).Msg("fatal error open Database aborting")
				panic(err)
			}
			time.Sleep(3 * time.Second) //backoff
			count = count + 1
			continue
		}
		break
	}

	// Open GRPC channel
	grpcClient, err  := client.StartGrpcClient(appServer.ApiService[2].Url)
	if err != nil {
		childLogger.Error().Err(err).Msg("Erro connect to grpc server")
	}
	// test connection
	err = grpcClient.TestConnection(ctx)
	if err != nil {
		childLogger.Error().Err(err).Msg("Erro connect to grpc server")
	} else {
		childLogger.Info().Msg("gprc channel openned sucessfull")
	}

	// wire
	database := database.NewWorkerRepository(&databasePGServer)
	workerService := service.NewWorkerService(database, appServer.ApiService, &grpcClient)
	httpRouters := api.NewHttpRouters(workerService)
	httpServer := server.NewHttpAppServer(appServer.Server)

	// start server
	httpServer.StartHttpAppServer(ctx, &httpRouters, &appServer)
}