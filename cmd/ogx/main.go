package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
	xhtml "golang.org/x/net/html"
)

const (
	defaultTimeout = 10 * time.Second

	// maxBodyBytes limits parsed HTML to ~2 MiB; OG tags are usually in <head>.
	maxBodyBytes = 2 << 20

	// userAgent identifies this CLI to upstream servers/proxies.
	userAgent = "utils-og-check/1.0"
	ogPrefix  = "og:"
)

var (
	timeout    time.Duration
	jsonOutput bool
)

// result contains the inspected URL and extracted OG tags.
type result struct {
	url  string
	tags map[string][]string
}

func main() {
	cmd := &cobra.Command{
		Use:     "ogx <url>",
		Short:   "Inspect OpenGraph metadata for a URL",
		Example: "  ogx https://example.com\n  ogx example.com\n  ogx example.com --json",
		Args:    cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			if err := run(args[0]); err != nil {
				fatalf("%v", err)
			}
		},
	}

	cmd.Flags().DurationVar(&timeout, "timeout", defaultTimeout, "HTTP request timeout")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")

	if err := cmd.Execute(); err != nil {
		fatalf("%v", err)
	}
}

// fatalf prints an error message and exits with status 1.
func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

// run executes a single OG inspection request.
func run(rawURL string) error {
	target, err := parseURL(rawURL)
	if err != nil {
		return err
	}

	tags, err := fetchOG(target, timeout)
	if err != nil {
		return err
	}

	r := result{url: target, tags: tags}
	if jsonOutput {
		return printJSON(r)
	}
	printText(r)
	return nil
}

// parseURL trims input, adds a default scheme, and validates the URL.
func parseURL(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("url cannot be empty")
	}
	if !strings.Contains(s, "://") {
		s = "https://" + s
	}

	u, err := url.ParseRequestURI(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid URL %q", s)
	}
	return u.String(), nil
}

// fetchOG fetches a page and extracts og:* metadata.
func fetchOG(target string, timeout time.Duration) (map[string][]string, error) {
	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		return nil, fmt.Errorf("unexpected HTTP status: %s", res.Status)
	}

	return parseOG(io.LimitReader(res.Body, maxBodyBytes)), nil
}

// printText writes the human-readable report to stdout.
func printText(r result) {
	fmt.Printf("URL: %s\n\n", r.url)
	if len(r.tags) == 0 {
		fmt.Println("No OpenGraph tags found.")
		return
	}

	for _, key := range sortedKeys(r.tags) {
		for _, value := range r.tags[key] {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
}

// printJSON writes machine-readable output to stdout.
func printJSON(r result) error {
	payload := struct {
		URL  string              `json:"url"`
		Tags map[string][]string `json:"tags"`
	}{
		URL:  r.url,
		Tags: r.tags,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

// parseOG extracts all og:* meta values from an HTML document stream.
func parseOG(r io.Reader) map[string][]string {
	tags := make(map[string][]string)
	tokenizer := xhtml.NewTokenizer(r)

	for {
		tokenType := tokenizer.Next()
		if tokenType == xhtml.ErrorToken {
			return tags
		}
		if tokenType != xhtml.StartTagToken && tokenType != xhtml.SelfClosingTagToken {
			continue
		}

		name, hasAttr := tokenizer.TagName()
		if !strings.EqualFold(string(name), "meta") {
			continue
		}

		var key, content string
		for hasAttr {
			attrKey, attrValue, more := tokenizer.TagAttr()
			hasAttr = more

			switch {
			case strings.EqualFold(string(attrKey), "property"), strings.EqualFold(string(attrKey), "name"):
				if key == "" {
					key = strings.ToLower(strings.TrimSpace(string(attrValue)))
				}
			case strings.EqualFold(string(attrKey), "content"):
				content = strings.TrimSpace(html.UnescapeString(string(attrValue)))
			}
		}

		if key == "" || content == "" || !strings.HasPrefix(key, ogPrefix) {
			continue
		}
		tags[key] = append(tags[key], content)
	}
}

// sortedKeys returns map keys in lexical order.
func sortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
