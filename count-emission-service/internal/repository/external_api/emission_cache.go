package externalapi

import (
	"count-emission-service/internal/domain"
	"count-emission-service/internal/model/thirdparty/carbonsutra"
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	value     *carbonsutra.CountEmisionThirdParty
	expiresAt time.Time
}

type CachedEmissionRepo struct {
	delegate domain.CarbonSutraRepository
	mu       sync.Mutex
	cache    map[string]cacheEntry
	ttl      time.Duration
}

func NewCachedEmissionRepo(delegate domain.CarbonSutraRepository) domain.CarbonSutraRepository {
	return &CachedEmissionRepo{
		delegate: delegate,
		cache:    make(map[string]cacheEntry),
		ttl:      24 * time.Hour,
	}
}

func (c *CachedEmissionRepo) GetCarbonEmission(payload carbonsutra.CountEmisionBodyPayload) (*carbonsutra.CountEmisionThirdParty, error) {
	key := fmt.Sprintf("%s|%s|%.4f|%s",
		payload.VehicleType,
		payload.FuelType,
		payload.DistanceValue,
		time.Now().Format("2006-01-02"),
	)

	c.mu.Lock()
	if entry, ok := c.cache[key]; ok && time.Now().Before(entry.expiresAt) {
		c.mu.Unlock()
		return entry.value, nil
	}
	c.mu.Unlock()

	result, err := c.delegate.GetCarbonEmission(payload)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cache[key] = cacheEntry{value: result, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()

	return result, nil
}
