package limit

import "net/http"

//Client struct represents a rate limited http client
type Client struct {
	client    *http.Client
	limiter   RateLimiter
	userAgent string
}

//NewClient returns a new http.Client instance
func NewClient(client *http.Client, limiter *RateLimiter, userAgent string) *Client {
	limiter.Start()
	return &Client{
		client:    client,
		limiter:   *limiter,
		userAgent: userAgent,
	}
}

//Do performs a given http.Request with ratelimiting in mind
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	c.limiter.Wait()
	req.Header.Add("User-Agent", c.userAgent)
	return c.client.Do(req)
}
