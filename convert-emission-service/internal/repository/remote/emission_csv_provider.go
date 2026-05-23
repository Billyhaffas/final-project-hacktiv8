package remote

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"convert-emission-service/internal/domain"
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

	req.Header.Set("User-Agent", "convert-emission-service/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch carbon price CSV, status code: %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)

	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var emissions []domain.Emission

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

		carbonPrice, _ := strconv.ParseFloat(record[3], 64)

		uniqueGroup := code
		if uniqueGroup == "" {
			uniqueGroup = entity
		}

		computedID := fmt.Sprintf("%s-%d", uniqueGroup, year)

		emissions = append(emissions, domain.Emission{
			Id:          computedID,
			Entity:      entity,
			Code:        code,
			Year:        year,
			CarbonPrice: carbonPrice,
		})
	}

	return emissions, nil
}
