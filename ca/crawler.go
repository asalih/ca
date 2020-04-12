package ca

import (
	"sync"

	"github.com/asalih/colly"
	"github.com/robertkrimen/otto/ast"
)

//Crawler Crawl model
type Crawler struct {
	ScriptsManager *ScriptsManager
	url            string
	domain         string
	crawlHandler   *colly.Collector
	finished       bool
	OnRequest      func(...*crawlData)
	OnCrawlFinish  func()
}

//NewCrawler Crawler initiator
func NewCrawler(url string, domain string, scriptsManager *ScriptsManager) *Crawler {
	crawler := &Crawler{scriptsManager, url, domain, nil, false, nil, nil}
	crawler.init()

	return crawler
}

func (crawler *Crawler) init() {

	crawler.crawlHandler = colly.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		//colly.Async(true),
		colly.AllowedDomains(crawler.domain),
	)

	// On every a element which has href attribute call callback
	crawler.crawlHandler.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		crawler.crawlHandler.Visit(e.Request.AbsoluteURL(link))
	})
	crawler.crawlHandler.OnHTML("form[action]", func(e *colly.HTMLElement) {
		link := e.Attr("action")
		method := e.Attr("method")

		params := make(map[string]string)
		e.ForEach("input[name]", func(i int, ie *colly.HTMLElement) {
			//TODO Add dummy form data
			params[ie.Attr("name")] = ""
		})

		uri := e.Request.AbsoluteURL(link)
		if method == "GET" || method == "" {
			crawler.crawlHandler.Visit(uri)
			return
		}

		crawler.crawlHandler.Post(uri, params)
	})

	// Before making a request print "Visiting ..."
	crawler.crawlHandler.OnRequest(func(r *colly.Request) {
		//reqURI := r.URL.RequestURI()
		//fmt.Println("Visiting", reqURI)

		if crawler.OnRequest != nil {
			cData := newCrawlData(r.URL, r.Method)
			if r.Method == "GET" {
				crawler.OnRequest(cData.appendURLParams(r.URL.Query()))
			} else {
				crawler.OnRequest(cData.appendParams(r.Body))
			}
		}
	})

	crawler.crawlHandler.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
}

//Start starts the crawling phase
func (crawler *Crawler) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	crawler.crawlHandler.Visit(crawler.url)

	for _, p := range crawler.ScriptsManager.CrawlScripts {
		go func(scr *ast.Program) {

			result := executeCrawlScript(scr, crawler.url)
			if result == nil {
				return
			}

			crawler.OnRequest(result...)
		}(p)
	}

	crawler.crawlHandler.Wait()
	crawler.OnCrawlFinish()
	crawler.finished = true
}
