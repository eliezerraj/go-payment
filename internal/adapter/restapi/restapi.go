package restapi

import(
	"errors"
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"context"
	"crypto/x509"
	"crypto/tls"
	"encoding/base64"

	"github.com/rs/zerolog/log"
	"github.com/go-payment/internal/erro"
	"github.com/go-payment/internal/core"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var childLogger = log.With().Str("adapter/restapi", "restapi").Logger()
//----------------------------------------
type RestApiService struct {
	restApiConfig *core.AppServer
}

func NewRestApiService(restApiConfig *core.AppServer) *RestApiService{
	childLogger.Debug().Msg("*** NewRestApiService")

	return &RestApiService {
		restApiConfig: restApiConfig,
	}
}
// ---------
func loadClientCertsTLS(cert *core.Cert) (*tls.Config, error){
	childLogger.Debug().Msg("loadClientCertsTLS")

	caPEM_Raw, err := base64.StdEncoding.DecodeString(string(cert.CaAccountPEM))
	if err != nil {
		childLogger.Error().Err(err).Msg("Erro caPEM_Raw !!!")
		return nil, err
	}

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(caPEM_Raw)

	clientTLSConf := &tls.Config{
		RootCAs: certpool,
	}

	return clientTLSConf ,nil
}

func (r *RestApiService) CallRestApi(	ctx context.Context, 
										method	string,
										path string, 
										x_apigw_api_id *string, 
										body interface{}) (interface{}, error) {
	childLogger.Debug().Msg("CallRestApi")

	childLogger.Debug().Str("method : ", method).Msg("")
	childLogger.Debug().Str("path : ", path).Msg("")
	childLogger.Debug().Interface("x_apigw_api_id : ", x_apigw_api_id).Msg("")
	childLogger.Debug().Interface("body : ", body).Msg("")


	transportHttp := &http.Transport{}
	// -------------- Load Certs -------------------------
	if string(r.restApiConfig.RestEndpoint.CaCert.CaAccountPEM) != "" {
		transportHttpConfig, err := loadClientCertsTLS(r.restApiConfig.RestEndpoint.CaCert)
		if err != nil {
			childLogger.Error().Err(err).Msg("Erro loadClientCertsTLS")
			return nil, err
		}
		transportHttp.TLSClientConfig = transportHttpConfig
	} 
	// -------------- Load Certs -------------------------

	client := http.Client{
		Transport: otelhttp.NewTransport(transportHttp),
		Timeout: time.Second * 29,
	}

	payload := new(bytes.Buffer) 
	if body != nil{
		json.NewEncoder(payload).Encode(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, path, payload)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Request")
		return false, errors.New(err.Error())
	}

	req.Header.Add("Content-Type", "application/json;charset=UTF-8");
	if (x_apigw_api_id != nil){
		req.Header.Add("x-apigw-api-id", *x_apigw_api_id)
	}
	req.Host = r.restApiConfig.RestEndpoint.ServerHost;

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		childLogger.Error().Err(err).Msg("error Do Request")
		return false, errors.New(err.Error())
	}

	childLogger.Debug().Int("StatusCode :", resp.StatusCode).Msg("")
	switch (resp.StatusCode) {
		case 401:
			return false, erro.ErrHTTPForbiden
		case 403:
			return false, erro.ErrHTTPForbiden
		case 200:
		case 400:
			return false, erro.ErrNotFound
		case 404:
			return false, erro.ErrNotFound
		default:
			return false, erro.ErrHTTPForbiden
	}

	result := body
	err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
		childLogger.Error().Err(err).Msg("error no ErrUnmarshal")
		return false, errors.New(err.Error())
    }

	return result, nil
}