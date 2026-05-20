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
	"strings"
	"testing"
)

func TestLoadConfigFromEnvUsesDefaultsForEmptyOptionalValues(t *testing.T) {
	clearSecretServerEnv(t)

	t.Setenv("SECRETSERVER_USERNAME", "username")
	t.Setenv("SECRETSERVER_PASSWORD", "password")
	t.Setenv("SECRETSERVER_URL", "https://example.com")
	t.Setenv("SECRETSERVER_SITE_ID", "")
	t.Setenv("SECRETSERVER_FOLDER_ID", "")
	t.Setenv("SECRETSERVER_TEMPLATE_ID", "")
	t.Setenv("SECRETSERVER_DATA_FIELD_ID", "")
	t.Setenv("SECRETSERVER_REQUIRED_PASSWORD_FIELD_ID", "")
	t.Setenv("SECRETSERVER_REQUIRED_PASSWORD_VALUE", "")

	cfg, err := loadConfigFromEnv()
	if err != nil {
		t.Fatalf("loadConfigFromEnv() returned error: %v", err)
	}

	if cfg.siteID != 1 {
		t.Fatalf("siteID = %d, want 1", cfg.siteID)
	}
	if cfg.folderID != 10 {
		t.Fatalf("folderID = %d, want 10", cfg.folderID)
	}
	if cfg.secretTemplateID != 6051 {
		t.Fatalf("secretTemplateID = %d, want 6051", cfg.secretTemplateID)
	}
	if cfg.dataFieldID != 329 {
		t.Fatalf("dataFieldID = %d, want 329", cfg.dataFieldID)
	}
	if cfg.requiredPasswordFieldID != 0 {
		t.Fatalf("requiredPasswordFieldID = %d, want 0", cfg.requiredPasswordFieldID)
	}
	if cfg.requiredPasswordValue != "external-secrets-e2e" {
		t.Fatalf("requiredPasswordValue = %q, want default", cfg.requiredPasswordValue)
	}
}

func TestLoadConfigFromEnvUsesOptionalValues(t *testing.T) {
	clearSecretServerEnv(t)

	t.Setenv("SECRETSERVER_USERNAME", "username")
	t.Setenv("SECRETSERVER_PASSWORD", "password")
	t.Setenv("SECRETSERVER_URL", "https://example.com")
	t.Setenv("SECRETSERVER_SITE_ID", "2")
	t.Setenv("SECRETSERVER_FOLDER_ID", "14")
	t.Setenv("SECRETSERVER_TEMPLATE_ID", "2")
	t.Setenv("SECRETSERVER_DATA_FIELD_ID", "60")
	t.Setenv("SECRETSERVER_REQUIRED_PASSWORD_FIELD_ID", "7")
	t.Setenv("SECRETSERVER_REQUIRED_PASSWORD_VALUE", "required-value")

	cfg, err := loadConfigFromEnv()
	if err != nil {
		t.Fatalf("loadConfigFromEnv() returned error: %v", err)
	}

	if cfg.siteID != 2 {
		t.Fatalf("siteID = %d, want 2", cfg.siteID)
	}
	if cfg.folderID != 14 {
		t.Fatalf("folderID = %d, want 14", cfg.folderID)
	}
	if cfg.secretTemplateID != 2 {
		t.Fatalf("secretTemplateID = %d, want 2", cfg.secretTemplateID)
	}
	if cfg.dataFieldID != 60 {
		t.Fatalf("dataFieldID = %d, want 60", cfg.dataFieldID)
	}
	if cfg.requiredPasswordFieldID != 7 {
		t.Fatalf("requiredPasswordFieldID = %d, want 7", cfg.requiredPasswordFieldID)
	}
	if cfg.requiredPasswordValue != "required-value" {
		t.Fatalf("requiredPasswordValue = %q, want required-value", cfg.requiredPasswordValue)
	}
}

func TestLoadConfigFromEnvRejectsInvalidOptionalInt(t *testing.T) {
	clearSecretServerEnv(t)

	t.Setenv("SECRETSERVER_USERNAME", "username")
	t.Setenv("SECRETSERVER_PASSWORD", "password")
	t.Setenv("SECRETSERVER_URL", "https://example.com")
	t.Setenv("SECRETSERVER_FOLDER_ID", "invalid")

	_, err := loadConfigFromEnv()
	if err == nil {
		t.Fatal("loadConfigFromEnv() returned nil error, want invalid integer error")
	}
	if !strings.Contains(err.Error(), "SECRETSERVER_FOLDER_ID") {
		t.Fatalf("error = %q, want SECRETSERVER_FOLDER_ID", err)
	}
}

func clearSecretServerEnv(t *testing.T) {
	t.Helper()

	for _, name := range []string{
		"SECRETSERVER_USERNAME",
		"SECRETSERVER_PASSWORD",
		"SECRETSERVER_URL",
		"SECRETSERVER_SITE_ID",
		"SECRETSERVER_FOLDER_ID",
		"SECRETSERVER_TEMPLATE_ID",
		"SECRETSERVER_DATA_FIELD_ID",
		"SECRETSERVER_REQUIRED_PASSWORD_FIELD_ID",
		"SECRETSERVER_REQUIRED_PASSWORD_VALUE",
	} {
		t.Setenv(name, "")
	}
}
