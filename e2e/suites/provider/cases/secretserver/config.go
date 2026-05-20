/*
Copyright © The ESO Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package secretserver

import (
	"fmt"
	"os"
	"strconv"
)

type config struct {
	username                string
	password                string
	serverURL               string
	siteID                  int
	folderID                int
	secretTemplateID        int
	dataFieldID             int
	requiredPasswordFieldID int
	requiredPasswordValue   string
}

func loadConfigFromEnv() (*config, error) {
	var cfg config
	var err error

	// Required settings
	cfg.username, err = getEnv("SECRETSERVER_USERNAME")
	if err != nil {
		return nil, err
	}
	cfg.password, err = getEnv("SECRETSERVER_PASSWORD")
	if err != nil {
		return nil, err
	}
	cfg.serverURL, err = getEnv("SECRETSERVER_URL")
	if err != nil {
		return nil, err
	}

	cfg.siteID, err = getOptionalIntEnv("SECRETSERVER_SITE_ID", 1)
	if err != nil {
		return nil, err
	}
	cfg.folderID, err = getOptionalIntEnv("SECRETSERVER_FOLDER_ID", 10)
	if err != nil {
		return nil, err
	}
	cfg.secretTemplateID, err = getOptionalIntEnv("SECRETSERVER_TEMPLATE_ID", 6051)
	if err != nil {
		return nil, err
	}
	cfg.dataFieldID, err = getOptionalIntEnv("SECRETSERVER_DATA_FIELD_ID", 329)
	if err != nil {
		return nil, err
	}
	cfg.requiredPasswordFieldID, err = getOptionalIntEnv("SECRETSERVER_REQUIRED_PASSWORD_FIELD_ID", 0)
	if err != nil {
		return nil, err
	}
	cfg.requiredPasswordValue = "external-secrets-e2e"
	if value := os.Getenv("SECRETSERVER_REQUIRED_PASSWORD_VALUE"); value != "" {
		cfg.requiredPasswordValue = value
	}

	return &cfg, nil
}

func getEnv(name string) (string, error) {
	value, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf("environment variable %q is not set", name)
	}
	return value, nil
}

func getOptionalIntEnv(name string, defaultValue int) (int, error) {
	value, ok := os.LookupEnv(name)
	if !ok || value == "" {
		return defaultValue, nil
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %q must be an integer: %w", name, err)
	}
	return intValue, nil
}
