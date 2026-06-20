package resources

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func ExtractBlogContent(url string) (string, string, error) {

	resp, err := http.Get(url)

	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return "", "", err
	}

	title := doc.Find("title").Text()

	content := doc.Find("main").Text()

if content == "" {
	content = doc.Find("article").Text()
}

if content == "" {
	content = doc.Find("body").Text()
}

	return title, content, nil
}