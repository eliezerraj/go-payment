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
	Msg			any `json:"msg"`
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:		err.Error(),
	}
}

// Middleware v02 - with decoratorDB
func (h *HttpWorkerAdapter) DecoratorDB(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		childLogger.Debug().Msg("-------------- Decorator - MiddleWareHandlerHeader (INICIO) --------------")
	
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers","Content-Type,access-control-allow-origin, access-control-allow-headers")
	
		w.Header().Set("strict-transport-security","max-age=63072000; includeSubdomains; preloa")
		w.Header().Set("content-security-policy","default-src 'none'; img-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'; frame-ancestors 'none'")
		w.Header().Set("x-content-type-option","nosniff")
		w.Header().Set("x-frame-options","DENY")
		w.Header().Set("x-xss-protection","1; mode=block")
		w.Header().Set("referrer-policy","same-origin")
		w.Header().Set("permission-policy","Content-Type,access-control-allow-origin, access-control-allow-headers")

		// If the user was informed then insert it in the session
		if string(r.Header.Get("client-id")) != "" {
			h.workerService.SetSessionVariable(r.Context(),string(r.Header.Get("client-id")))
		} else {
			h.workerService.SetSessionVariable(r.Context(),"NO_INFORMED")
		}

		childLogger.Debug().Msg("-------------- Decorator- MiddleWareHandlerHeader (FIM) ----------------")

		next.ServeHTTP(w, r)
	})
}

func (h *HttpWorkerAdapter) Health(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Health")

	health := true
	json.NewEncoder(rw).Encode(health)
	return
}

func (h *HttpWorkerAdapter) Live(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Live")

	live := true
	json.NewEncoder(rw).Encode(live)
	return
}

func (h *HttpWorkerAdapter) Header(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Header")
	
	span := lib.Span(req.Context(), "handler.header")	
    defer span.End()

	json.NewEncoder(rw).Encode(req.Header)
	return
}

func (h *HttpWorkerAdapter) Auth(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Auth")

	span := lib.Span(req.Context(), "handler.auth")	
    defer span.End()

	authUser := core.AuthUser{}
	authUser.User = h.appServer.AuthUser.User
	authUser.Password = h.appServer.AuthUser.Password

	res, err := h.workerService.Auth(req.Context(), authUser)
	if err != nil {
		switch err {
			case erro.ErrNotFound:
				rw.WriteHeader(404)
				span.RecordError(err)
				json.NewEncoder(rw).Encode(err.Error())
				return
			default:
				rw.WriteHeader(500)
				span.RecordError(err)
				json.NewEncoder(rw).Encode(err.Error())
				return
			}
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) Get(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Get")

	span := lib.Span(req.Context(), "handler.get")	
    defer span.End()

	vars := mux.Vars(req)
	payment := core.Payment{}

	varID, err := strconv.Atoi(vars["id"]) 
    if err != nil {
		apiError := NewAPIError(400, erro.ErrInvalidId)
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
    } 
  
	payment.ID = varID
	res, err := h.workerService.Get(req.Context(), &payment)
	if err != nil {
		var apiError APIError
		switch err {
			case erro.ErrNotFound:
				apiError = NewAPIError(404, err)
			default:
				apiError = NewAPIError(500, err)
		}
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) Pay( rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Pay")

	span := lib.Span(req.Context(), "handler.pay")	
    defer span.End()

	payment := core.Payment{}
	err := json.NewDecoder(req.Body).Decode(&payment)
    if err != nil {
		apiError := NewAPIError(400, erro.ErrUnmarshal)
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
    }

	res, err := h.workerService.Pay(req.Context(), &payment)
	if err != nil {
		var apiError APIError
		switch err {
			case erro.ErrNotFound:
				apiError = NewAPIError(404, err)
			default:
				apiError = NewAPIError(409, err)
		}
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) GetPodInfoGrpc(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("GetPodInfoGrpc")

	span := lib.Span(req.Context(), "handler.getPodInfoGrpc")	
    defer span.End()

	res, err := h.workerService.GetPodInfoGrpc(req.Context())
	if err != nil {
		var apiError APIError
		switch err {
			default:
				apiError = NewAPIError(500, err)
		}
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) CheckPaymentFraudGrpc(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("CheckPaymentFraudGrpc")

	span := lib.Span(req.Context(), "handler.checkPaymentFraudGrpc")	
    defer span.End()

	paymentFraud := core.PaymentFraud{}
	err := json.NewDecoder(req.Body).Decode(&paymentFraud)
    if err != nil {
		apiError := NewAPIError(400, erro.ErrUnmarshal)
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
    }

	res, err := h.workerService.CheckPaymentFraudGrpc(req.Context(), &paymentFraud)
	if err != nil {
		var apiError APIError
		switch err {
			default:
				apiError = NewAPIError(500, err)
		}
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) PayWithCheckFraud( rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("PayWithCheckFraud")

	span := lib.Span(req.Context(), "handler.payWithCheckFraud")	
    defer span.End()

	payment := core.Payment{}
	err := json.NewDecoder(req.Body).Decode(&payment)
    if err != nil {
		apiError := NewAPIError(400, erro.ErrUnmarshal)
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
    }

	res, err := h.workerService.PayWithCheckFraud(req.Context(), &payment)
	if err != nil {
		var apiError APIError
		switch err {
			case erro.ErrNotFound:
				apiError = NewAPIError(404, err)
			default:
				apiError = NewAPIError(409, err)
		}
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
	}

	json.NewEncoder(rw).Encode(res)
	return
}