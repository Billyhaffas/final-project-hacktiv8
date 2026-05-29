package externalapi

import (
	"count-emission-service/internal/domain"
	"count-emission-service/internal/model/thirdparty/carbonsutra"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type CarbonSutraExternalRepo struct {
	client *http.Client
}

func NewCarbonSutraRepository(client *http.Client) domain.CarbonSutraRepository {
	return &CarbonSutraExternalRepo{client: client}
}

func (repo *CarbonSutraExternalRepo) GetCarbonEmission(payload carbonsutra.CountEmisionBodyPayload) (*carbonsutra.CountEmisionThirdParty, error) {
	urlEndpoint := "https://carbonsutra1.p.rapidapi.com/vehicle_estimate_by_type"
	data := url.Values{}
	data.Set("vehicle_type", payload.VehicleType)
	data.Set("fuel_type", payload.FuelType)
	data.Set("distance_value", fmt.Sprintf("%f", payload.DistanceValue))
	data.Set("distance_unit", payload.DistanceUnit)
	data.Set("include_wtt", payload.IncludeWtt)
	data.Set("cluster_name", "")

	req, err := http.NewRequest(http.MethodPost, urlEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Go-http-client/1.1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("CARBON_TOKEN"))
	req.Header.Set("x-rapidapi-host", "carbonsutra1.p.rapidapi.com")
	req.Header.Set("x-rapidapi-key", os.Getenv("RAPID_API_KEY"))

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result carbonsutra.CountEmisionThirdParty

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CarbonSutra API error: HTTP %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
