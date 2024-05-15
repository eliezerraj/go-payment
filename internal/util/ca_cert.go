package util

import(
	"os"
	"io/ioutil"

	"github.com/joho/godotenv"
	"github.com/go-payment/internal/core"
)

func GetCaCertEnv() core.Cert {
	childLogger.Debug().Msg("GetCaCertEnv")

	err := godotenv.Load(".env")
	if err != nil {
		childLogger.Info().Err(err).Msg("No .env File !!!!")
	}

	var cert						core.Cert

	if os.Getenv("TLS_FRAUD") ==  "true" {	
		childLogger.Debug().Msg("*** Loading ca_fraud_B64.crt ***")

		cert.CaFraudPEM, err = ioutil.ReadFile("/var/pod/cert/ca_fraud_B64.crt")
		if err != nil {
			childLogger.Info().Err(err).Msg("ca_fraud_B641.crt not found")
		}
	}

	if os.Getenv("TLS_ACCOUNT") ==  "true" {	
		childLogger.Debug().Msg("*** Loading ca_account_B64.crt ***")

		cert.CaAccountPEM, err = ioutil.ReadFile("/var/pod/cert/ca_account_B64.crt")
		if err != nil {
			childLogger.Info().Err(err).Msg("ca_account_B64.crt not found")
		}
	}

	return cert
}