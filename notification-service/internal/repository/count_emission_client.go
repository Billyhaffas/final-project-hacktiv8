package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UserDailyEmission struct {
	UserId             int8    `json:"UserId"`
	Date               string  `json:"Date"`
	TotalEmissionKgCo2 float64 `json:"TotalEmissionKgCo2"`
}

type EmissionClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewEmissionClient(baseURL string) *EmissionClient {
	return &EmissionClient{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    baseURL,
	}
}

type getDailyUserEmissionTypeRespon struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Data    *UserDailyEmission `json:"data"`
}

func (c *EmissionClient) GetDailyEmission(ctx context.Context, userID int) (*UserDailyEmission, error) {
	url := fmt.Sprintf("%s/api/v1/emissions/daily?user_id=%d", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("emission service returned status: %d", resp.StatusCode)
	}

	var apiRes getDailyUserEmissionTypeRespon
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		return nil, err
	}

	if apiRes.Data == nil {
		return nil, fmt.Errorf("no emission data found for user")
	}

	return apiRes.Data, nil
}
