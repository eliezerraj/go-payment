package controller

import (
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/gorilla/mux"

	"github.com/go-payment/internal/service"
	"github.com/go-payment/internal/core"
	"github.com/go-payment/internal/erro"
	"github.com/go-payment/internal/lib"
)

var childLogger = log.With().Str("handler", "controller").Logger()

//-------------------------------------------
type HttpWorkerAdapter struct {
	workerService 	*service.WorkerService
	appServer 		*core.AppServer
}

func NewHttpWorkerAdapter(workerService *service.WorkerService,	appServer *core.AppServer) HttpWorkerAdapter {
	childLogger.Debug().Msg("NewHttpWorkerAdapter")
	return HttpWorkerAdapter{
		workerService: workerService,
		appServer: appServer,
	}
}

type APIError struct {
	StatusCode	int  `json:"statusCode"`
	Msg			string `json:"msg"`
}

func (e APIError) Error() string {
	return e.Msg
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:		err.Error(),
	}
}

func WriteJSON(rw http.ResponseWriter, code int, v any) error{
	rw.WriteHeader(code)
	return json.NewEncoder(rw).Encode(v)
}

func (h *HttpWorkerAdapter) Health(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Health")

	health := true
	json.NewEncoder(rw).Encode(health)
}

func (h *HttpWorkerAdapter) Live(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Live")

	live := true
	json.NewEncoder(rw).Encode(live)
}

func (h *HttpWorkerAdapter) Header(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Header")
	
	span := lib.Span(req.Context(), "handler.header")	
    defer span.End()

	json.NewEncoder(rw).Encode(req.Header)
}

func (h *HttpWorkerAdapter) Auth(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Debug().Msg("Auth")

	span := lib.Span(req.Context(), "handler.auth")	
    defer span.End()

	authUser := core.AuthUser{}
	authUser.User = h.appServer.AuthUser.User
	authUser.Password = h.appServer.AuthUser.Password

	res, err := h.workerService.Auth(req.Context(), authUser)
	if err != nil {
		var apiError APIError
		switch err {
			case erro.ErrNotFound:
				apiError = NewAPIError(http.StatusNotFound, err)
			default:
				apiError = NewAPIError(http.StatusInternalServerError, err)
		}
		return apiError
	}

	return WriteJSON(rw, http.StatusOK, res)
}

func (h *HttpWorkerAdapter) Get(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Debug().Msg("Get")

	span := lib.Span(req.Context(), "handler.get")	
    defer span.End()

	vars := mux.Vars(req)
	payment := core.Payment{}

	varID, err := strconv.Atoi(vars["id"]) 
    if err != nil {
		apiError := NewAPIError(http.StatusBadRequest, erro.ErrInvalidId)
		return apiError
    } 
  
	payment.ID = varID
	res, err := h.workerService.Get(req.Context(), &payment)
	if err != nil {
		var apiError APIError
		switch err {
			case erro.ErrNotFound:
				apiError = NewAPIError(http.StatusNotFound, err)
			default:
				apiError = NewAPIError(http.StatusInternalServerError, err)
		}
		return apiError
	}

	return WriteJSON(rw, http.StatusOK, res)
}

func (h *HttpWorkerAdapter) Pay( rw http.ResponseWriter, req *http.Request) error {
	childLogger.Debug().Msg("Pay")

	span := lib.Span(req.Context(), "handler.pay")	
    defer span.End()

	payment := core.Payment{}

	payment.TenantID = req.Context().Value("tenant_id").(string)

	err := json.NewDecoder(req.Body).Decode(&payment)
    if err != nil {
		apiError := NewAPIError(http.StatusBadRequest, erro.ErrUnmarshal)
		return apiError
    }

	res, err := h.workerService.Pay(req.Context(), &payment)
	if err != nil {
		var apiError APIError
		switch err {
			case erro.ErrNotFound:
				apiError = NewAPIError(http.StatusNotFound, err)
			default:
				apiError = NewAPIError(http.StatusInternalServerError, err)
		}
		return apiError
	}

	return WriteJSON(rw, http.StatusOK, res)
}

func (h *HttpWorkerAdapter) GetPodInfoGrpc(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Debug().Msg("GetPodInfoGrpc")

	span := lib.Span(req.Context(), "handler.getPodInfoGrpc")	
    defer span.End()

	res, err := h.workerService.GetPodInfoGrpc(req.Context())
	if err != nil {
		var apiError APIError
		switch err {
			default:
				apiError = NewAPIError(http.StatusInternalServerError, err)
		}
		return apiError
	}

	return WriteJSON(rw, http.StatusOK, res)
}

func (h *HttpWorkerAdapter) CheckPaymentFraudGrpc(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Debug().Msg("CheckPaymentFraudGrpc")

	span := lib.Span(req.Context(), "handler.checkPaymentFraudGrpc")	
    defer span.End()

	paymentFraud := core.PaymentFraud{}
	err := json.NewDecoder(req.Body).Decode(&paymentFraud)
    if err != nil {
		apiError := NewAPIError(http.StatusBadRequest, erro.ErrUnmarshal)
		return apiError
    }

	res, err := h.workerService.CheckPaymentFraudGrpc(req.Context(), &paymentFraud)
	if err != nil {
		var apiError APIError
		switch err {
			default:
				apiError = NewAPIError(http.StatusInternalServerError, err)
		}
		return apiError
	}

	return WriteJSON(rw, http.StatusOK, res)
}

func (h *HttpWorkerAdapter) PayWithCheckFraud( rw http.ResponseWriter, req *http.Request) error {
	childLogger.Debug().Msg("PayWithCheckFraud")

	span := lib.Span(req.Context(), "handler.payWithCheckFraud")	
    defer span.End()

	payment := core.Payment{}
	err := json.NewDecoder(req.Body).Decode(&payment)
    if err != nil {
		apiError := NewAPIError(http.StatusBadRequest, erro.ErrUnmarshal)
		return apiError
    }

	res, err := h.workerService.PayWithCheckFraud(req.Context(), &payment)
	if err != nil {
		var apiError APIError
		switch err {
			case erro.ErrNotFound:
				apiError = NewAPIError(http.StatusNotFound, err)
			default:
				apiError = NewAPIError(http.StatusInternalServerError, err)
		}
		return apiError
	}

	return WriteJSON(rw, http.StatusOK, res)
}