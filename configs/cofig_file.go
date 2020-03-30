package configs

import "os"

const MediaFolder = "media"
const ServerUrl = "http://localhost:8080"
const FrontEndUrl = "http://localhost:3000"
const ApiVersion = "api/v1"

var PostgresPreferences postgresPreferencesStruct

var RedisPreferences redisPreferencesStruct

const AppleWWDR = "configs/apple_config_files/WWDR.pem"
const AppleCertificate = "configs/apple_config_files/passcertificate.pem"
const AppleKey = "configs/apple_config_files/passkey.pem"

var ApplePassword = os.Getenv("ApplePassword")

func init() {
	PostgresPreferences = postgresPreferencesStruct{
		User:     "postgres",
		Password: "",
		Port:     "5431",
	}

	RedisPreferences = redisPreferencesStruct{
		Size:      10,
		Network:   "tcp",
		Address:   ":6379",
		Password:  "",
		SecretKey: []byte(os.Getenv("SESSION_KEY")),
	}
}
