package ca

import (
	"sync"

	"github.com/asalih/colly"
)

type Crawler struct {
	ScriptsManager *ScriptsManager
	url            string
	domain         string
	crawlHandler   *colly.Collector
	finished       bool
	OnRequest      func(*CrawlData)
	OnCrawlFinish  func()
}

func NewCrawler(url string, domain string, scriptsManager *ScriptsManager) *Crawler {
	return &Crawler{scriptsManager, url, domain, nil, false, nil, nil}
}

func (crawler *Crawler) Init() {

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

		if method == "GET" || method == "" {
			crawler.crawlHandler.Visit(e.Request.AbsoluteURL(link))
			return
		}

		uri := e.Request.AbsoluteURL(link)
		visited, err := crawler.crawlHandler.HasVisited(uri)

		if visited || err != nil {
			return
		}

		crawler.crawlHandler.Post(uri, params)
	})

	// Before making a request print "Visiting ..."
	crawler.crawlHandler.OnRequest(func(r *colly.Request) {
		//reqURI := r.URL.RequestURI()
		//fmt.Println("Visiting", reqURI)

		if crawler.OnRequest != nil {
			crawler.OnRequest(NewCrawlData(r.URL, r.Method).AppendParams(r.URL.Query()))
		}
	})

	crawler.crawlHandler.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"

}

func (crawler *Crawler) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	crawler.crawlHandler.Visit(crawler.url)
	crawler.crawlHandler.Wait()
	crawler.OnCrawlFinish()
	crawler.finished = true
}
