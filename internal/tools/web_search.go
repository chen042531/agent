package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/ycchen/ai-agent-go/internal/ollama"
)

type WebSearchTool struct{}

func (t *WebSearchTool) Name() string {
	return "web_search"
}

func (t *WebSearchTool) Description() string {
	return "Searches the web using DuckDuckGo and returns top results"
}

func (t *WebSearchTool) Parameters() ollama.ParameterSchema {
	return ollama.ParameterSchema{
		Type: "object",
		Properties: map[string]ollama.PropertyDefinition{
			"query": {
				Type:        "string",
				Description: "Search query string",
			},
		},
		Required: []string{"query"},
	}
}

type SearchResult struct {
	Title   string
	Snippet string
	URL     string
}

func (t *WebSearchTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("query must be a string")
	}

	results, err := t.searchDuckDuckGo(ctx, query)
	if err != nil {
		return "", fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		return "No results found", nil
	}

	// Format results
	var output strings.Builder
	output.WriteString(fmt.Sprintf("Search results for '%s':\n\n", query))
	for i, result := range results {
		if i >= 5 {
			break
		}
		output.WriteString(fmt.Sprintf("%d. %s\n", i+1, result.Title))
		output.WriteString(fmt.Sprintf("   %s\n", result.Snippet))
		output.WriteString(fmt.Sprintf("   %s\n\n", result.URL))
	}

	return output.String(), nil
}

func (t *WebSearchTool) searchDuckDuckGo(ctx context.Context, query string) ([]SearchResult, error) {
	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; AI-Agent/1.0)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	html := string(body)
	return t.parseResults(html), nil
}

func (t *WebSearchTool) parseResults(html string) []SearchResult {
	var results []SearchResult

	// Simple regex-based parsing
	// Match result blocks
	resultPattern := regexp.MustCompile(`(?s)<div class="result[^"]*">.*?</div>`)
	titlePattern := regexp.MustCompile(`<a class="result__a"[^>]*>(.*?)</a>`)
	urlPattern := regexp.MustCompile(`<a class="result__url"[^>]*href="([^"]+)"`)
	snippetPattern := regexp.MustCompile(`<a class="result__snippet"[^>]*>(.*?)</a>`)

	blocks := resultPattern.FindAllString(html, -1)
	for _, block := range blocks {
		var result SearchResult

		if matches := titlePattern.FindStringSubmatch(block); len(matches) > 1 {
			result.Title = t.stripHTML(matches[1])
		}
		if matches := urlPattern.FindStringSubmatch(block); len(matches) > 1 {
			result.URL = matches[1]
		}
		if matches := snippetPattern.FindStringSubmatch(block); len(matches) > 1 {
			result.Snippet = t.stripHTML(matches[1])
		}

		if result.Title != "" && result.URL != "" {
			results = append(results, result)
		}
	}

	return results
}

func (t *WebSearchTool) stripHTML(s string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	s = re.ReplaceAllString(s, "")
	// Decode HTML entities
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	return strings.TrimSpace(s)
}
