package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/rs/zerolog/log"

	"github.com/go-payment/internal/core/service"
	"github.com/go-payment/internal/core/model"
	"github.com/go-payment/internal/core/erro"
	go_core_observ "github.com/eliezerraj/go-core/observability"
	"github.com/eliezerraj/go-core/coreJson"
	"github.com/gorilla/mux"
)

var childLogger = log.With().Str("adapter", "api.router").Logger()

var core_json coreJson.CoreJson
var core_apiError coreJson.APIError
var tracerProvider go_core_observ.TracerProvider

type HttpRouters struct {
	workerService 	*service.WorkerService
}

func NewHttpRouters(workerService *service.WorkerService) HttpRouters {
	return HttpRouters{
		workerService: workerService,
	}
}

// About return a health
func (h *HttpRouters) Health(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Msg("Health")

	health := true
	json.NewEncoder(rw).Encode(health)
}

// About return a live
func (h *HttpRouters) Live(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Msg("Live")

	live := true
	json.NewEncoder(rw).Encode(live)
}

// About show all header received
func (h *HttpRouters) Header(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Msg("Header")
	
	json.NewEncoder(rw).Encode(req.Header)
}

// About add payment
func (h *HttpRouters) AddPayment(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Msg("AddPayment")

	span := tracerProvider.Span(req.Context(), "adapter.api.AddPayment")
	defer span.End()

	payment := model.Payment{}
	err := json.NewDecoder(req.Body).Decode(&payment)
    if err != nil {
		core_apiError = core_apiError.NewAPIError(err, http.StatusBadRequest)
		return &core_apiError
    }
	defer req.Body.Close()

	res, err := h.workerService.AddPayment(req.Context(), &payment)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

// About get payment
func (h *HttpRouters) GetPayment(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Msg("GetPayment")

	// Trace
	span := tracerProvider.Span(req.Context(), "adapter.api.GetPayment")
	defer span.End()

	// Parameter
	vars := mux.Vars(req)
	varID, err := strconv.Atoi(vars["id"]) 
    if err != nil { 
		core_apiError = core_apiError.NewAPIError(err, http.StatusBadRequest)
		return  &core_apiError
    } 

	payment := model.Payment{}
	payment.ID = varID

	// GetPayment service
	res, err := h.workerService.GetPayment(req.Context(), &payment)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

// About get information from a grpc server (pod information)
func (h *HttpRouters) GetInfoPodGrpc(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Msg("GetInfoPodGrpc")

	// Trace
	span := tracerProvider.Span(req.Context(), "adapter.api.GetInfoPodGrpc")
	defer span.End()

	// GetInfoPodGrpc service
	res, err := h.workerService.GetInfoPodGrpc(req.Context())
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

// About check the score from payment´s features
func (h *HttpRouters) CheckFeaturePaymentFraudGrpc(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Msg("CheckFeaturePaymentFraudGrpc")

	// Trace
	span := tracerProvider.Span(req.Context(), "adapter.api.CheckFeaturePaymentFraudGrpc")
	defer span.End()

	// prepare parameter
	paymentFraud := model.PaymentFraud{}
	err := json.NewDecoder(req.Body).Decode(&paymentFraud)
    if err != nil {
		core_apiError = core_apiError.NewAPIError(err, http.StatusBadRequest)
		return  &core_apiError
    }

	// CheckFeaturePaymentFraudGrpc service
	res, err := h.workerService.CheckFeaturePaymentFraudGrpc(req.Context(), &paymentFraud)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

// About add a payment with the fraud score
func (h *HttpRouters) AddPaymentWithCheckFraud(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Msg("AddPaymentWithCheckFraud")

	// Trace
	span := tracerProvider.Span(req.Context(), "adapter.api.AddPaymentWithCheckFraud")
	defer span.End()

	// prepare parameter
	payment := model.Payment{}
	err := json.NewDecoder(req.Body).Decode(&payment)
    if err != nil {
		core_apiError = core_apiError.NewAPIError(err, http.StatusBadRequest)
		return  &core_apiError
    }

	// AddPaymentWithCheckFraud service
	res, err := h.workerService.AddPaymentWithCheckFraud(req.Context(), &payment)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}