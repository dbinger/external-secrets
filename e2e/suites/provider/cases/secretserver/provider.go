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
	"encoding/json"

	"github.com/DelineaXPM/tss-sdk-go/v3/server"
	"github.com/external-secrets/external-secrets-e2e/framework"
	"github.com/onsi/gomega"
)

type secretStoreProvider struct {
	api       *server.Server
	cfg       *config
	framework *framework.Framework
	secretID  map[string]int
}

func (p *secretStoreProvider) init(cfg *config, f *framework.Framework) {
	p.cfg = cfg
	p.secretID = make(map[string]int)
	p.framework = f
	secretserverClient, err := server.New(server.Configuration{
		Credentials: server.UserCredential{
			Username: cfg.username,
			Password: cfg.password,
		},
		ServerURL: cfg.serverURL,
	})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	p.api = secretserverClient
}

func (p *secretStoreProvider) CreateSecret(key string, val framework.SecretEntry) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(val.Value), &data)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	fields := []server.SecretField{
		{
			FieldID:   p.cfg.dataFieldID,
			ItemValue: val.Value,
		},
	}

	if p.cfg.requiredPasswordFieldID > 0 {
		fields = append(fields, server.SecretField{
			FieldID:   p.cfg.requiredPasswordFieldID,
			ItemValue: p.cfg.requiredPasswordValue,
		})
	}

	template, err := p.api.SecretTemplate(p.cfg.secretTemplateID)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	for _, field := range template.Fields {
		if field.IsRequired {
			fields = append(fields, server.SecretField{
				FieldID:   field.SecretTemplateFieldID,
				ItemValue: p.cfg.requiredPasswordValue,
			})
		}
	}
	fields = uniqueFields(fields)

	s, err := p.api.CreateSecret(server.Secret{
		SecretTemplateID: p.cfg.secretTemplateID,
		SiteID:           p.cfg.siteID,
		FolderID:         p.cfg.folderID,
		Name:             key,
		Fields:           fields,
	})
	gomega.Expect(err).ToNot(gomega.HaveOccurred(),
		"failed creating SecretServer secret with site=%d folder=%d template=%d dataField=%d fieldIDs=%v",
		p.cfg.siteID, p.cfg.folderID, p.cfg.secretTemplateID, p.cfg.dataFieldID, fieldIDs(fields))
	p.secretID[key] = s.ID
}

func (p *secretStoreProvider) DeleteSecret(key string) {
	err := p.api.DeleteSecret(p.secretID[key])
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

func uniqueFields(fields []server.SecretField) []server.SecretField {
	seen := make(map[int]struct{}, len(fields))
	result := make([]server.SecretField, 0, len(fields))
	for _, field := range fields {
		if _, ok := seen[field.FieldID]; ok {
			continue
		}
		seen[field.FieldID] = struct{}{}
		result = append(result, field)
	}
	return result
}

func fieldIDs(fields []server.SecretField) []int {
	ids := make([]int, 0, len(fields))
	for _, field := range fields {
		ids = append(ids, field.FieldID)
	}
	return ids
}
