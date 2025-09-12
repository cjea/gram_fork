package env

import (
	"fmt"
	"os"
	"strings"
)

const (
	// EnvVarGramAPIKey is the environment variable name which points to the user's
	// API key.
	EnvVarGramAPIKey = "GRAM_API_KEY" // #nosec G101

	// EnvVarProjectSlug is the environment variable name which points to the user's
	// intended project.
	EnvVarProjectSlug = "GRAM_PROJECT_SLUG" // #nosec G101
)

func ReadApiKey() string {
	return validateApiKey(Must(EnvVarGramAPIKey))
}

func ReadProjectSlug() string {
	return Must(EnvVarProjectSlug)
}

const apiKeyPrefix = "gram"

func validateApiKey(key string) string {
	ok := strings.HasPrefix(key, apiKeyPrefix)

	if ok {
		return key
	} else {
		panic(fmt.Errorf("key is malformed: expected prefix '%s'", apiKeyPrefix))
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
