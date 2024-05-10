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

type RestApiSConfig struct {
	ServerUrlDomain			string
	GatewayMlHost			string
	XApigwId				string
	XApigwIdMlHost			string
	ServerHost				string
	Cert					core.Cert
}

func NewRestApi(serverUrlDomain string,
				gatewayMlHost	string,
				serverHost 		string,
				xApigwId 		string,
				xApigwIdMlHost	string,
				cert core.Cert) (*RestApiSConfig){
	childLogger.Debug().Msg("*** NewRestApi")

	return &RestApiSConfig {
							ServerUrlDomain: 	serverUrlDomain,
							GatewayMlHost:		gatewayMlHost,
							XApigwId: 			xApigwId,
							XApigwIdMlHost:		xApigwIdMlHost,
							ServerHost:			serverHost,
							Cert:				cert,
						}
}

func (r *RestApiSConfig) GetData(ctx context.Context, 
								serverUrlDomain string, 
								serverHost string, 
								xApigwId string, 
								path string, 
								id string) (interface{}, error) {
	childLogger.Debug().Msg("GetData")

	domain := serverUrlDomain + path +"/" + id

	data_interface, err := r.makeGet(ctx, domain, serverHost,xApigwId ,id)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Request")
		return nil, err
	}
    
	return data_interface, nil
}

func (r *RestApiSConfig) PostData(	ctx context.Context, 
									serverUrlDomain string, 
									serverHost string, 
									xApigwId string, 
									path string,
									data interface{}) (interface{}, error) {
	childLogger.Debug().Msg("PostData")

	domain := serverUrlDomain + path 

	data_interface, err := makePost(ctx, domain, serverHost, xApigwId, data)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Request")
		return nil, err
	}
    
	return data_interface, nil
}

func loadClientCertsTLS(cert *core.Cert) (*tls.Config, error){
	childLogger.Debug().Msg("loadClientCertsTLS")

	caPEM_Raw, err := base64.StdEncoding.DecodeString(string(cert.CertAccountPEM))
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

func (r *RestApiSConfig) makeGet(ctx context.Context, 
								url string, 
								serverHost string, 
								xApigwId string, 
								id interface{}) (interface{}, error) {
	childLogger.Debug().Msg("makeGet")

	transportHttp := &http.Transport{}
	// -------------- Load Certs -------------------------
	if string(r.Cert.CertAccountPEM) != "" {
		transportHttpConfig, err := loadClientCertsTLS(&r.Cert)
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

	childLogger.Debug().Str("url : ", url).Msg("")
	childLogger.Debug().Str("serverHost : ", serverHost).Msg("")
	childLogger.Debug().Str("xApigwId : ", xApigwId).Msg("")

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Request")
		return false, errors.New(err.Error())
	}

	req.Header.Add("Content-Type", "application/json;charset=UTF-8");
	req.Header.Add("x-apigw-api-id", xApigwId);
	req.Host = serverHost;

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
			return false, erro.ErrServer
	}

	var result interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
		childLogger.Error().Err(err).Msg("error no ErrUnmarshal")
		return false, errors.New(err.Error())
    }

	return result, nil
}

func makePost(ctx context.Context, url string, serverHost string, xApigwId string, data interface{}) (interface{}, error) {
	childLogger.Debug().Msg("makePost")

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout: time.Second * 29,
	}

	childLogger.Debug().Str("url : ", url).Msg("")
	childLogger.Debug().Str("serverHost : ", serverHost).Msg("")
	childLogger.Debug().Str("xApigwId : ", xApigwId).Msg("")

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(data)

	req, err := http.NewRequestWithContext(ctx ,"POST", url, payload)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Request")
		return false, errors.New(err.Error())
	}

	req.Header.Add("Content-Type", "application/json;charset=UTF-8");
	req.Header.Add("x-apigw-api-id", xApigwId);
	req.Host = serverHost;

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

	result := data
	err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
		childLogger.Error().Err(err).Msg("error no ErrUnmarshal")
		return false, errors.New(err.Error())
    }

	return result, nil
}
