package configs

import (
	"os"
	"time"
)

const MediaFolder = "media"
const ServerUrl = "http://0.0.0.0:8080"
const FrontEndUrl = "http://localhost:3000"
const ApiVersion = "api/v1"

var PostgresPreferences postgresPreferencesStruct

var RedisPreferences redisPreferencesStruct

const AppleWWDR = "configs/apple_config_files/WWDR.pem"
const AppleCertificate = "configs/apple_config_files/passcertificate.pem"
const AppleKey = "configs/apple_config_files/passkey.pem"

//on urls with this ports frontend need to send data
const HTTPStaffUrl = "0.0.0.0:8084"
const GRPCStaffUrl = "0.0.0.0:8083"
const HTTPSurveyUrl = "0.0.0.0:8086"
const GRPCCafeUrl = "0.0.0.0:8085"
const GRPCCustomerUrl = "0.0.0.0:8082"

var ApplePassword = os.Getenv("ApplePassword")

type timeouts struct {
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	ContextTimeout time.Duration
}

var Timeouts timeouts

func init() {
	PostgresPreferences = postgresPreferencesStruct{
		User:     os.Getenv("PostgresUser"),
		Password: os.Getenv("PostgresPassword"),
		Port:     os.Getenv("PostgresPort"),
	}

	RedisPreferences = redisPreferencesStruct{
		Size:      10,
		Network:   "tcp",
		Address:   os.Getenv("RedisAddress"),
		Password:  os.Getenv("RedisPassword"),
		SecretKey: []byte(os.Getenv("SESSION_KEY")),
	}

	Timeouts = timeouts{
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		ContextTimeout: time.Second * 2,
	}
}
