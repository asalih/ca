package ca

import "net/url"

type CrawlData struct {
	URL    *url.URL
	Method string
	Params []string
}

func NewCrawlData(u *url.URL, method string) *CrawlData {
	return &CrawlData{u, method, nil}
}

func (c *CrawlData) AppendParams(v url.Values) *CrawlData {
	c.Params = []string{}
	for qk := range v {
		c.Params = append(c.Params, qk)
	}
	return c
}
