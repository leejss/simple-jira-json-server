package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func Save(data []byte, outputPath string) error {

	dir := filepath.Dir(outputPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	var jsonBuffer bytes.Buffer

	if err := json.Indent(&jsonBuffer, data, "", "  "); err != nil {
		return fmt.Errorf("failed to indent JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, jsonBuffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println("JSON saved to:", outputPath)
	return nil
}
