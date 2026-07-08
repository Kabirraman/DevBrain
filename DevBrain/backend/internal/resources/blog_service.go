package resources

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExtractBlogContent(url string) (string, string, error) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", "", err
	}

	// Some sites (including Wikipedia) serve stripped-down or blocked
	// content to Go's default User-Agent. Look like a normal browser.
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36",
	)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return "", "", err
	}

	title := strings.TrimSpace(doc.Find("title").Text())

	if title == "" {
		title, _ = doc.Find("meta[property='og:title']").Attr("content")
		title = strings.TrimSpace(title)
	}

	content := extractText(doc, "article")

	if content == "" {
		content = extractText(doc, "main")
	}

	if content == "" {
		content = extractText(doc, "#mw-content-text")
	}

	if content == "" {
		content = extractText(doc, "body")
	}

	return title, content, nil
}

func extractText(doc *goquery.Document, selector string) string {

	sel := doc.Find(selector)

	// Strip elements that add noise but no learning value.
	sel.Find("script, style, nav, footer, header, noscript").Remove()

	text := sel.Text()

	// Collapse runs of whitespace left behind by removed elements.
	fields := strings.Fields(text)

	return strings.Join(fields, " ")
}