package ca

import (
	"bufio"
	"io"
	"net/url"

	"github.com/robertkrimen/otto/ast"
)

type crawlData struct {
	URL    *url.URL
	Method string
	Params []string
	Script *ast.Program
}

func newCrawlData(u *url.URL, method string) *crawlData {
	return &crawlData{u, method, nil, nil}
}

func (c *crawlData) appendURLParams(v url.Values) *crawlData {
	c.Params = []string{}
	for qk := range v {
		c.Params = append(c.Params, qk)
	}
	return c
}

func (c *crawlData) appendParams(r io.Reader) *crawlData {
	v, _ := url.ParseQuery(bufio.NewScanner(r).Text())
	c.Params = []string{}
	for k := range v {
		c.Params = append(c.Params, k)
	}

	return c
}
