package restapi

import(
	"errors"
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"context"

	"github.com/rs/zerolog/log"
	"github.com/go-payment/internal/erro"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var childLogger = log.With().Str("adapter/restapi", "restapi").Logger()

type RestApiSConfig struct {
	ServerUrlDomain			string
	XApigwId				string
	ServerHost				string
}

func NewRestApi(serverUrlDomain string, serverHost string ,xApigwId string) (*RestApiSConfig){
	childLogger.Debug().Msg("*** NewRestApi")
	return &RestApiSConfig {
		ServerUrlDomain: 	serverUrlDomain,
		XApigwId: 			xApigwId,
		ServerHost:			serverHost,
	}
}

func (r *RestApiSConfig) GetData(ctx context.Context, serverUrlDomain string, serverHost string, xApigwId string, path string, id string) (interface{}, error) {
	childLogger.Debug().Msg("GetData")

	domain := serverUrlDomain + path +"/" + id

	data_interface, err := makeGet(ctx, domain, serverHost,xApigwId ,id)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Request")
		return nil, err
	}
    
	return data_interface, nil
}

func (r *RestApiSConfig) PostData(ctx context.Context, serverUrlDomain string, serverHost string, xApigwId string, path string ,data interface{}) (interface{}, error) {
	childLogger.Debug().Msg("PostData")

	domain := serverUrlDomain + path 

	data_interface, err := makePost(ctx, domain, serverHost,xApigwId ,data)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Request")
		return nil, err
	}
    
	return data_interface, nil
}

func makeGet(ctx context.Context, url string, serverHost string, xApigwId string, id interface{}) (interface{}, error) {
	childLogger.Debug().Msg("makeGet")

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
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
