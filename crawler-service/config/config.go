package config

import (
	"encoding/json"
	"log"
	"os"
)

type configData map[string]string

var jsonData *configData

// LoadConfigFromJSON loads JSON config data
func LoadConfigFromJSON() {
	const file = "env.json"

	envFile, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("\n\n%s does not exist\n%v\n", file, err)
	}

	data := configData{}
	jsonData = &data

	if json.Unmarshal(envFile, &data) != nil {
		log.Fatalf("\nUnable to parse JSON\n%v", err)
	}
}

// GetEnv returns config value of the key.
// It returns an empty string if the key is not present
func GetEnv(key string) string {
	if jsonData == nil {
		LoadConfigFromJSON()
	}

	if value, exists := (*jsonData)[key]; exists {
		return value
	}

	return ""
}
