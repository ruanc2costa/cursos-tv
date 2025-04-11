package models

import "time"

// TokenRequest representa a estrutura para solicitar token de acesso OAuth
type TokenRequest struct {
	ClientID     string `json:"client_id"`
	Scope        string `json:"scope"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

// TokenResponse representa a resposta do endpoint de autenticação OAuth
type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// GenerateTokenRequest representa a estrutura para solicitar token de incorporação
type GenerateTokenRequest struct {
	AccessLevel       string              `json:"accessLevel"`
	LifetimeInMinutes int                 `json:"lifetimeInMinutes,omitempty"`
	Identities        []EffectiveIdentity `json:"identities,omitempty"`
}

// EffectiveIdentity representa uma identidade para RLS (segurança em nível de linha)
type EffectiveIdentity struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles,omitempty"`
	Datasets []string `json:"datasets"`
}

// EmbedToken representa o token de incorporação retornado pela API do Power BI
type EmbedToken struct {
	Token      string    `json:"token"`
	TokenID    string    `json:"tokenId"`
	Expiration time.Time `json:"expiration"`
}

// EmbedConfig representa a configuração completa para incorporação do Power BI
type EmbedConfig struct {
	EmbedToken  EmbedToken `json:"embedToken"`
	EmbedURL    string     `json:"embedUrl"`
	DashboardID string     `json:"dashboardId"`
	ReportID    string     `json:"reportId,omitempty"`
	GroupID     string     `json:"groupId"`
	Type        string     `json:"type"`
}
