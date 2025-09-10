package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	key := apiKeyFromEnv()
	ok := strings.HasPrefix(key, "gram")

	fmt.Printf("API key looks good: %v", ok)
}

func apiKeyFromEnv() string {
	return mustEnv("GRAM_API_KEY")
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		panic(fmt.Errorf("missing env: %s", key))
	}

	return val
}
