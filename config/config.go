package config

import (
	"encoding/json"
	"io/ioutil"
)

//Config describes the whole process of generating sitemap
type Config struct {
	URL     string         `json:"url"`
	Parsing *ParsingConfig `json:"parsing"`
	Output  []*Filter      `json:"output"`
}

//ParsingConfig describes site parsing process
type ParsingConfig struct {
	Workers int `json:"workers"`
	//ParseExclusions regex array defines what sites content is not parsed but included in sitemap
	ParseExclusions []*Regex `json:"parseExclusions,omitempty"`
	RespectRobots   bool     `json:"respectRobots,omitempty"`
	//ResultExclusions regex array defines what sites are being excluded from the sitemap
	Params            []*ParamsFilter `json:"params,omitempty"`
	CutProtocol       bool            `json:"cutProtocol,omitempty"`
	Proxies           []Proxy         `json:"proxies,omitempty"`
	UserAgent         string          `json:"userAgent,omitempty"`
	RequestsPerSecond int64           `json:"requestsPerSec,omitempty"`
	Burst             int             `json:"burst,omitempty"`
	StripQueryString  bool            `json:"stripQueryString,omitempty"`
	StripWWW          bool            `json:"stripWWW,omitempty"`
	NoProxyClient     bool            `json:"noProxyClient,omitempty"`
}

type ParamsFilter struct {
	Regex   *Regex   `json:"regex"`
	Params  []string `json:"params"`
	Include bool     `json:"include"`
}

//Filter defines how urls are being put into separated files sitemap
type Filter struct {
	Regex                   *Regex     `json:"regex"`
	PerFile                 int        `json:"perFile"`
	FilePrefix              string     `json:"filePrefix"`
	Modifiers               Changeable `json:"modifiers,omitempty"`
	IncludeModificationDate bool       `json:"modificationDate,omitempty"`
}

type Changeable struct {
	ChangeFrequency string  `json:"changeFrequency,omitempty" xml:"changefreq,omitempty"`
	Priority        float64 `json:"priority,omitempty" xml:"priority,omitempty"`
}

//Proxy represents a proxy configuration to be used for parsing
type Proxy struct {
	Address  string `json:"address"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

//FromFile parses file into a Config
func FromFile(file string) (*Config, error) {
	fileC, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	cfg := new(Config)
	err = json.Unmarshal(fileC, cfg)
	return cfg, err
}

//FromFiles parses regular config file for website setting and proxy config file for proxies
func FromFiles(config string, proxies string) (*Config, error) {
	cfg, err := FromFile(config)
	if err != nil {
		return nil, err
	}

	if proxies == "" {
		return cfg, err
	}

	proxyConfig, err := ioutil.ReadFile(proxies)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(proxyConfig, &cfg.Parsing.Proxies)
	return cfg, err
}
