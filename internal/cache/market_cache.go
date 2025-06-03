package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// MarketCache interface for caching market data
type MarketCache interface {
	SetMarkets(exchange string, markets interface{}, ttl time.Duration) error
	GetMarkets(exchange string, v interface{}) error
	Exists(exchange string) (bool, error)
}

// FileMarketCache implements file-based caching (simple, no dependencies)
type FileMarketCache struct {
	dir string
}

func NewFileMarketCache(dir string) (*FileMarketCache, error) {
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}
	return &FileMarketCache{dir: dir}, nil
}

func (f *FileMarketCache) SetMarkets(exchange string, markets interface{}, ttl time.Duration) error {
	data, err := json.MarshalIndent(markets, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal markets: %w", err)
	}

	filename := f.getFilename(exchange)
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	// Also write a metadata file with expiration time
	metadata := map[string]interface{}{
		"created": time.Now().Unix(),
		"ttl":     ttl.Seconds(),
		"expires": time.Now().Add(ttl).Unix(),
	}
	metaData, _ := json.Marshal(metadata)
	metaFile := filename + ".meta"
	os.WriteFile(metaFile, metaData, 0644)

	log.Printf("[Cache] Saved %s markets to %s (TTL: %v)", exchange, filename, ttl)
	return nil
}

func (f *FileMarketCache) GetMarkets(exchange string, v interface{}) error {
	filename := f.getFilename(exchange)

	// Check if cache is expired
	metaFile := filename + ".meta"
	if metaData, err := os.ReadFile(metaFile); err == nil {
		var meta map[string]interface{}
		if json.Unmarshal(metaData, &meta) == nil {
			if expires, ok := meta["expires"].(float64); ok {
				if time.Now().Unix() > int64(expires) {
					log.Printf("[Cache] Cache for %s is expired", exchange)
					return fmt.Errorf("cache expired")
				}
			}
		}
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cache miss: %w", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal cache: %w", err)
	}

	log.Printf("[Cache] Loaded %s markets from cache", exchange)
	return nil
}

func (f *FileMarketCache) Exists(exchange string) (bool, error) {
	filename := f.getFilename(exchange)
	_, err := os.Stat(filename)
	return err == nil, nil
}

func (f *FileMarketCache) getFilename(exchange string) string {
	return fmt.Sprintf("%s/%s_markets.json", f.dir, exchange)
}

// CachedMarketLoader wraps market loading with caching
type CachedMarketLoader struct {
	cache MarketCache
	ttl   time.Duration
}

func NewCachedMarketLoader(cache MarketCache, ttl time.Duration) *CachedMarketLoader {
	return &CachedMarketLoader{
		cache: cache,
		ttl:   ttl,
	}
}

func (c *CachedMarketLoader) LoadDeriveMarkets() (map[string]DeriveInstrument, error) {
	// Try to load from cache first
	var markets map[string]DeriveInstrument
	if err := c.cache.GetMarkets("derive", &markets); err == nil {
		return markets, nil
	}

	// Cache miss or expired - load from API
	log.Printf("[Cache] Cache miss for Derive markets, loading from API...")
	markets, err := LoadAllDeriveMarkets()
	if err != nil {
		return nil, err
	}

	// Save to cache
	if err := c.cache.SetMarkets("derive", markets, c.ttl); err != nil {
		log.Printf("[Cache] Warning: Failed to save markets to cache: %v", err)
	}

	return markets, nil
}
