package repositories

import (
	"fmt"
	"os"
	"strings"
)

// SQL 파일에서 쿼리 로드
func LoadQueries(filename string) (map[string]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read query file: %w", err)
	}

	queries := make(map[string]string)
	content := string(data)
	sections := strings.Split(content, "-- name: ")

	for _, section := range sections {
		lines := strings.SplitN(section, "\n", 2)
		if len(lines) < 2 {
			continue
		}
		name := strings.TrimSpace(lines[0])
		query := strings.TrimSpace(lines[1])
		queries[name] = query
	}
	return queries, nil
}
