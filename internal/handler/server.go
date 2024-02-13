package handler

import (
	"time"
	"encoding/json"
	"net/http"
	"strconv"
	"os"
	"os/signal"
	"syscall"
	"context"

	"github.com/gorilla/mux"

	"github.com/go-payment/internal/service"
	"github.com/go-payment/internal/core"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

)

type HttpWorkerAdapter struct {
	workerService 	*service.WorkerService
}

func NewHttpWorkerAdapter(workerService *service.WorkerService) *HttpWorkerAdapter {
	childLogger.Debug().Msg("NewHttpWorkerAdapter")
	return &HttpWorkerAdapter{
		workerService: workerService,
	}
}

type HttpServer struct {
	httpAppServer 	core.HttpAppServer
}

func NewHttpAppServer(httpAppServer core.HttpAppServer) HttpServer {
	childLogger.Debug().Msg("NewHttpAppServer")

	return HttpServer{ httpAppServer: httpAppServer	}
}

func (h HttpServer) StartHttpAppServer(ctx context.Context, httpWorkerAdapter *HttpWorkerAdapter) {
	childLogger.Info().Msg("StartHttpAppServer")
		
	// ---------------------- OTEL ---------------
	var OTEL_EXPORTER_OTLP_ENDPOINT = h.httpAppServer.InfoPod.OtelExportEndpoint
	childLogger.Info().Str("OTEL_EXPORTER_OTLP_ENDPOINT :", OTEL_EXPORTER_OTLP_ENDPOINT).Msg("")

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(),
												otlptracegrpc.WithEndpoint(OTEL_EXPORTER_OTLP_ENDPOINT),
											)
	if err != nil {
		childLogger.Error().Err(err).Msg("ERRO otlptracegrpc")
	}
	idg := xray.NewIDGenerator()

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("go-payment"),
	)

	tp := sdktrace.NewTracerProvider(
									sdktrace.WithSampler(sdktrace.AlwaysSample()),
									sdktrace.WithBatcher(traceExporter),
									sdktrace.WithResource(res),
									sdktrace.WithIDGenerator(idg),
									)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	defer func() { 
		err = tp.Shutdown(ctx)
		if err != nil{
			childLogger.Error().Err(err).Msg("Erro closing OTEL tracer !!!")
		}
	}()
	// ----------------------------------

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		childLogger.Debug().Msg("/")
		json.NewEncoder(rw).Encode(h.httpAppServer)
	})

	myRouter.HandleFunc("/info", func(rw http.ResponseWriter, req *http.Request) {
		childLogger.Debug().Msg("/info")
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(h.httpAppServer)
	})
	
	health := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    health.HandleFunc("/health", httpWorkerAdapter.Health)

	live := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    live.HandleFunc("/live", httpWorkerAdapter.Live)

	header := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    header.HandleFunc("/header", httpWorkerAdapter.Header)
	header.Use(MiddleWareHandlerHeader)

	payPayment := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	payPayment.Handle("/payment/pay", 
						http.HandlerFunc(httpWorkerAdapter.Pay),)
	payPayment.Use(httpWorkerAdapter.DecoratorDB)
	payPayment.Use(otelmux.Middleware("go-payment"))

	getPayment := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	getPayment.Handle("/payment/get/{id}", 
						http.HandlerFunc(httpWorkerAdapter.Get),)
	getPayment.Use(MiddleWareHandlerHeader)
	getPayment.Use(otelmux.Middleware("go-payment"))

	podGrpc := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    podGrpc.Handle("/podGrpc", 
					http.HandlerFunc(httpWorkerAdapter.GetPodGrpc),)
	podGrpc.Use(MiddleWareHandlerHeader)
	podGrpc.Use(otelmux.Middleware("go-payment"))

	srv := http.Server{
		Addr:         ":" +  strconv.Itoa(h.httpAppServer.Server.Port),      	
		Handler:      myRouter,                	          
		ReadTimeout:  time.Duration(h.httpAppServer.Server.ReadTimeout) * time.Second,   
		WriteTimeout: time.Duration(h.httpAppServer.Server.WriteTimeout) * time.Second,  
		IdleTimeout:  time.Duration(h.httpAppServer.Server.IdleTimeout) * time.Second, 
	}

	childLogger.Info().Str("Service Port : ", strconv.Itoa(h.httpAppServer.Server.Port)).Msg("Service Port")

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			childLogger.Error().Err(err).Msg("Cancel http mux server !!!")
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		childLogger.Error().Err(err).Msg("WARNING Dirty Shutdown !!!")
		return
	}

	childLogger.Info().Msg("Stop Done !!!!")
}