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
//-------------------------------------------
type HttpWorkerAdapter struct {
	workerService 	*service.WorkerService
}

func NewHttpWorkerAdapter(workerService *service.WorkerService) HttpWorkerAdapter {
	childLogger.Debug().Msg("NewHttpWorkerAdapter")
	return HttpWorkerAdapter{
		workerService: workerService,
	}
}
//-------------------------------------------
type HttpServer struct {
	httpServer	*core.Server
}

func NewHttpAppServer(httpServer *core.Server) HttpServer {
	childLogger.Debug().Msg("NewHttpAppServer")

	return HttpServer{httpServer: httpServer }
}
//-------------------------------------------
func (h HttpServer) StartHttpAppServer(	ctx context.Context, 
										httpWorkerAdapter *HttpWorkerAdapter,
										appServer *core.AppServer) {
	childLogger.Info().Msg("StartHttpAppServer")
		
	// ---------------------- OTEL ---------------
	childLogger.Info().Str("OTEL_EXPORTER_OTLP_ENDPOINT :", appServer.ConfigOTEL.OtelExportEndpoint).Msg("")

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(),
												otlptracegrpc.WithEndpoint(appServer.ConfigOTEL.OtelExportEndpoint),
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
	myRouter.Use(MiddleWareHandlerHeader)

	myRouter.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		childLogger.Debug().Msg("/")
		json.NewEncoder(rw).Encode(appServer)
	})

	myRouter.HandleFunc("/info", func(rw http.ResponseWriter, req *http.Request) {
		childLogger.Debug().Msg("/info")
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(appServer)
	})
	
	health := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    health.HandleFunc("/health", httpWorkerAdapter.Health)

	live := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    live.HandleFunc("/live", httpWorkerAdapter.Live)

	header := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    header.HandleFunc("/header", httpWorkerAdapter.Header)

	payPayment := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	payPayment.Handle("/payment/pay", 
						http.HandlerFunc(httpWorkerAdapter.Pay),)
	payPayment.Use(httpWorkerAdapter.DecoratorDB)
	payPayment.Use(otelmux.Middleware("go-payment"))

	getPayment := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	getPayment.Handle("/payment/get/{id}", 
						http.HandlerFunc(httpWorkerAdapter.Get),)
	getPayment.Use(otelmux.Middleware("go-payment"))

	podGrpc := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    podGrpc.Handle("/getPodInfoGrpc", 
					http.HandlerFunc(httpWorkerAdapter.GetPodInfoGrpc),)
	podGrpc.Use(otelmux.Middleware("go-payment"))

	paymentFraudGrpc := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
    paymentFraudGrpc.Handle("/checkPaymentFraudGrpc", 
					http.HandlerFunc(httpWorkerAdapter.CheckPaymentFraudGrpc),)
	paymentFraudGrpc.Use(otelmux.Middleware("go-payment"))

	paymentFraudFeature := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
    paymentFraudFeature.Handle("/payment/payWithCheckFraud", 
					http.HandlerFunc(httpWorkerAdapter.PayWithCheckFraud),)
	paymentFraudFeature.Use(otelmux.Middleware("go-payment"))

	srv := http.Server{
		Addr:         ":" +  strconv.Itoa(h.httpServer.Port),      	
		Handler:      myRouter,                	          
		ReadTimeout:  time.Duration(h.httpServer.ReadTimeout) * time.Second,   
		WriteTimeout: time.Duration(h.httpServer.WriteTimeout) * time.Second,  
		IdleTimeout:  time.Duration(h.httpServer.IdleTimeout) * time.Second, 
	}

	childLogger.Info().Str("Service Port : ", strconv.Itoa(h.httpServer.Port)).Msg("Service Port")

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