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
	username         string
	password         string
	serverURL        string
	siteID           int
	folderID         int
	secretTemplateID int
	dataFieldID      int
}

// secretServerEnabled reports whether the SecretServer e2e suite should run.
// SECRETSERVER_ENABLED is an explicit opt-in switch: unless it is set to a
// truthy value (e.g. "true") the specs skip, so the suite never runs unless a
// maintainer deliberately turns it on. When enabled, loadConfigFromEnv surfaces
// any missing or invalid configuration as a failure rather than a silent skip.
func secretServerEnabled() bool {
	enabled, _ := strconv.ParseBool(os.Getenv("SECRETSERVER_ENABLED"))
	return enabled
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

	// Instance-specific IDs. These have no defaults: values differ per Secret
	// Server instance, so discover them for the target instance (the example
	// values shown are only illustrative). SECRETSERVER_DATA_FIELD_ID must be a
	// non-file template field, which holds the JSON payload.
	cfg.siteID, err = getIntEnv("SECRETSERVER_SITE_ID") // e.g. 1
	if err != nil {
		return nil, err
	}
	cfg.folderID, err = getIntEnv("SECRETSERVER_FOLDER_ID") // e.g. 14
	if err != nil {
		return nil, err
	}
	cfg.secretTemplateID, err = getIntEnv("SECRETSERVER_TEMPLATE_ID") // e.g. 2 (a "Password" template)
	if err != nil {
		return nil, err
	}
	cfg.dataFieldID, err = getIntEnv("SECRETSERVER_DATA_FIELD_ID") // e.g. 60 (a non-file field)
	if err != nil {
		return nil, err
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

func getIntEnv(name string) (int, error) {
	value, ok := os.LookupEnv(name)
	if !ok || value == "" {
		return 0, fmt.Errorf("environment variable %q is not set", name)
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %q must be an integer: %w", name, err)
	}
	return intValue, nil
}
