package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type configData map[string]string

var configJsonData *configData

// LoadConfigFromJSON loads JSON config data
func LoadConfigFromJSON() {
	const file = "env.json"

	envFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("\n\n%s does not exist\n%v\n consider creating an env.json file", file, err)
	}

	data := configData{}
	configJsonData = &data

	if json.Unmarshal(envFile, &data) != nil {
		log.Fatalf("\nUnable to parse JSON\n%v", err)
	}
}

// GetEnv returns config value of the key.
// It returns an empty string if the key is not present
func GetEnv(key string) string {
	if configJsonData == nil {
		LoadConfigFromJSON()
	}

	if value, exists := (*configJsonData)[key]; exists {
		return value
	}

	return ""
}
