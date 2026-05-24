package remote

import (
	"context"
	"convert-emission-service/internal/domain"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type CSVPriceProvider struct {
	url string
}

func NewCSVPriceProvider(url string) domain.CarbonPriceProvider {
	return &CSVPriceProvider{url: url}
}

func (p *CSVPriceProvider) FetchPrices(ctx context.Context) ([]domain.CarbonPrice, error) {
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
		return nil, fmt.Errorf("failed to fetch CSV data, status: %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)

	// Skip CSV header (Entity, Code, Year, Emissions-weighted carbon price...)
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var records []domain.CarbonPrice
	now := time.Now().UTC()
	const defaultIdrCurRate = 16250.0

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
		yearStr := record[2]
		priceStr := record[3]

		priceUsd, _ := strconv.ParseFloat(priceStr, 64)

		countryIdentifier := code
		if countryIdentifier == "" {
			countryIdentifier = entity
		}

		computedID := fmt.Sprintf("%s-%s", countryIdentifier, yearStr)

		records = append(records, domain.CarbonPrice{
			Id:             computedID,
			PricePerTonUsd: priceUsd,
			UsdCurRate:     defaultIdrCurRate,
			Source:         p.url,
			FetchedAt:      now,
		})
	}

	return records, nil
}
