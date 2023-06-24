package envs

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

type Env struct {
	jWTToken        string
	jWTRefreshToken string
	dbHost          string
	dbPassword      string
	dbUser          string
	dbPort          string
	dbName          string
	redisURL        string
	webServerPort   string
}

func newEnv() *Env {
	InitEnv()
	return &Env{
		jWTToken:        os.Getenv("JWT_SECRET"),
		jWTRefreshToken: os.Getenv("JWT_REFRESH_SECRET"),
		dbHost:          os.Getenv("DB_HOST"),
		dbPassword:      os.Getenv("DB_PASSWORD"),
		dbUser:          os.Getenv("DB_USERNAME"),
		dbPort:          os.Getenv("DB_PORT"),
		dbName:          os.Getenv("DB_DATABASE"),
		redisURL:        os.Getenv("UPSTASH_REDIS_URL"),
		webServerPort:   os.Getenv("WEB_SERVER_PORT"),
	}
}

func (env *Env) GetJWTToken() string {
	return env.jWTToken
}

func (env *Env) GetJWTRefreshToken() string {
	return env.jWTRefreshToken
}

func (env *Env) GetDbHost() string {
	return env.dbHost
}

func (env *Env) GetDbPassword() string {
	return env.dbPassword
}

func (env *Env) GetDbUser() string {
	return env.dbUser
}

func (env *Env) GetDbPort() string {
	return env.dbPort
}

func (env *Env) GetDbName() string {
	return env.dbName
}

func (env *Env) GetRedisURL() string {
	return env.redisURL
}

func (env *Env) GetWebServerPort() string {
	return env.webServerPort
}

var instance *Env

func GetInstance() *Env {
	if instance == nil {
		instance = newEnv()
	}
	return instance
}

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		execPath, err := os.Executable()
		if err != nil {
			panic(err)
		}
		envPath := filepath.Join(filepath.Dir(execPath), ".env")
		err = godotenv.Load(envPath)
		if err != nil {
			panic(err)
		}
	}
}
