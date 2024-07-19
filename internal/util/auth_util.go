package util

import(
	"os"

	"github.com/joho/godotenv"
	"github.com/go-payment/internal/core"
)

func GetAuthEnv() (core.AuthUser) {
	childLogger.Debug().Msg("GetAuthEnv")

	err := godotenv.Load(".env")
	if err != nil {
		childLogger.Info().Err(err).Msg("No .env File !!!!")
	}

	authUser := core.AuthUser{}

	if os.Getenv("USER_AUTH") !=  "" {	
		authUser.User = os.Getenv("USER_AUTH")
	}
	if os.Getenv("PASSWORD_AUTH") !=  "" {	
		authUser.Password = os.Getenv("PASSWORD_AUTH")
	}

	return authUser
}
