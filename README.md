# Sitemap Generator
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

## Building:
1. Download required libraries
 - `go get https://github.com/eapache/channels`
 - `go get github.com/temoto/robotstxt-go`
2. Build it
 - `cd cli`
 - `go build`

## Usage:
  `./sitemap-generator -config config.json` - where config.json is location of config file

## ToDo:
 - Unit tests
 - Config documentation
 - Benchamarks
 - Adapt for library usage
