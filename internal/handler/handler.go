package handler

import (	
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/gorilla/mux"

	"github.com/go-payment/internal/core"
	"github.com/go-payment/internal/erro"
	"go.opentelemetry.io/otel"
)

var childLogger = log.With().Str("handler", "handler").Logger()

// Middleware v01
func MiddleWareHandlerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		childLogger.Debug().Msg("-------------- MiddleWareHandlerHeader (INICIO)  --------------")
	
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
		
		childLogger.Debug().Msg("-------------- MiddleWareHandlerHeader (FIM) ----------------")

		next.ServeHTTP(w, r)
	})
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
	
	json.NewEncoder(rw).Encode(req.Header)
	return
}

func (h *HttpWorkerAdapter) Get(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Get")

	ctx, hdlspan := otel.Tracer("go-payment").Start(req.Context(),"handler.Get")
	defer hdlspan.End()

	vars := mux.Vars(req)
	payment := core.Payment{}

	varID, err := strconv.Atoi(vars["id"]) 
    if err != nil { 
		rw.WriteHeader(500)
		json.NewEncoder(rw).Encode(erro.ErrInvalidId.Error())
		return
    } 
  
	payment.ID = varID
	res, err := h.workerService.Get(ctx, payment)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			rw.WriteHeader(404)
			json.NewEncoder(rw).Encode(err.Error())
			return
		default:
			rw.WriteHeader(500)
			json.NewEncoder(rw).Encode(err.Error())
			return
		}
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) Pay( rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Pay")

	ctx, hdlspan := otel.Tracer("go-payment").Start(req.Context(),"handler.Pay")
	defer hdlspan.End()

	payment := core.Payment{}
	err := json.NewDecoder(req.Body).Decode(&payment)
    if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(erro.ErrUnmarshal.Error())
        return
    }

	res, err := h.workerService.Pay(ctx, payment)
	if err != nil {
		switch err {
			case erro.ErrNotFound:
				rw.WriteHeader(404)
				json.NewEncoder(rw).Encode(err.Error())
				return
			default:
				rw.WriteHeader(409)
				json.NewEncoder(rw).Encode(err.Error())
				return
		}
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) GetPodGrpc(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("GetPodGrpc")

	ctx, hdlspan := otel.Tracer("go-payment").Start(req.Context(),"handler.GetPodGrpc")
	defer hdlspan.End()

	res, err := h.workerService.GetInfoPodGrpc(ctx)
	if err != nil {
		switch err {
		default:
			rw.WriteHeader(500)
			json.NewEncoder(rw).Encode(err.Error())
			return
		}
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) CheckPaymentFraudGrpc(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("CheckPaymentFraudGrpc")

	ctx, hdlspan := otel.Tracer("go-payment").Start(req.Context(),"handler.CheckPaymentFraudGrpc")
	defer hdlspan.End()

	paymentFraud := core.PaymentFraud{}
	err := json.NewDecoder(req.Body).Decode(&paymentFraud)
    if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(erro.ErrUnmarshal.Error())
        return
    }

	res, err := h.workerService.CheckPaymentFraudGrpc(ctx, &paymentFraud)
	if err != nil {
		switch err {
		default:
			rw.WriteHeader(500)
			json.NewEncoder(rw).Encode(err.Error())
			return
		}
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) PayFraudFeature( rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("PayFraudFeature")

	ctx, hdlspan := otel.Tracer("go-payment").Start(req.Context(),"handler.PayFraudFeature")
	defer hdlspan.End()

	payment := core.Payment{}
	err := json.NewDecoder(req.Body).Decode(&payment)
    if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(erro.ErrUnmarshal.Error())
        return
    }

	res, err := h.workerService.PayFraudFeature(ctx, payment)
	if err != nil {
		switch err {
			case erro.ErrNotFound:
				rw.WriteHeader(404)
				json.NewEncoder(rw).Encode(err.Error())
				return
			default:
				rw.WriteHeader(409)
				json.NewEncoder(rw).Encode(err.Error())
				return
		}
	}

	json.NewEncoder(rw).Encode(res)
	return
}