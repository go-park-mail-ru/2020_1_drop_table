package configs

import "os"

const MediaFolder = "media"
const ServerUrl = "http://localhost:8080"
const FrontEndUrl = "http://localhost:3000"

var PostgresPreferences postgresPreferencesStruct

var RedisPreferences redisPreferencesStruct

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
