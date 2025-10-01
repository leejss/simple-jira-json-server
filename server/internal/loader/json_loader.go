package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/leejss/simple-json-server/server/internal/models"
)

type FileLoader struct {
	path  string
	cache map[string][]models.Issue
	mu    sync.RWMutex
}

func NewFileLoader(path string) *FileLoader {
	return &FileLoader{
		path:  path,
		cache: make(map[string][]models.Issue),
	}
}

func (f *FileLoader) LoadIssues(year string) ([]models.Issue, error) {
	cacheKey := fmt.Sprintf("issues_%s", year)

	f.mu.RLock()

	if cached, exists := f.cache[cacheKey]; exists {
		f.mu.RUnlock()
		return cached, nil
	}

	f.mu.RUnlock()

	filename := filepath.Join(f.path, fmt.Sprintf("issues_%s.json", year))
	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var issues []models.Issue
	if err := json.Unmarshal(data, &issues); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	f.mu.Lock()
	f.cache[cacheKey] = issues
	f.mu.Unlock()

	return issues, nil
}
