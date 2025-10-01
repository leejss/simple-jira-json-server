package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func (f *FileLoader) LoadIssues(year int) ([]models.Issue, error) {
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

func (f *FileLoader) GetAvailableYears() ([]int, error) {
	entries, err := os.ReadDir(f.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var years []int

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "jira_") {
			continue
		}

		yearStr := strings.TrimPrefix(entry.Name(), "jira_")
		yearStr = strings.TrimSuffix(yearStr, ".json")

		if year, err := strconv.Atoi(yearStr); err == nil {
			years = append(years, year)
		}
	}

	return years, nil
}

func (f *FileLoader) LoadAllYears() ([]models.Issue, error) {
	years, err := f.GetAvailableYears()
	if err != nil {
		return nil, err
	}

	var allIssues []models.Issue
	for _, year := range years {
		issues, err := f.LoadIssues(year)
		if err != nil {
			continue
		}

		allIssues = append(allIssues, issues...)
	}

	return allIssues, nil
}

func (f *FileLoader) ClearCache() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.cache = make(map[string][]models.Issue)
}
