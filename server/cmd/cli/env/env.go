package env

import (
	"fmt"
	"os"
	"strings"
)

// EnvVarGramAPIKey is the environment variable name which points to the user's
// API key.
const EnvVarGramAPIKey = "GRAM_API_KEY" // #nosec G101

func ReadApiKey() string {
	return validateApiKey(Must(EnvVarGramAPIKey))
}

func validateApiKey(key string) string {
	ok := strings.HasPrefix(key, "gram")

	if ok {
		return key
	} else {
		panic(fmt.Errorf("key is malformed"))
	}
}

func Must(key string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		panic(fmt.Errorf("missing env: %s", key))
	}

	return val
}

func Fallback(key string, fallback string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return fallback
	} else {
		return val
	}
}
