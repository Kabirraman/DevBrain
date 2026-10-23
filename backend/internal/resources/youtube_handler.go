package resources

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type youtubeOEmbed struct {
	Title      string `json:"title"`
	AuthorName string `json:"author_name"`
}

type transcriptDocument struct {
	Lines []struct {
		Text string `xml:",chardata"`
	} `xml:"text"`
}

var youtubeIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

func ImportYouTubeHandler(c *gin.Context) {
	var request ImportBlogRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter a valid YouTube URL."})
		return
	}
	videoID, err := parseYouTubeID(request.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter a valid YouTube video URL."})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()
	metadata, err := fetchYouTubeMetadata(ctx, videoID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Couldn't load video details. Check the URL and try again."})
		return
	}
	description, _ := fetchYouTubeDescription(ctx, videoID)
	transcript, _ := fetchYouTubeTranscript(ctx, videoID)
	content := strings.TrimSpace(strings.Join([]string{transcript, description}, "\n\n"))
	if content == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "This video has no available transcript or description."})
		return
	}
	if len(content) > 2<<20 {
		content = content[:2<<20]
	}
	title := strings.TrimSpace(metadata.Title)
	if title == "" {
		title = "YouTube video"
	}
	err = CreateResource(c.MustGet("userID").(string), CreateResourceRequest{
		Title:       title,
		Type:        "youtube",
		Author:      metadata.AuthorName,
		Description: description,
		SourceURL:   "https://www.youtube.com/watch?v=" + videoID,
		Content:     content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't save this resource."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "YouTube video imported successfully"})
}

func parseYouTubeID(rawURL string) (string, error) {
	parsed, err := url.ParseRequestURI(strings.TrimSpace(rawURL))
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", fmt.Errorf("invalid URL")
	}
	host := strings.ToLower(strings.TrimPrefix(parsed.Hostname(), "www."))
	var id string
	switch host {
	case "youtube.com", "m.youtube.com":
		id = parsed.Query().Get("v")
		if id == "" {
			parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
			if len(parts) == 2 && (parts[0] == "shorts" || parts[0] == "embed") {
				id = parts[1]
			}
		}
	case "youtu.be":
		id = strings.Trim(parsed.Path, "/")
	}
	if !youtubeIDPattern.MatchString(id) {
		return "", fmt.Errorf("invalid video id")
	}
	return id, nil
}

func fetchYouTubeMetadata(ctx context.Context, videoID string) (youtubeOEmbed, error) {
	endpoint := "https://www.youtube.com/oembed?format=json&url=" + url.QueryEscape("https://www.youtube.com/watch?v="+videoID)
	var metadata youtubeOEmbed
	if err := youtubeGET(ctx, endpoint, &metadata); err != nil {
		return metadata, err
	}
	return metadata, nil
}

func fetchYouTubeDescription(ctx context.Context, videoID string) (string, error) {
	requestURL := "https://www.youtube.com/watch?v=" + url.QueryEscape(videoID)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (compatible; DevBrain/1.0)")
	response, err := safeHTTPClient().Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("video page returned HTTP %d", response.StatusCode)
	}
	page, err := io.ReadAll(io.LimitReader(response.Body, 2<<20))
	if err != nil {
		return "", err
	}
	match := regexp.MustCompile(`<meta name="description" content="([^"]*)"`).FindSubmatch(page)
	if len(match) < 2 {
		return "", nil
	}
	return strings.TrimSpace(strings.NewReplacer("&quot;", "\"", "&#39;", "'", "&amp;", "&", "&lt;", "<", "&gt;", ">").Replace(string(match[1]))), nil
}

func fetchYouTubeTranscript(ctx context.Context, videoID string) (string, error) {
	for _, language := range []string{"en", "en-US", "en-GB"} {
		for _, kind := range []string{"", "asr"} {
			query := url.Values{"fmt": {"srv3"}, "v": {videoID}, "lang": {language}}
			if kind != "" {
				query.Set("kind", kind)
			}
			endpoint := "https://www.youtube.com/api/timedtext?" + query.Encode()
			request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
			if err != nil {
				return "", err
			}
			request.Header.Set("User-Agent", "Mozilla/5.0 (compatible; DevBrain/1.0)")
			response, err := safeHTTPClient().Do(request)
			if err != nil {
				continue
			}
			data, readErr := io.ReadAll(io.LimitReader(response.Body, 2<<20))
			response.Body.Close()
			if readErr != nil || response.StatusCode < 200 || response.StatusCode >= 300 || len(data) == 0 {
				continue
			}
			var document transcriptDocument
			if xml.Unmarshal(data, &document) != nil {
				continue
			}
			lines := make([]string, 0, len(document.Lines))
			for _, line := range document.Lines {
				if text := strings.TrimSpace(line.Text); text != "" {
					lines = append(lines, text)
				}
			}
			if len(lines) > 0 {
				return strings.Join(lines, " "), nil
			}
		}
	}
	return "", nil
}

func youtubeGET(ctx context.Context, requestURL string, target any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (compatible; DevBrain/1.0)")
	response, err := safeHTTPClient().Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("YouTube returned HTTP %d", response.StatusCode)
	}
	return json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(target)
}
