package sitemapgen

import (
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/eapache/channels"
	"github.com/maciekmm/sitemap-generator/config"
	"github.com/temoto/robotstxt-go"
)

//Validator manages address flow by pushing them to certain creation proccesses and makes sure no links are parsed twice.
type Validator struct {
	sites       map[string]bool
	workerQueue *channels.InfiniteChannel
	Input       chan *url.URL
	config      config.Config
	url         *url.URL
	waitGroup   *sync.WaitGroup
	robots      *robotstxt.RobotsData
	generator   chan string
}

//NewValidator creates a new validator instance
func NewValidator(config config.Config, workerQueue *channels.InfiniteChannel, waitGroup *sync.WaitGroup, robots *robotstxt.RobotsData, sitemapGenerator chan string) *Validator {
	parsedURL, err := url.Parse(config.URL)
	parsedURL.Host = StripWWW(parsedURL.Host)
	if err != nil {
		return nil
	}
	return &Validator{make(map[string]bool), workerQueue, make(chan *url.URL, 256), config, parsedURL, waitGroup, robots, sitemapGenerator}
}

func (v *Validator) start() {
	for {
		select {
		case url, ok := <-v.Input:
			//Stop worker if channel is closed
			if !ok {
				log.Println("Validator: channel closed -> stopping")
				return
			}
			ReorderAndCrop(v.config.Parsing, url)
			stringURL := url.String()
			URLProtStripped := stringURL
			if v.config.Parsing.CutProtocol {
				URLProtStripped = StripProtocol(stringURL)
			}
			//log.Println(stringurl)
			if _, ok := v.sites[URLProtStripped]; ok {
				//log.Println("Validator: Skipping ", stringurl, " - already looked up")
				v.waitGroup.Done()
				continue
			}
			//Mark as already processed
			v.sites[URLProtStripped] = true

			//Check for robots
			if v.config.Parsing.RespectRobots && v.robots != nil {
				if !v.robots.TestAgent(url.Path, v.config.Parsing.UserAgent) {
					log.Println("Validator: Skipping ", stringURL, " - denied by robots")
					v.waitGroup.Done()
					continue
				}
			}

			if !strings.HasSuffix(strings.ToLower(url.Host), strings.ToLower(v.url.Host)) {
				log.Println("Validator: Skipping ", stringURL, " - invalid host: ", url.Host, " expected: ", v.url.Host)
				v.waitGroup.Done()
				continue
			}

			//Push to workers
			if !ShallParse(v.config.Parsing, stringURL) {
				//log.Println("Validator: Skipping ", stringurl, " - excluded from parsing")

				//Excluding from parsing does not exclude from adding to sitemap files
				v.waitGroup.Add(1)
				v.generator <- url.String()
				v.waitGroup.Done()
				continue
			}
			v.workerQueue.In() <- url
		}
	}
}

//ReorderAndCrop removes the anchor (#sth) Fragment,
//sorts, removes and encodes query string parameters
//and lowercases Host
func ReorderAndCrop(conf *config.ParsingConfig, url *url.URL) {
	url.Path = strings.TrimSuffix(strings.TrimSpace(url.Path), "/")
	if conf.StripQueryString {
		url.RawQuery = ""
	} else {
		stringURL := url.String()
		query := url.Query()
		for _, filter := range conf.Params {
			if filter.Regex.MatchString(stringURL) {
				//If only specified params are relevent
				if filter.Include {
					for key := range query {
						//Check if param is allowed
						found := false
						for _, param := range filter.Params {
							if param == key {
								found = true
							}
						}
						//If not remove
						if !found {
							query.Del(key)
						}
					}
					//If params are irrelevant
				} else {
					for _, param := range filter.Params {
						query.Del(param)
					}
				}
			}
		}
		url.RawQuery = query.Encode()
	}
	url.Fragment = ""
	url.Host = strings.ToLower(url.Host)
	if conf.StripWWW {
		url.Host = StripWWW(url.Host)
	}
}

//StripWWW strips a www. prefix/subdomain from URL represented in string
func StripWWW(host string) string {
	if strings.HasPrefix(host, "www.") {
		return host[4:]
	}
	return host
}

//StripProtocol strips a protocol from URL represented in string
func StripProtocol(url string) string {
	if url[6] == '/' && strings.HasPrefix(url, "http") {
		if url[4] == ':' {
			return url[7:]
		}
		return url[8:]
	}
	return url
}

//ShallParse checks whether site's source should be parsed
func ShallParse(conf *config.ParsingConfig, url string) bool {
	for _, excl := range conf.ParseExclusions {
		if excl.MatchString(url) {
			return false
		}
	}
	return true
}
