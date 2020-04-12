package ca

import (
	"container/list"
	"net/url"
	"sync"
)

//Manager General manager for crawl&attack process
type Manager struct {
	Crawler  *Crawler
	Attacker *Attacker
	queue    *list.List
	mutex    sync.Mutex
}

//NewManager Manager initator
func NewManager(uriStr string) *Manager {
	uriParsed, _ := url.Parse(uriStr)

	scriptsManager := NewScriptsManager("./scripts/")
	return &Manager{NewCrawler(uriParsed.String(), uriParsed.Host, scriptsManager), NewAttacker(scriptsManager), list.New(), sync.Mutex{}}
}

//CrawlAndAttack Starts crawl and attack process
func (m *Manager) CrawlAndAttack() {

	m.Crawler.OnRequest = func(c ...*crawlData) {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		for _, i := range c {
			m.queue.PushBack(i)
		}
	}

	m.Crawler.OnCrawlFinish = func() {
		m.Attacker.Finalize()
	}

	m.Attacker.CrawlDataReader = func() *crawlData {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		e := m.queue.Front()
		if e == nil {
			return nil
		}

		m.queue.Remove(e)

		if e.Value == nil {
			return nil
		}

		return e.Value.(*crawlData)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go m.Crawler.Start(&wg)
	go m.Attacker.Start(&wg)

	wg.Wait()
}

//IsFinished Returns whether crawl and attack has been finished
func (m *Manager) IsFinished() bool {
	return m.Crawler.finished && m.Attacker.finished
}
