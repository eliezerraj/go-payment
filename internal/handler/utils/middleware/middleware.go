package middleware

import (
	"context"	
	"net/http"
	"github.com/rs/zerolog/log"

	"github.com/go-payment/internal/handler/controller"
)

var childLogger = log.With().Str("handler.utils", "middleware").Logger()

// Middleware v01
func MiddleWareHandlerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		childLogger.Debug().Msg("-------------- MiddleWareHandlerHeader (INICIO)  --------------")
	
		/*if reqHeadersBytes, err := json.Marshal(r.Header); err != nil {
			log.Error().Err(err).Msg("Could not Marshal http headers !!!")
		} else {
			log.Debug().Str("Headers : ", string(reqHeadersBytes) ).Msg("")
		}
		log.Debug().Str("Method : ", r.Method ).Msg("")
		log.Debug().Str("URL : ", r.URL.Path ).Msg("")*/
		//log.Println(r.Header.Get("Host"))
		//log.Println(r.Header.Get("User-Agent"))
		//log.Println(r.Header.Get("X-Forwarded-For"))

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

		ctx := context.WithValue(r.Context(), "tenant_id", string(r.Header.Get("tenant_id")))
		
		childLogger.Debug().Msg("-------------- MiddleWareHandlerHeader (FIM) ----------------")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func MiddleWareErrorHandler(h apiFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		 if err := h(rw, req); err != nil {
			if e, ok := err.(controller.APIError); ok{
				controller.WriteJSON(rw, e.StatusCode, e)
			}
		 }
	 }
}