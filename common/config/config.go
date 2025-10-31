package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

func Load() *viper.Viper {
	cfg := viper.New()
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := strings.ToLower(pair[0])
		val := pair[1]
		key = strings.ReplaceAll(key, "_", ".")
		key = strings.ToLower(key)
		cfg.Set(key, val)
	}
	return cfg
}
