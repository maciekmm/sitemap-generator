# Sitemap Generator
[![Build Status](https://travis-ci.org/maciekmm/sitemap-generator.svg?branch=master)](https://travis-ci.org/maciekmm/sitemap-generator)

A CLI app made for generating detailed sitemaps for websites which don't have one.

**This piece of software is in early development stage and has not been fully tested, use at your own risk.**

## Features:
- Multiple workers for faster parsing and generation
- QueryString inclusion/exclusion for urls which meet regex requirements
- URL unification (avoiding duplicates)
- Streamed flow (low memory usage)
- Output filtering based on  URLs
- Proxy support for avoiding rate limiting
- Simple and powerful configuration
- `robots.txt` support

## Building CLI:
1. Download required libraries
 - `go get https://github.com/eapache/channels`
 - `go get github.com/temoto/robotstxt-go`
2. Build it
 - `cd sitemap-generator`
 - `go build`

## Usage:
  `./sitemap-generator -config config.json` - where config.json is location of config file

## Example config:

```json
{
    "url": "http://example.com/", "#":"URL that will be parsed first",
    "parsing": {
        "workers": 2, "#":"Amount of parallel workers parsing pages",
        "parseExclusions": [
            "example.com/p/",
            "http://.+\\.example\\.com", "#Don't parse subdomains",
            "\\.jsp"
        ], "#":"Which sites not to parse, this doesn't exclude it from being but into sitemap file",
        "params": [
            {
              "regex": "",
              "include": true, "#": "Include means that only params specified below will be kept, Exclude will remove given params.",
              "params": ["id"]
            }
        ], "#":"Specifiec which params should be kept and which one should be stripped",
        "respectRobots": true, "#": "Whether robots.txt should be respected",
        "userAgent": "BOT-SitemapGenerator", "#": "UserAgent for requests",
        "noProxyClient": true, "#": "Whether to create http client without proxy",
        "requestsPerSec": 1, "#": "Amount of requests per client ",
        "stripQueryString": false, "#": "Whether to completely ignore query string",
        "stripWWW": true, "#": "Whether to treat www.example.com and example.com as thesame page.",
        "burst": 2, "#": "Request burst - accumulation of unused request opportunities from request per sec",
        "cutProtocol": true, "#": "Whether to remove http(s) protocol",
        "proxies": [
            {
                "address": "http://000.000.000.000", "#": "Proxy address",
                "username": "username", "#": "Username",
                "password": "password", "#": "Password"
            }
        ]
    },
    "output": [
        {
            "perFile": 4000, "#": "How many sites per file",
            "regex": "example.com/p/", "#": "Which sites should apply",
            "filePrefix": "products",
            "modifiers": {
              "changeFrequency": "daily",
              "priority": 0.1
            }
        }
    ]
}
```

## ToDo:
 - Unit tests
 - Config documentation
 - Benchamarks
 - Adapt for library usage
