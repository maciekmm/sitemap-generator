package sitemapgen

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/eapache/channels"
	"github.com/maciekmm/sitemap-generator/limit"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Worker struct {
	workQueue   *channels.InfiniteChannel
	validator   chan<- *url.URL
	waitGroup   *sync.WaitGroup
	generator   chan<- string
	httpClients chan *limit.Client
}

func NewWorker(workQueue *channels.InfiniteChannel, validator chan<- *url.URL, waitGroup *sync.WaitGroup, generator chan<- string, httpClients chan *limit.Client) *Worker {
	return &Worker{
		workQueue:   workQueue,
		validator:   validator,
		waitGroup:   waitGroup,
		generator:   generator,
		httpClients: httpClients,
	}
}

func (w *Worker) Start() {
	for {
		select {
		case job, ok := <-w.workQueue.Out():
			if !ok {
				return
			}

			stringURL := job.(*url.URL).String()
			req, err := http.NewRequest("GET", stringURL, nil)
			if err != nil {
				log.Println("Worker: Could not parse: ", stringURL, " error: ", err.Error())
				w.waitGroup.Done()
				continue
			}
			client := <-w.httpClients
			w.httpClients <- client
			//DEBUG: log.Println("In pool:", strconv.Itoa(w.workQueue.Len()), " - ", stringURL)
			resp, err := client.Do(req)
			if err != nil {
				log.Println("Worker: Could not connect to: ", stringURL, " error: ", err.Error())
				if strings.Contains(err.Error(), "http: error connecting to proxy") || strings.Contains(err.Error(), "while waiting for connection") {
					w.workQueue.In() <- job
				} else {
					w.waitGroup.Done()
				}
				continue
			}

			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				log.Println("Worker: Invalid status code for: ", stringURL, " code: ", resp.StatusCode)
				//TODO: return to pool on certain errors
				w.waitGroup.Done()
				continue
			}

			if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
				resp.Body.Close()
				log.Println("Worker: Invalid content-type for: ", stringURL, " content-type: ", resp.Header.Get("Content-Type"))
				w.waitGroup.Done()
				continue
			}

			//Push to file generation
			w.waitGroup.Add(1)
			w.generator <- stringURL
			//DEBUG: log.Println("Parsing ", stringURL)
			doc := html.NewTokenizer(resp.Body)
			for tokenType := doc.Next(); tokenType != html.ErrorToken; {
				token := doc.Token()
				switch tokenType {
				case html.StartTagToken:
					if token.DataAtom != atom.A {
						tokenType = doc.Next()
						continue
					}
					for _, attr := range token.Attr {
						if attr.Key == "href" {
							parsedURL, err := toAbsURL(job.(*url.URL), attr.Val)
							if err != nil {
								log.Println("Worker: Could not get an absolute path for: ", attr.Val, " error: ", err.Error())
								continue
							}
							w.waitGroup.Add(1)
							w.validator <- parsedURL
						}
					}
				}
			}
			resp.Body.Close()
			w.waitGroup.Done()
		}
	}
}

func toAbsURL(baseurl *url.URL, weburl string) (*url.URL, error) {
	relurl, err := url.Parse(weburl)
	if err != nil {
		return nil, err
	}
	absurl := baseurl.ResolveReference(relurl)
	return absurl, nil
}
