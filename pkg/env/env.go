package env

import "os"

const (
	AppEnv  = "APP_ENV"
	AppPort = "APP_PORT"
)

func Env() string {
	return os.Getenv(AppEnv)
}

func Port() string {
	return os.Getenv(AppPort)
}

func Development() bool {
	return Env() == "development"
}

func Production() bool {
	return Env() == "production"
}

func Test() bool {
	return Env() == "test"
}
