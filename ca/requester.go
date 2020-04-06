package ca

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type httpRequestHandler struct {
	crawlData *CrawlData
	attack string
}

type httpResponseHandler struct{
	URL string
	Headers map[string]string
	Body string
	Cookies []*http.Cookie
}

func (h *httpRequestHandler) Do() *httpResponseHandler{
	client := &http.Client{
		Timeout: time.Second * 45,
	}

	if h.attack != "" && h.crawlData.Method == "GET"{
		var qs string
		for _, p:= range h.crawlData.Params{
			if qs != ""{
				qs += "&"
			}

			qs += p + "=" + h.attack
		}

		h.crawlData.URL.RawQuery = qs
	}

	req, _ := http.NewRequest(h.crawlData.Method, h.crawlData.URL.String(), nil)

	resp, _ := client.Do(req)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
    }
    bodyString := string(bodyBytes)

	respMap := &httpResponseHandler{resp.Request.URL.String(), h.HeadersToString(resp.Header), bodyString, resp.Cookies()}

	return respMap
}

//HeadersToString ...
func (h *httpRequestHandler) HeadersToString(header http.Header) (map[string]string) {
	headers := make(map[string]string)

	for name, values := range header {
		hval := ""
		for _, value := range values {
			hval += value
		}
		headers[name] = hval
	}
	return headers
}