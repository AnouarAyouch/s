package currencycache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type cacheEntry struct {
	CreatedAt time.Time `json:"created_at"`
	Value     string    `json:"value"`
}
type Cache struct {
	mu       sync.Mutex
	cacheMap map[string]cacheEntry
	ttl      time.Duration
	filePath string
}

func NewCache(ttl time.Duration) *Cache {
	currentDir, err := os.Getwd()
	if err != nil {

		panic(fmt.Sprintf("failed to get current directory: %v", err))
	}
	file_path := filepath.Join(currentDir, "cache.json")

	c := &Cache{
		cacheMap: make(map[string]cacheEntry),
		ttl:      ttl,
		filePath: file_path,
	}
	c.loadFromFile()
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheMap[key] = cacheEntry{
		CreatedAt: time.Now(),
		Value:     string(val),
	}

	c.saveToFile()
}
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, found := c.cacheMap[key]
	if !found {
		return nil, false
	}

	if time.Since(entry.CreatedAt) > c.ttl {
		delete(c.cacheMap, key)
		c.saveToFile()
		return nil, false
	}
	return []byte(entry.Value), true
}
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(time.Hour) // check every hour
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		changed := false
		for k, v := range c.cacheMap {
			if now.Sub(v.CreatedAt) > c.ttl {
				delete(c.cacheMap, k)
				changed = true
			}
		}
		if changed {
			c.saveToFile()
		}
		c.mu.Unlock()
	}
}

func (c *Cache) saveToFile() {
	data, err := json.MarshalIndent(c.cacheMap, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling cache:", err)
		return
	}
	if err := os.WriteFile(c.filePath, data, 0644); err != nil {
		fmt.Println("Error writing cache file:", err)
	}
}

func (c *Cache) loadFromFile() {
	file, err := os.ReadFile(c.filePath)
	if err != nil {
		return // file doesn't exist yet
	}

	var data map[string]cacheEntry
	if err := json.Unmarshal(file, &data); err != nil {
		fmt.Println("Error unmarshaling cache file:", err)
		return
	}

	now := time.Now()
	for k, v := range data {
		if now.Sub(v.CreatedAt) <= c.ttl {
			c.cacheMap[k] = v
		}
	}
}
