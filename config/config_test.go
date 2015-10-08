package config

import (
	"encoding/json"
	"testing"
)

const exampleConfig = `{
    "url": "http://example.pl",
    "parsing": {
        "workers": 3,
        "parseExclusions": [],
        "params": [],
        "respectRobots": true,
        "userAgent": "Googlebot 2.1",
        "excludeNonTextFiles": true,
        "requestsPerSec": 1,
        "burst": 30,
        "proxies": [
            {
                "address": "http://localhost",
                "username": "",
                "password": ""
            }
        ],
        "cutProtocol": true
    },
    "output": [
        {
            "perFile": 1000,
            "regex": "/pl/product,",
            "filePrefix": "products",
            "modifiers": {
                "priority": 1.2,
                "changeFrequency": "daily"
            }
        }
    ]
}
`

func TestConfigParsing(t *testing.T) {
	cfg := new(Config)
	err := json.Unmarshal([]byte(exampleConfig), cfg)
	if err != nil {
		t.Error(err)
	}
	//Integer parsing
	if cfg.Parsing.Workers != 3 {
		t.Errorf("Error occured while parsing workers, expected 20, got %d", cfg.Parsing.Workers)
	}
	if cfg.Parsing.CutProtocol != true {
		t.Errorf("Errror occured while parsing regex arrays, expected true got %b", cfg.Parsing.CutProtocol)
	}
}
