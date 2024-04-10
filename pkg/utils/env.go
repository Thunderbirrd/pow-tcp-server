package utils

import (
	"os"
	"strconv"
	"time"
)

func EnvToInt(value *int, key string, defaultValue int) {
	*value = defaultValue

	envValue, exists := os.LookupEnv(key)
	if !exists || envValue == "" {
		return
	}

	if res, err := strconv.Atoi(envValue); err == nil {
		*value = res
	}
}

func EnvToStr(value *string, key string, defaultValue string) {
	*value = getEnv(key, defaultValue)
}

func EnvToDuration(value *time.Duration, key string, defaultValue time.Duration) {
	if strVal, exists := os.LookupEnv(key); exists && strVal != "" {
		var err error
		*value, err = time.ParseDuration(strVal)
		if err != nil {
			*value = defaultValue
		}
	} else {
		*value = defaultValue
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}

	return defaultValue
}
