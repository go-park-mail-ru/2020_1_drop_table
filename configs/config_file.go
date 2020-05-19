package configs

import (
	"os"
	"time"
)

const MediaFolder = "media"
const ApiVersion = "api/v1"

var FrontEndUrl = os.Getenv("FrontEndUrl")
var ServerUrl = os.Getenv("ServerUrl")

var PostgresPreferences postgresPreferencesStruct

var RedisPreferences redisPreferencesStruct

const AppleWWDR = "configs/apple_config_files/WWDR.pem"
const AppleCertificate = "configs/apple_config_files/passcertificate.pem"
const AppleKey = "configs/apple_config_files/passkey.pem"

//on urls with this ports frontend need to send data
const MainService = "0.0.0.0:8080"
const HTTPStaffUrl = "0.0.0.0:8084"
const GRPCStaffUrl = "staff:8083"
const HTTPSurveyUrl = "0.0.0.0:8086"
const GRPCCafeUrl = "main_service:8085"
const GRPCCustomerUrl = "0.0.0.0:8082"

var ApplePassword = os.Getenv("ApplePassword")

var GoogleMapAPIKey = os.Getenv("GoogleMapAPIKey")

type timeouts struct {
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	ContextTimeout time.Duration
}

var Timeouts timeouts

type sessionName string

const (
	SessionStaffID sessionName = "session"
)

func init() {
	PostgresPreferences = postgresPreferencesStruct{
		User:     os.Getenv("PostgresUser"),
		Password: os.Getenv("PostgresPassword"),
		Port:     os.Getenv("PostgresPort"),
		Host:     os.Getenv("PostgresHost"),
		DBName:   os.Getenv("PostgresDBName"),
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
