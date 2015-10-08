package sitemapgen

import (
	"net/url"
	"testing"
)

func TestToAbsURL(t *testing.T) {
	parsed, err := url.Parse("http://example.com/somepage/someproduct?querystring=fd")
	if err != nil {
		t.Error(err)
	}
	abs, err := toAbsURL(parsed, "http://somesite.com/relative")
	if err != nil || abs.String() != "http://somesite.com/relative" {
		t.Error("Error getting abs url, expected http://somesite.com/relative got ", abs)
	}
	abs, err = toAbsURL(parsed, "/relative")
	if err != nil || abs.String() != "http://example.com/relative" {
		t.Error("Error getting abs url, expected http://example.com/relative got ", abs)
	}
	abs, err = toAbsURL(parsed, "./relative")
	if err != nil || abs.String() != "http://example.com/somepage/relative" {
		t.Error("Error getting abs url, expected http://example.com/somepage/relative got ", abs)
	}
}
