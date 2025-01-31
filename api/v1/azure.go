package v1

import "github.com/flanksource/duty/types"

type AzureDevops struct {
	BaseScraper         `json:",inline"`
	Organization        string       `yaml:"organization" json:"organization"`
	PersonalAccessToken types.EnvVar `yaml:"personalAccessToken" json:"personalAccessToken"`
	Projects            []string     `yaml:"projects" json:"projects"`
	Pipelines           []string     `yaml:"pipelines" json:"pipelines"`
}
type Azure struct {
	BaseScraper    `json:",inline"`
	ConnectionName string       `yaml:"connection,omitempty" json:"connection,omitempty"`
	SubscriptionID string       `yaml:"subscriptionID" json:"subscriptionID"`
	Organisation   string       `yaml:"organisation" json:"organisation"`
	ClientID       types.EnvVar `yaml:"clientID,omitempty" json:"clientID,omitempty"`
	ClientSecret   types.EnvVar `yaml:"clientSecret,omitempty" json:"clientSecret,omitempty"`
	TenantID       string       `yaml:"tenantID" json:"tenantID"`
}
