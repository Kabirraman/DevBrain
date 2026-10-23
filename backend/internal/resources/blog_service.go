package resources

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ExtractBlogContent(ctx context.Context, rawURL string) (string, string, string, string, error) {
	parsed, err := parseImportURL(rawURL)
	if err != nil {
		return "", "", "", "", err
	}
	if err := validatePublicURL(parsed); err != nil {
		return "", "", "", "", err
	}

	client := safeHTTPClient()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return "", "", "", "", err
	}
	request.Header.Set("User-Agent", "DevBrain/1.0 (+https://devbrain.app; resource importer)")
	resp, err := client.Do(request)
	if err != nil {
		return "", "", "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", "", "", "", fmt.Errorf("source returned HTTP %d", resp.StatusCode)
	}
	if resp.ContentLength > 5*1024*1024 {
		return "", "", "", "", fmt.Errorf("page is larger than the 5 MB import limit")
	}

	page, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024+1))
	if err != nil || len(page) > 5*1024*1024 {
		return "", "", "", "", fmt.Errorf("page is larger than the 5 MB import limit")
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
	if err != nil {
		return "", "", "", "", err
	}

	title := doc.Find("title").Text()
	title = strings.Join(strings.Fields(title), " ")
	author := strings.TrimSpace(doc.Find(`meta[name="author"]`).AttrOr("content", ""))
	if author == "" {
		author = strings.TrimSpace(doc.Find(`meta[property="article:author"]`).AttrOr("content", ""))
	}
	description := strings.TrimSpace(doc.Find(`meta[name="description"]`).AttrOr("content", ""))

	doc.Find("script, style, nav, footer, header, noscript").Remove()
	content := strings.Join(strings.Fields(doc.Find("main").Text()), " ")
	if content == "" {
		content = strings.Join(strings.Fields(doc.Find("article").Text()), " ")
	}
	if content == "" {
		content = strings.Join(strings.Fields(doc.Find("body").Text()), " ")
	}
	if len(content) > 2*1024*1024 {
		content = content[:2*1024*1024]
	}
	if content == "" {
		return "", "", "", "", fmt.Errorf("no readable article content was found")
	}

	return title, author, description, content, nil
}

func parseImportURL(rawURL string) (*url.URL, error) {
	parsed, err := url.ParseRequestURI(strings.TrimSpace(rawURL))
	if err != nil {
		return nil, fmt.Errorf("enter a valid public URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("only HTTP and HTTPS URLs can be imported")
	}
	if parsed.Hostname() == "" || parsed.User != nil {
		return nil, fmt.Errorf("enter a valid public URL")
	}
	return parsed, nil
}

func validatePublicURL(parsed *url.URL) error {
	if ip := net.ParseIP(parsed.Hostname()); ip != nil {
		if !isPublicIP(ip) {
			return fmt.Errorf("private network URLs cannot be imported")
		}
		return nil
	}
	addresses, err := net.LookupIP(parsed.Hostname())
	if err != nil {
		return fmt.Errorf("could not resolve the source host")
	}
	if len(addresses) == 0 {
		return fmt.Errorf("could not resolve the source host")
	}
	for _, address := range addresses {
		if !isPublicIP(address) {
			return fmt.Errorf("private network URLs cannot be imported")
		}
	}
	return nil
}

func isPublicIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
		ip = ip4
		if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast() {
			return false
		}
		if ip[0] == 100 && ip[1] >= 64 && ip[1] <= 127 {
			return false
		}
		if ip[0] == 169 && ip[1] == 254 {
			return false
		}
		if ip[0] == 192 && ip[1] == 0 && (ip[2] == 0 || ip[2] == 2) {
			return false
		}
		if ip[0] == 192 && ip[1] == 88 && ip[2] == 99 {
			return false
		}
		if ip[0] == 198 && (ip[1] == 18 || ip[1] == 19 || ip[1] == 51 && ip[2] == 100) {
			return false
		}
		if ip[0] == 203 && ip[1] == 0 && ip[2] == 113 {
			return false
		}
		return true
	}
	if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast() {
		return false
	}
	if len(ip) == net.IPv6len {
		if ip[0]&0xfe == 0xfc {
			return false
		}
		if ip[0] == 0x20 && ip[1] == 0x01 && ip[2] == 0x0d && ip[3] == 0xb8 {
			return false
		}
	}
	return true
}

func safeHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 12 * time.Second,
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(address)
				if err != nil {
					return nil, err
				}
				ips, err := net.LookupIP(host)
				if err != nil {
					return nil, err
				}
				var publicIPs []net.IP
				for _, ip := range ips {
					if isPublicIP(ip) {
						publicIPs = append(publicIPs, ip)
					}
				}
				if len(publicIPs) == 0 {
					return nil, fmt.Errorf("private network URLs cannot be imported")
				}
				dialer := net.Dialer{Timeout: 5 * time.Second}
				return dialer.DialContext(ctx, network, net.JoinHostPort(publicIPs[0].String(), port))
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return fmt.Errorf("too many redirects")
			}
			return validatePublicURL(req.URL)
		},
	}
}
