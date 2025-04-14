/*package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"tvtec/models"
)

type PowerBIService interface {
	GetDashboardEmbedToken(dashboardID, groupID, datasetID string) (*models.EmbedConfig, error)
}

type powerBIService struct {
	httpClient *http.Client
}

func NewPowerBIService() PowerBIService {
	return &powerBIService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *powerBIService) getAccessToken() (string, error) {
	// Obter credenciais das variáveis de ambiente
	tenantID := os.Getenv("POWERBI_TENANT_ID")
	clientID := os.Getenv("POWERBI_CLIENT_ID")
	clientSecret := os.Getenv("POWERBI_CLIENT_SECRET")

	if tenantID == "" || clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("variáveis de ambiente do Power BI não configuradas")
	}

	// Criar corpo da solicitação
	tokenReq := models.TokenRequest{
		ClientID:     clientID,
		Scope:        "https://analysis.windows.net/powerbi/api/.default",
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}

	jsonData, err := json.Marshal(tokenReq)
	if err != nil {
		return "", err
	}

	// Enviar solicitação
	url := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("erro na solicitação de token: %s", body)
	}

	// Decodificar resposta
	var tokenResp models.TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func (s *powerBIService) GetDashboardEmbedToken(dashboardID, groupID, datasetID string) (*models.EmbedConfig, error) {
	// Obter token de acesso
	accessToken, err := s.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("falha ao obter token de acesso: %v", err)
	}

	// Criar corpo da solicitação
	generateReq := models.GenerateTokenRequest{
		AccessLevel:       "View",
		LifetimeInMinutes: 60,
		Identities: []models.EffectiveIdentity{
			{
				Username: "powerbi-embedded-user",
				Roles:    []string{"Viewer"},
				Datasets: []string{datasetID},
			},
		},
	}

	jsonData, err := json.Marshal(generateReq)
	if err != nil {
		return nil, err
	}

	// Enviar solicitação
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dashboards/%s/GenerateToken", groupID, dashboardID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("erro na geração do token: %s", body)
	}

	// Decodificar resposta
	var embedToken models.EmbedToken
	err = json.Unmarshal(body, &embedToken)
	if err != nil {
		return nil, err
	}

	// Construir URL de incorporação
	embedURL := fmt.Sprintf("https://app.powerbi.com/dashboardEmbed?dashboardId=%s&groupId=%s", dashboardID, groupID)

	// Construir configuração completa
	embedConfig := &models.EmbedConfig{
		EmbedToken:  embedToken,
		EmbedURL:    embedURL,
		DashboardID: dashboardID,
		GroupID:     groupID,
		Type:        "dashboard",
	}

	return embedConfig, nil
}
*/