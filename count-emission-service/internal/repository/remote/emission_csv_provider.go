package remote

import (
	"context"
	"count-emission-service/internal/domain"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type CSVProvider struct {
	url string
}

func NewCSVProvider(url string) domain.EmissionProvider {
	return &CSVProvider{url: url}
}

func (p *CSVProvider) FetchEmissions(ctx context.Context) ([]domain.Emission, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch CSV, status code: %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)

	// Skip header row (Entity, Code, Year, "Annual CO₂ emissions (per capita)")
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var emissions []domain.Emission
	now := time.Now().UTC()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) < 4 {
			continue
		}

		entity := record[0]
		code := record[1]
		year, _ := strconv.Atoi(record[2])
		perCapitaAnnualTonnes, _ := strconv.ParseFloat(record[3], 64)

		// Fallback logic if regional code is missing
		countryCode := code
		if countryCode == "" {
			countryCode = entity
		}

		// Calculate Daily Limit in kg from Annual Tonnes:
		// (Tonnes * 1000) / 365.25 days
		dailyLimitKg := (perCapitaAnnualTonnes * 1000.0) / 365.25

		// Establish primary key pattern based on context requirements
		computedID := fmt.Sprintf("%s-%d", countryCode, year)

		emissions = append(emissions, domain.Emission{
			Id:              computedID,
			CountryCode:     countryCode,
			DailyLimitKgCo2: dailyLimitKg,
			SourceUrl:       p.url,
			UpdatedAt:       now,
		})
	}

	return emissions, nil
}
