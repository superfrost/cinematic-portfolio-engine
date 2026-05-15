package internal

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/mozillazg/go-unidecode"
)

var (
	youtubeRegex = regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/embed/)([a-zA-Z0-9_-]{11})`)
	rutubeRegex  = regexp.MustCompile(`rutube\.ru/(?:video|play/embed)/([a-zA-Z0-9]+)`)
	specChars    = regexp.MustCompile(`[^a-zA-Z0-9\s-]`)
	multiDash    = regexp.MustCompile(`-{2,}`)
)

func ParseYouTubeID(rawURL string) (string, error) {
	matches := youtubeRegex.FindStringSubmatch(rawURL)
	if len(matches) < 2 {
		return "", fmt.Errorf("YouTube ID not found in: %s", rawURL)
	}
	return matches[1], nil
}

func ParseRutubeID(rawURL string) (string, error) {
	matches := rutubeRegex.FindStringSubmatch(rawURL)
	if len(matches) < 2 {
		return "", fmt.Errorf("Rutube ID not found in: %s", rawURL)
	}
	return matches[1], nil
}

type VideoInfo struct {
	Provider string
	ID       string
}

func ParseVideoLink(rawURL string) (*VideoInfo, error) {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if id, err := ParseYouTubeID(u.String()); err == nil {
		return &VideoInfo{Provider: "youtube", ID: id}, nil
	}

	if id, err := ParseRutubeID(u.String()); err == nil {
		return &VideoInfo{Provider: "rutube", ID: id}, nil
	}

	return nil, fmt.Errorf("unsupported video link: %s", rawURL)
}

func Transliterate(name string) string {
	lat := unidecode.Unidecode(name)
	lat = strings.ToLower(lat)
	lat = specChars.ReplaceAllString(lat, "")
	lat = strings.ReplaceAll(lat, " ", "-")
	lat = multiDash.ReplaceAllString(lat, "-")
	lat = strings.Trim(lat, "-")
	return lat
}

var filenameRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})-(.*)\.json$`)

type ParsedFilename struct {
	Date string
	Slug string
}

func ParseFilename(name string) (*ParsedFilename, bool) {
	matches := filenameRegex.FindStringSubmatch(name)
	if len(matches) < 3 {
		return nil, false
	}
	return &ParsedFilename{Date: matches[1], Slug: matches[2]}, true
}
