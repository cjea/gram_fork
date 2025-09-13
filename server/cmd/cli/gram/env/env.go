package env

import (
	"fmt"
	"os"
	"strings"
)

const (
	// VarNameProducerKey is the environment variable name which points to the user's
	// API key.
	VarNameProducerKey = "GRAM_API_KEY"

	// VarNameProjectSlug is the environment variable name which points to the user's
	// intended project.
	VarNameProjectSlug = "GRAM_PROJECT_SLUG"
)

func MustApiKey() string {
	return validateApiKey(Must(VarNameProducerKey))
}

func MustProjectSlug() string {
	return Must(VarNameProjectSlug)
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
