package sitemapgen

import (
	"net/url"
	"regexp"
	"testing"

	"github.com/maciekmm/sitemap-generator/config"
)

func TestReorderAndCrop(t *testing.T) {
	parsingConfig := &config.ParsingConfig{
		CutProtocol: true,
		Params: []*config.ParamsFilter{&config.ParamsFilter{
			Regex: &config.Regex{
				Regexp: regexp.MustCompile("somepath"),
			},
			Params:  []string{"abc"},
			Include: true,
		}, &config.ParamsFilter{
			Regex: &config.Regex{
				Regexp: regexp.MustCompile(""),
			},
			Params:  []string{"id"},
			Include: false,
		}},
		StripWWW: true,
		//StripQueryString: true,
	}

	cases := map[string]string{
		"http://www.maciekmm.net/somepath?abc=xde&id=zxc&nnn=sfv":                   "http://maciekmm.net/somepath?abc=xde",
		"http://www.example.sub.com.pl/path?querystring=abc&id=abee&abc=zqx#anchor": "http://example.sub.com.pl/path?abc=zqx&querystring=abc",
	}

	for key, res := range cases {
		url, err := url.Parse(key)
		if err != nil {
			t.Error(err)
		}
		ReorderAndCrop(parsingConfig, url)
		if url.String() != res {
			t.Error("Reordering and cropping failed, got: ", url.String(), " expected: ", res)
		}
	}
}

func TestStripProtocol(t *testing.T) {
	cases := map[string]string{
		"https://www.maciekmm.net/": "www.maciekmm.net/",
		"http://www.maciekmm.net/":  "www.maciekmm.net/",
		"www.maciekmm.net/":         "www.maciekmm.net/",
		"://www.maciekmm.net/":      "://www.maciekmm.net/",
	}
	for key, result := range cases {
		parsed := StripProtocol(key)
		if parsed != result {
			t.Error("Cutting protocol failed, got: ", parsed, " expected: ", result)
		}
	}
}

func TestStripWWW(t *testing.T) {
	cases := map[string]string{
		"www.maciekmm.net": "maciekmm.net",
		"wwe.maciekmm.net": "wwe.maciekmm.net",
		"maciekmm.net":     "maciekmm.net",
		"net.maciekmm.www": "net.maciekmm.www",
	}
	for key, result := range cases {
		parsed := StripWWW(key)
		if parsed != result {
			t.Error("Stripping WWW prefix failed, got: ", parsed, " expected: ", result)
		}
	}
}

func TestShallParse(t *testing.T) {
	parsingConfig := &config.ParsingConfig{ParseExclusions: []*config.Regex{
		&config.Regex{Regexp: regexp.MustCompile("example\\.com")},
		&config.Regex{Regexp: regexp.MustCompile(",product,[\\d]+")},
	}}
	cases := map[string]bool{
		"www.example.com":               false,
		"www.maciekmm.net":              true,
		"dummy.com/product,123":         true,
		"bear.com/page,product,234,abc": false,
		"bear.com/page,product,234":     false,
	}
	for key, result := range cases {
		parsed := ShallParse(parsingConfig, key)
		if parsed != result {
			t.Error("Testing for shall parse failed, for:", key, "got: ", parsed, " expected: ", result)
		}
	}
}
