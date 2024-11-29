package webserver

import(
	"time"
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/go-payment/internal/core"
	"github.com/go-payment/internal/util"
	"github.com/go-payment/internal/repository/pg"
	"github.com/go-payment/internal/repository/storage"
	"github.com/go-payment/internal/handler/controller"
	"github.com/go-payment/internal/service"
	"github.com/go-payment/internal/adapter/restapi"
	"github.com/go-payment/internal/handler"
	"github.com/go-payment/internal/adapter/grpc"
)

var(
	logLevel = zerolog.DebugLevel
	appServer	core.AppServer
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)
	
	infoPod, server, restEndpoint := util.GetInfoPod()
	database := util.GetDatabaseEnv()
	configOTEL := util.GetOtelEnv()
	caCert := util.GetCaCertEnv()
	authUser := util.GetAuthEnv()

	appServer.InfoPod = &infoPod
	appServer.AuthUser = &authUser
	appServer.Database = &database
	appServer.Server = &server
	appServer.RestEndpoint = &restEndpoint
	appServer.RestEndpoint.CaCert = &caCert
	appServer.ConfigOTEL = &configOTEL
}

func Server(){
	log.Debug().Msg("----------------------------------------------------")
	log.Debug().Msg("main")
	log.Debug().Msg("----------------------------------------------------")
	log.Debug().Interface("appServer :",appServer).Msg("")
	log.Debug().Msg("----------------------------------------------------")

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( appServer.Server.ReadTimeout ) * time.Second)
	defer cancel()
	
	// Open Database
	count := 1
	var databasePG	pg.DatabasePG
	var err error
	for {
		databasePG, err = pg.NewDatabasePGServer(ctx, appServer.Database)
		if err != nil {
			if count < 3 {
				log.Error().Err(err).Msg("Erro open Database... trying again !!")
			} else {
				log.Error().Err(err).Msg("Fatal erro open Database aborting")
				panic(err)
			}
			time.Sleep(3 * time.Second)
			count = count + 1
			continue
		}
		break
	}

	repoDB := storage.NewWorkerRepository(databasePG)
	grpcClient, err  := grpc.StartGrpcClient(appServer.RestEndpoint.GrpcHost, appServer.RestEndpoint.CaCert)
	if err != nil {
		log.Error().Err(err).Msg("Erro connect to grpc server")
	}

	restApiService	:= restapi.NewRestApiService(&appServer)

	workerService 		:= service.NewWorkerService(&repoDB, &appServer, restApiService, &grpcClient)
	httpWorkerAdapter 	:= controller.NewHttpWorkerAdapter(workerService, &appServer)
	httpServer 			:= handler.NewHttpAppServer(appServer.Server)

	httpServer.StartHttpAppServer(ctx, &httpWorkerAdapter, &appServer)
}