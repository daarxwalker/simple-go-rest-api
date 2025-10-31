package tests

import (
	"testing"

	"gocourse/pkg/env"
)

func SetTestEnv(t testing.TB) {
	t.Setenv(env.AppEnv, "test")
}
