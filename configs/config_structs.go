package configs

type postgresPreferencesStruct struct {
	User     string
	Password string
	Port     string
	Host     string
	DBName   string
}

type redisPreferencesStruct struct {
	Size      int
	Network   string
	Address   string
	Password  string
	SecretKey []byte
}
