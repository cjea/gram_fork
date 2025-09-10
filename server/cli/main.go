package main

import (
	"fmt"
	"os"
	"strings"
)

var API_KEY string = apiKeyFromEnv()

func main() {
	fmt.Printf("Starting CLI.")
}

func apiKeyFromEnv() string {
	return validateApiKey(mustEnv("GRAM_API_KEY"))
}

func validateApiKey(key string) string {
	ok := strings.HasPrefix(key, "gram")

	if ok {
		return key
	} else {
		panic(fmt.Errorf("key is malformed"))
	}
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		panic(fmt.Errorf("missing env: %s", key))
	}

	return val
}
